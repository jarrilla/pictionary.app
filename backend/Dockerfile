# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /app

# Install git (required for go get)
RUN apk add --no-cache git

# Copy go.mod (and go.sum if it exists)
COPY go.mod ./
COPY go.sum* ./

# Copy source code
COPY . .

# Initialize modules and download dependencies
RUN go mod tidy && \
    go get github.com/gorilla/mux && \
    go get github.com/joho/godotenv && \
    go get github.com/rs/cors && \
    go get github.com/sashabaranov/go-openai && \
    go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pictionary-server .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/pictionary-server .

# Copy .env file
COPY .env* ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./pictionary-server"] 