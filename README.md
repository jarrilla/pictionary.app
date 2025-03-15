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
- **Database**: MongoDB for caching generated images

## Prerequisites

- Node.js (v16+)
- Go (v1.21+)
- MongoDB
- OpenAI API key

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/pictionary-app.git
   cd pictionary-app
   ```

2. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```

3. Install backend dependencies:
   ```bash
   cd ../backend
   go mod download
   ```

4. Configure environment variables:
   - Create a `.env` file in the backend directory with:
     ```
     OPENAI_API_KEY=your_api_key_here
     MONGODB_URI=mongodb://localhost:27017
     PORT=8080
     ```
   - Create a `.env` file in the frontend directory with:
     ```
     VITE_API_URL=http://localhost:8080
     ```

## Running the Application

1. Start MongoDB:
   ```bash
   mongod
   ```

2. Start the backend server:
   ```bash
   cd backend
   go run main.go
   ```

3. In a new terminal, start the frontend development server:
   ```bash
   cd frontend
   npm run dev
   ```

4. Open your browser and navigate to `http://localhost:5173`

## Building for Production

1. Build the frontend:
   ```bash
   cd frontend
   npm run build
   ```

2. Build the backend:
   ```bash
   cd ../backend
   go build -o pictionary-server
   ```

3. Run the production server:
   ```bash
   ./pictionary-server
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [Free Dictionary API](https://dictionaryapi.dev/) for providing word definitions
- [OpenAI](https://openai.com/) for the DALL-E 3 image generation API 