# Use a minimal Go base image
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the application
RUN go build -o tower-of-song

# Use a minimal base image for the final container
FROM ubuntu:latest

# Install necessary dependencies
RUN apt update && apt install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/tower-of-song /app/tower-of-song

# Ensure the binary is executable
RUN chmod +x /app/tower-of-song

# Copy the config file (if needed)
COPY ./config.json /app/config.json

# Expose the port our Go server listens on
EXPOSE 8080

# Run the server
CMD ["/app/tower-of-song"]
