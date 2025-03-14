#!/bin/bash

echo "Setting up Go modules for the backend..."

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
  echo "Initializing go.mod..."
  go mod init pictionary-app/backend
else
  echo "go.mod already exists."
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

echo "Go module setup complete."
echo "You can now build the Docker containers with: ./run-docker.sh build" 