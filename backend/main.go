package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new router
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/api/generate-image", generateImageHandler).Methods("POST")
	r.HandleFunc("/api/health", healthCheckHandler).Methods("GET")

	// Serve static files from the frontend build directory in production
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/dist")))

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := c.Handler(r)

	// Start the server
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func generateImageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req ImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.Word == "" || req.Definition == "" {
		http.Error(w, "Word and definition are required", http.StatusBadRequest)
		return
	}

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
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
		log.Printf("Error generating image: %v", err)

		// Check for content policy violation
		errMsg := err.Error()
		if strings.Contains(errMsg, "400 Bad Request") &&
			strings.Contains(errMsg, "safety system") {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
