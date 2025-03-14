#!/bin/bash

# Function to display help message
show_help() {
  echo "Usage: ./run-docker.sh [OPTION]"
  echo "Run the Pictionary app using Docker."
  echo ""
  echo "Options:"
  echo "  dev       Run in development mode with hot-reloading"
  echo "  dev-d     Run in development mode with hot-reloading (detached)"
  echo "  prod      Run in production mode (default)"
  echo "  prod-d    Run in production mode (detached)"
  echo "  build     Build the Docker images"
  echo "  down      Stop and remove containers"
  echo "  init      Initialize Go modules and dependencies"
  echo "  logs      View logs from the containers"
  echo "  help      Display this help message"
  echo ""
}

# Function to initialize Go modules
init_go_modules() {
  echo "Initializing Go modules..."
  cd backend
  ./setup-go-modules.sh
  cd ..
}

# Function to display access information
show_access_info() {
  local mode=$1
  local detached=$2
  
  echo ""
  echo "==================================================="
  if [ "$mode" == "dev" ]; then
    echo "Development server started!"
    echo "Access the application at: http://localhost:5173"
    echo "Backend API available at: http://localhost:8080"
  else
    echo "Production server started!"
    echo "Access the application at: http://localhost"
    echo "Backend API available at: http://localhost:8080"
  fi
  echo "==================================================="
  echo ""
  
  if [ "$detached" == "true" ]; then
    echo "Containers are running in the background."
    echo "To view logs, run: ./run-docker.sh logs"
    echo "To stop containers, run: ./run-docker.sh down"
    echo ""
  fi
}

# Check if .env file exists in backend directory
if [ ! -f "backend/.env" ]; then
  echo "Warning: backend/.env file not found."
  echo "Creating a copy from .env.example..."
  
  if [ -f "backend/.env.example" ]; then
    cp backend/.env.example backend/.env
    echo "Created backend/.env from .env.example."
    echo "Please edit backend/.env to add your OpenAI API key."
  else
    echo "Error: backend/.env.example not found."
    echo "Please create a backend/.env file with your configuration."
    exit 1
  fi
fi

# Process command line arguments
case "$1" in
  dev)
    echo "Starting in development mode..."
    docker compose -f docker-compose.dev.yml up
    show_access_info "dev" "false"
    ;;
  dev-d)
    echo "Starting in development mode (detached)..."
    docker compose -f docker-compose.dev.yml up -d
    show_access_info "dev" "true"
    ;;
  prod)
    echo "Starting in production mode..."
    docker compose up
    show_access_info "prod" "false"
    ;;
  prod-d)
    echo "Starting in production mode (detached)..."
    docker compose up -d
    show_access_info "prod" "true"
    ;;
  build)
    echo "Building Docker images..."
    docker compose build
    echo "Building development Docker images..."
    docker compose -f docker-compose.dev.yml build
    ;;
  down)
    echo "Stopping containers..."
    docker compose down
    docker compose -f docker-compose.dev.yml down
    ;;
  logs)
    echo "Showing logs from containers..."
    echo "Development containers:"
    docker compose -f docker-compose.dev.yml logs -f
    ;;
  init)
    init_go_modules
    ;;
  help)
    show_help
    ;;
  *)
    if [ -z "$1" ]; then
      echo "Starting in production mode..."
      docker compose up
      show_access_info "prod" "false"
    else
      echo "Unknown option: $1"
      show_help
      exit 1
    fi
    ;;
esac 