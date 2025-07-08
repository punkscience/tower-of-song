#!/bin/bash

# Tower of Song - Local Test Script
# This script builds and runs the Docker image locally for testing

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Source test folder config if present
CONFIG_FILE="$(dirname "$0")/test-folders.conf"
if [ -f "$CONFIG_FILE" ]; then
    echo -e "${BLUE}Using test folder config: $CONFIG_FILE${NC}"
    source "$CONFIG_FILE"
fi

# Configuration
IMAGE_NAME="tower-of-song"
CONTAINER_NAME="tower-of-song-test"
PORT="8080"
MUSIC_DIR="${MUSIC_DIR:-./test-music}"
DATA_DIR="${DATA_DIR:-./test-data}"

echo -e "${BLUE}🎵 Tower of Song - Local Test Script${NC}"
echo "=================================="

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    docker stop $CONTAINER_NAME 2>/dev/null || true
    docker rm $CONTAINER_NAME 2>/dev/null || true
    echo -e "${GREEN}Cleanup complete!${NC}"
}

# Set up cleanup on script exit
trap cleanup EXIT

# Create test directories if they don't exist
echo -e "${BLUE}Setting up test directories...${NC}"
mkdir -p "$MUSIC_DIR"
mkdir -p "$DATA_DIR"

# Check if music directory has any files
if [ -z "$(ls -A $MUSIC_DIR 2>/dev/null)" ]; then
    echo -e "${YELLOW}Warning: Music directory is empty.${NC}"
    echo -e "${YELLOW}You can add some MP3 files to $MUSIC_DIR for testing.${NC}"
    echo -e "${YELLOW}For now, the server will start with an empty library.${NC}"
fi

# Build the Docker image
echo -e "${BLUE}Building Docker image...${NC}"
docker build -t $IMAGE_NAME .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Docker image built successfully!${NC}"
else
    echo -e "${RED}❌ Failed to build Docker image${NC}"
    exit 1
fi

# Stop and remove existing container if it exists
echo -e "${BLUE}Stopping any existing test container...${NC}"
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm $CONTAINER_NAME 2>/dev/null || true

# Run the container
echo -e "${BLUE}Starting Tower of Song server...${NC}"
docker run -d \
    --name $CONTAINER_NAME \
    -p $PORT:8080 \
    -v "$MUSIC_DIR:/app/music:ro" \
    -v "$DATA_DIR:/app/data" \
    $IMAGE_NAME

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Container started successfully!${NC}"
else
    echo -e "${RED}❌ Failed to start container${NC}"
    exit 1
fi

# Wait for server to start
echo -e "${BLUE}Waiting for server to start...${NC}"
sleep 5

# Check if server is running
if curl -s http://localhost:$PORT > /dev/null; then
    echo -e "${GREEN}✅ Server is running and responding!${NC}"
    echo -e "${GREEN}🌐 Open your browser and go to: http://localhost:$PORT${NC}"
    echo -e "${GREEN}🔑 Default login: admin / (check config.json for password)${NC}"
    echo ""
    echo -e "${BLUE}Container logs:${NC}"
    echo "=================================="
    
    # Show logs and keep container running
    echo -e "${YELLOW}Press Ctrl+C to stop the server and cleanup${NC}"
    docker logs -f $CONTAINER_NAME
else
    echo -e "${RED}❌ Server is not responding${NC}"
    echo -e "${BLUE}Container logs:${NC}"
    docker logs $CONTAINER_NAME
    exit 1
fi 