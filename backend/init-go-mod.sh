#!/bin/sh

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
  echo "Initializing go.mod..."
  go mod init pictionary-app/backend
fi

# Install required dependencies
echo "Installing dependencies..."
go get github.com/gorilla/mux
go get github.com/joho/godotenv
go get github.com/rs/cors
go get github.com/sashabaranov/go-openai

# Ensure dependencies are downloaded and go.sum is created
echo "Running go mod tidy..."
go mod tidy

echo "Go module initialization complete." 