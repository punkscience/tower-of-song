#!/bin/bash

# Tower of Song - Cleanup Script
# Clean up test containers and data

echo "🧹 Tower of Song - Cleanup"
echo "=========================="

# Stop and remove test containers
echo "Stopping test containers..."
docker stop tower-of-song-test tower-of-song-quick 2>/dev/null || true
docker rm tower-of-song-test tower-of-song-quick 2>/dev/null || true

# Remove test image
echo "Removing test image..."
docker rmi tower-of-song 2>/dev/null || true

# Clean up test directories (optional)
read -p "Do you want to remove test directories? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Removing test directories..."
    rm -rf ./test-music ./test-data
    echo "✅ Test directories removed"
else
    echo "Test directories preserved"
fi

echo "✅ Cleanup complete!" 