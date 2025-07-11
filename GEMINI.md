
# Gemini Project Context: Tower of Song

This document provides context for the Gemini CLI to understand the Tower of Song project.

## Project Overview

**Tower of Song** is a lightweight, self-hosted music streaming server written in Go. It's designed to be run in a Docker container, and it provides a RESTful API for accessing a music library. The project also includes a modern, dark-themed web interface for browsing and playing music, built with Tailwind CSS.

### Key Features

*   **Music Library Management**: Scans and indexes music files from specified folders.
*   **Metadata Extraction**: Reads ID3v2 tags to get song information.
*   **RESTful API**: Provides endpoints for authentication, library statistics, file listing, searching, and streaming.
*   **Web Interface**: A user-friendly, responsive web client for interacting with the music library.
*   **Dockerized**: The application is designed to be run in a Docker container, with a multi-stage Dockerfile for optimized builds.
*   **Authentication**: Uses a token-based authentication system.

## How to Work with This Project

### Building and Running

The project can be built and run using standard Go commands:

*   **Build**: `go build -o tower-of-song`
*   **Run**: `./tower-of-song`

The server will start on port `8080`.

### Testing

The project includes shell scripts for testing:

*   **Quick Test**: `./scripts/quick-test.sh`
*   **Full Test**: `./scripts/test-local.sh`
*   **Cleanup**: `./scripts/cleanup.sh`

### Configuration

The application is configured through `config.json`. This file specifies the music folders to scan and the credentials for authentication.

### API

The API is documented in `docs/api-specification.md`. Key endpoints include:

*   `POST /login`: Authenticate and get a token.
*   `GET /stats`: Get library statistics.
*   `GET /list`: List all music files.
*   `GET /search`: Search the music library.
*   `GET /stream`: Stream a music file.
*   `GET /trackinfo`: Get information about a specific track.

### Code Structure

*   `main.go`: The main application file, containing the web server and API endpoint handlers.
*   `go.mod`, `go.sum`: Go module files for dependency management.
*   `Dockerfile`: For building the Docker image.
*   `templates/`: Contains the HTML for the web interface.
*   `static/`: Contains static assets like CSS and favicons.
*   `docs/`: Contains project documentation.
*   `scripts/`: Contains helper scripts for testing and other tasks.

## Gemini's Role

When working on this project, Gemini should:

*   **Understand the Go code in `main.go`**: This is the core of the application.
*   **Be aware of the API structure**: Refer to `docs/api-specification.md` when making changes to the API.
*   **Use the provided test scripts**: To verify changes and ensure that the application is working correctly.
*   **Follow existing code style**: Maintain the existing code style and conventions.
