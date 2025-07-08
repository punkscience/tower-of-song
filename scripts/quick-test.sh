#!/bin/bash

# Tower of Song - Quick Test Script
# Simple script to quickly test the Docker image

set -e

# Source test folder config if present
CONFIG_FILE="$(dirname "$0")/test-folders.conf"
if [ -f "$CONFIG_FILE" ]; then
    echo "Using test folder config: $CONFIG_FILE"
    source "$CONFIG_FILE"
fi

# Configuration
IMAGE_NAME="tower-of-song"
CONTAINER_NAME="tower-of-song-quick"
PORT="8080"
MUSIC_DIR="${MUSIC_DIR:-./test-music}"
DATA_DIR="${DATA_DIR:-./test-data}"

echo "🎵 Tower of Song - Quick Test"
echo "============================="

# Create test directories
mkdir -p "$MUSIC_DIR" "$DATA_DIR"

# Build image
echo "Building Docker image..."
docker build -t $IMAGE_NAME .

# Stop existing container
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm $CONTAINER_NAME 2>/dev/null || true

# Run container
echo "Starting container..."
docker run -d \
    --name $CONTAINER_NAME \
    -p $PORT:8080 \
    -v "$MUSIC_DIR:/app/music:ro" \
    -v "$DATA_DIR:/app/data" \
    $IMAGE_NAME

echo "✅ Server started!"
echo "🌐 Open: http://localhost:$PORT"
echo "📋 Logs: docker logs -f $CONTAINER_NAME"
echo "🛑 Stop: docker stop $CONTAINER_NAME" 