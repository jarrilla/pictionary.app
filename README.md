# Pictionary App

A dictionary application that provides word definitions and AI-generated illustrations.

## Features

- Look up word definitions using the [Free Dictionary API](https://dictionaryapi.dev/)
- View AI-generated illustrations for each word using OpenAI's DALL-E 3
- Responsive design for desktop and mobile devices
- Support the project through one-time or monthly donations

## Tech Stack

- **Frontend**: Vue 3, TypeScript, Vite
- **Backend**: Go (Golang)
- **APIs**: Free Dictionary API, OpenAI API (DALL-E 3)
- **Development**: Air (Go hot-reloading), Vite (Vue hot-reloading)
- **Deployment**: Docker, Docker Compose

## Getting Started

### Prerequisites

- Node.js (v16+)
- Go (v1.21+)
- OpenAI API key
- Docker and Docker Compose (for containerized deployment)

### Installation

#### Option 1: Standard Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/pictionary-app.git
   cd pictionary-app
   ```

2. Install frontend dependencies:
   ```
   cd frontend
   npm install
   ```

3. Install backend dependencies:
   ```
   cd ../backend
   go mod download
   ```

4. Configure environment variables:
   - Copy the `.env.example` file to `.env` in the backend directory
   - Add your OpenAI API key to the `.env` file

### Running the Application

#### Option 1: Standard Run

1. Start the frontend development server:
   ```
   cd frontend
   npm run dev
   ```

2. Start the backend server:
   ```
   cd ../backend
   go run main.go
   ```

3. Open your browser and navigate to `http://localhost:5173`

#### Option 2: Docker Deployment

1. Configure environment variables:
   - Copy the `.env.example` file to `.env` in the backend directory
   - Add your OpenAI API key to the `.env` file

2. Initialize Go modules (if needed):
   ```
   ./run-docker.sh init
   ```

3. Use the provided script to run the application:
   ```
   ./run-docker.sh [OPTION]
   ```
   
   Available options:
   - `dev`: Run in development mode with hot-reloading (Vue + Air)
   - `dev-d`: Run in development mode with hot-reloading (detached)
   - `prod`: Run in production mode (default)
   - `prod-d`: Run in production mode (detached)
   - `build`: Build the Docker images
   - `down`: Stop and remove containers
   - `logs`: View logs from the containers
   - `init`: Initialize Go modules and dependencies
   - `help`: Display help message

   For example, to run in development mode in the background:
   ```
   ./run-docker.sh dev-d
   ```

4. Access the application:
   - Development mode: `http://localhost:5173`
   - Production mode: `http://localhost`

## Development Features

### Hot Reloading

The development environment is configured with hot-reloading for both frontend and backend:

- **Frontend**: Vite provides hot module replacement for Vue components
- **Backend**: Air automatically rebuilds and restarts the Go application when files change

This allows for a smooth development experience without manual restarts.

### Docker Development Environment

The Docker development setup includes:

- Volume mounts for real-time code changes
- Proper port mapping for accessing the application
- Persistent Go module and build caches for faster rebuilds
- Automatic restart of containers if they crash

## Troubleshooting Docker Setup

If you encounter issues with the Docker setup, try the following:

1. Make sure Docker and Docker Compose are installed and running:
   ```
   docker --version
   docker compose version
   ```

2. Initialize Go modules and dependencies:
   ```
   ./run-docker.sh init
   ```

3. Check if the `.env` file exists in the backend directory:
   ```
   ls -la backend/.env
   ```

4. Rebuild the Docker images:
   ```
   ./run-docker.sh build
   ```

5. Check Docker logs:
   ```
   ./run-docker.sh logs
   ```

6. If you can't access the application at http://localhost:5173, check if the ports are properly exposed:
   ```
   docker ps
   ```
   
   You should see the frontend container with port 5173 mapped to the host.

## Building for Production

### Option 1: Standard Build

1. Build the frontend:
   ```
   cd frontend
   npm run build
   ```

2. Build the backend:
   ```
   cd ../backend
   go build -o pictionary-server
   ```

3. Run the production server:
   ```
   ./pictionary-server
   ```

### Option 2: Docker Build

1. Build and run the Docker images:
   ```
   ./run-docker.sh build
   ./run-docker.sh prod-d
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [Free Dictionary API](https://dictionaryapi.dev/) for providing word definitions
- [OpenAI](https://openai.com/) for the DALL-E 3 image generation API
- [Air](https://github.com/cosmtrek/air) for Go hot-reloading 