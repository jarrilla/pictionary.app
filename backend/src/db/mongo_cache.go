package db

import (
	"context"
	"fmt"
	"time"

	"pictionary-app/backend/src/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCache represents a MongoDB-based cache for image URLs
type MongoCache struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// CacheEntry represents a cached image with its metadata
type CacheEntry struct {
	Word         string    `bson:"word"`
	PartOfSpeech string    `bson:"partOfSpeech"`
	Definition   string    `bson:"definition"`
	ImageData    string    `bson:"imageData"`
	CreatedAt    time.Time `bson:"createdAt"`
}

// NewMongoCache creates a new MongoDB cache instance
func NewMongoCache(uri string) (*MongoCache, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("Failed to connect to MongoDB: %v", err)
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Error("Failed to ping MongoDB: %v", err)
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	logger.Info("Successfully connected to MongoDB")

	// Get collection
	collection := client.Database("pictionary-app").Collection("cache")

	// Create indexes
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "word", Value: 1},
			{Key: "partOfSpeech", Value: 1},
			{Key: "definition", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Warning("Failed to create index: %v", err)
	} else {
		logger.Info("Successfully created MongoDB indexes")
	}

	return &MongoCache{
		client:     client,
		collection: collection,
	}, nil
}

// Close closes the MongoDB connection
func (c *MongoCache) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.client.Disconnect(ctx)
	if err != nil {
		logger.Error("Failed to close MongoDB connection: %v", err)
		return err
	}

	logger.Info("Successfully closed MongoDB connection")
	return nil
}

// Set adds or updates a cache entry
func (c *MongoCache) Set(word, partOfSpeech, definition, imageData string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	entry := CacheEntry{
		Word:         word,
		PartOfSpeech: partOfSpeech,
		Definition:   definition,
		ImageData:    imageData,
		CreatedAt:    time.Now(),
	}

	filter := bson.M{
		"word":         word,
		"partOfSpeech": partOfSpeech,
		"definition":   definition,
	}

	opts := options.Replace().SetUpsert(true)
	result, err := c.collection.ReplaceOne(ctx, filter, entry, opts)
	if err != nil {
		logger.Error("Failed to set cache entry: %v", err)
		return err
	}

	if result.UpsertedCount > 0 {
		logger.Info("Created new cache entry for word: %s", word)
	} else {
		logger.Info("Updated existing cache entry for word: %s", word)
	}

	return nil
}

// Get retrieves a cache entry
func (c *MongoCache) Get(word, partOfSpeech, definition string) (*CacheEntry, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"word":         word,
		"partOfSpeech": partOfSpeech,
		"definition":   definition,
	}

	var entry CacheEntry
	err := c.collection.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Debug("Cache miss for word: %s", word)
			return nil, false, nil
		}
		logger.Error("Failed to get cache entry: %v", err)
		return nil, false, err
	}

	logger.Debug("Cache hit for word: %s", word)
	return &entry, true, nil
}

// Delete removes a cache entry
func (c *MongoCache) Delete(word, partOfSpeech, definition string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"word":         word,
		"partOfSpeech": partOfSpeech,
		"definition":   definition,
	}

	result, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error("Failed to delete cache entry: %v", err)
		return err
	}

	if result.DeletedCount > 0 {
		logger.Info("Deleted cache entry for word: %s", word)
	} else {
		logger.Debug("No cache entry found to delete for word: %s", word)
	}

	return nil
}

// Clear removes all cache entries
func (c *MongoCache) Clear() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to clear cache: %v", err)
		return err
	}

	logger.Info("Cleared %d entries from cache", result.DeletedCount)
	return nil
}
