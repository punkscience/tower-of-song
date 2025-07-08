# Use a minimal Go base image
FROM golang:1.23.10 AS builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Copy the templates directory explicitly for static file serving
COPY ./templates /app/templates
# Copy the static directory for CSS and other assets
COPY ./static /app/static

# Build the application
RUN go build -o tower-of-song

# Use a minimal base image for the final container
FROM ubuntu:latest

# Install necessary dependencies
RUN apt update && apt install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Ensure persistent data directory exists
RUN mkdir -p /app/data

# Copy the compiled binary from the builder stage
COPY --from=builder /app/tower-of-song /app/tower-of-song
# Copy the templates directory from the builder stage
COPY --from=builder /app/templates /app/templates
# Copy the static directory from the builder stage
COPY --from=builder /app/static /app/static

# Ensure the binary is executable
RUN chmod +x /app/tower-of-song

# Copy the config file (if needed)
COPY ./config.json /app/config.json

# Expose the port our Go server listens on
EXPOSE 8080

# Declare a volume for persistent database storage
VOLUME ["/app/data"]

# Run the server
CMD ["/app/tower-of-song"]
