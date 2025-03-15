package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"pictionary-app/backend/src/db"
	"pictionary-app/backend/src/logger"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/sashabaranov/go-openai"
)

// ImageRequest represents the request body for generating an image
type ImageRequest struct {
	Word         string `json:"word"`
	PartOfSpeech string `json:"partOfSpeech"`
	Definition   string `json:"definition"`
}

// ImageResponse represents the response body for the generated image
type ImageResponse struct {
	ImageURL string `json:"imageUrl"`
	Error    string `json:"error,omitempty"`
}

// CacheLookupRequest represents the request body for cache lookups
type CacheLookupRequest struct {
	Word         string `json:"word"`
	PartOfSpeech string `json:"partOfSpeech"`
	Definition   string `json:"definition"`
}

// Global cache instance
var imageCache *db.MongoCache

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize logger
	logDir := filepath.Join("logs")
	if err := logger.Init(logDir); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Starting Pictionary App backend server...")

	// Get MongoDB URI from environment variable or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
		logger.Warning("MONGODB_URI not set, using default: %s", mongoURI)
	}

	// Initialize MongoDB cache
	cache, err := db.NewMongoCache(mongoURI)
	if err != nil {
		logger.Error("Failed to initialize MongoDB cache: %v", err)
		os.Exit(1)
	}
	imageCache = cache
	defer imageCache.Close()

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Warning("PORT not set, using default: %s", port)
	}

	// Create a new router
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/api/generate-image", generateImageHandler).Methods("POST")
	r.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	r.HandleFunc("/api/cache", cacheHandler).Methods("GET")

	// Add 404 handler for undefined API routes
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// Serve static files from the frontend build directory in production
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/dist")))

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("FRONTEND_URL"), ","),
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := c.Handler(r)

	// Start the server
	logger.Info("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Error("Server failed to start: %v", err)
		os.Exit(1)
	}
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	// Get parameters from query string
	word := r.URL.Query().Get("word")
	partOfSpeech := r.URL.Query().Get("partOfSpeech")
	definition := r.URL.Query().Get("definition")

	logger.Debug("Cache lookup request - Word: %s, PartOfSpeech: %s", word, partOfSpeech)

	// Validate required parameters
	if word == "" || definition == "" {
		logger.Warning("Invalid cache request - missing required parameters")
		http.Error(w, "Word and definition are required", http.StatusBadRequest)
		return
	}

	// Try to get the image from cache
	entry, exists, err := imageCache.Get(word, partOfSpeech, definition)
	if err != nil {
		logger.Error("Failed to retrieve from cache: %v", err)
		http.Error(w, "Failed to check cache", http.StatusInternalServerError)
		return
	}
	if !exists {
		logger.Debug("Cache miss for word: %s", word)
		http.Error(w, "Image not found in cache", http.StatusNotFound)
		return
	}

	// Return the cached image URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ImageResponse{
		ImageURL: entry.ImageURL,
	})
}

func generateImageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req ImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Warning("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger.Info("Generating image for word: %s (%s)", req.Word, req.PartOfSpeech)

	// Validate the request
	if req.Word == "" || req.Definition == "" {
		logger.Warning("Missing required fields in request")
		http.Error(w, "Word and definition are required", http.StatusBadRequest)
		return
	}

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Error("OpenAI API key not configured")
		http.Error(w, "OpenAI API key not configured", http.StatusInternalServerError)
		return
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Create the prompt for DALL-E
	prompt := fmt.Sprintf("Please draw a simple sketch of this: %s (%s) %s",
		req.Word,
		req.PartOfSpeech,
		req.Definition)

	logger.Debug("DALL-E prompt: %s", prompt)

	// Generate the image using DALL-E 3
	resp, err := client.CreateImage(r.Context(), openai.ImageRequest{
		Model:          openai.CreateImageModelDallE3,
		Prompt:         prompt,
		N:              1,
		Size:           openai.CreateImageSize1024x1024,
		Quality:        openai.CreateImageQualityStandard,
		Style:          openai.CreateImageStyleNatural,
		ResponseFormat: openai.CreateImageResponseFormatURL,
	})

	if err != nil {
		logger.Error("Failed to generate image: %v", err)

		// Check for content policy violation
		errMsg := err.Error()
		if strings.Contains(errMsg, "400 Bad Request") &&
			strings.Contains(errMsg, "safety system") {
			logger.Warning("Content policy violation for word: %s", req.Word)
			// Return a more user-friendly error message
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ImageResponse{
				Error: "Your request was rejected by the AI safety system. Please try a different word or definition.",
			})
			return
		}

		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully generated image for word: %s", req.Word)

	// Cache the generated image URL
	err = imageCache.Set(req.Word, req.PartOfSpeech, req.Definition, resp.Data[0].URL)
	if err != nil {
		logger.Error("Failed to cache image: %v", err)
	}

	// Return the image URL
	imageResp := ImageResponse{
		ImageURL: resp.Data[0].URL,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode the response
	json.NewEncoder(w).Encode(imageResp)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Health check request received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger.Warning("404 Not Found: %s", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Route not found"})
}
