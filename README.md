# Tower of Song

Tower of Song is a lightweight, self-hosted music streaming server written in Go. It scans your local music folders, extracts metadata, and provides a secure RESTful API for browsing, searching, and streaming your music collection. It features a beautiful, responsive web interface that works seamlessly on desktops, tablets, and mobile devices.

## Features

- **🎵 Music Management**: Scan and index music files (MP3, FLAC, WAV) from configurable folders.
- **📊 Metadata Extraction**: Automatically extracts metadata (title, artist, album) from ID3v2 tags.
- **🔐 Secure API**: Token-based authentication for all endpoints.
- **🌐 Responsive Web UI**: A beautiful, modern web interface that adapts to any screen size.
- **🐳 Docker Ready**: Multi-platform Docker images (ARM64 for Raspberry Pi, AMD64 for servers).
- **🔄 CI/CD Pipeline**: Automated testing, security scanning, and Docker image builds on every commit.

## Installation

### With Docker (Recommended)

1.  **Pull the latest image:**
    ```bash
    docker pull punkscience/tower-of-song:latest
    ```

2.  **Run the container:**
    ```bash
    docker run -d --name tower-of-song \
      -p 8080:8080 \
      -v /path/to/your/music:/app/music:ro \
      -v /path/to/data:/app/data \
      punkscience/tower-of-song:latest
    ```
    - Replace `/path/to/your/music` with the path to your music library.
    - Replace `/path/to/data` with a path to store the application's database.

### From Source

**Prerequisites:**
- Go 1.21+
- Git

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/punkscience/tower-of-song.git
    cd tower-of-song
    ```

2.  **Download dependencies:**
    ```bash
    go mod download
    ```

3.  **Build and run the application:**
    ```bash
    go build -o tower-of-song
    ./tower-of-song
    ```

## Configuration

Edit `config.json` to specify your music folders and user credentials:
```json
{
    "music_folders": ["/app/music"],
    "username": "admin",
    "password": "your-secure-password"
}
```
- The `music_folders` array can contain multiple paths.
- **Important:** Change the default `username` and `password` in a production environment.

## Usage

1.  **Access the web interface:**
    Open your browser and navigate to `http://localhost:8080`.

2.  **Login:**
    Use the credentials you configured in `config.json`.

3.  **Enjoy your music:**
    - Search for songs, artists, or albums.
    - Play, pause, and skip tracks.
    - Add your favorite songs to a dedicated list.

## CI/CD Pipeline

This project uses GitHub Actions for continuous integration and deployment. The pipeline includes the following stages:
- **Testing**: Unit and integration tests are run on every commit to any branch.
- **Security Scanning**: The codebase is scanned for vulnerabilities using `govulncheck`.
- **Docker Build & Push**: A multi-platform Docker image is built and pushed to Docker Hub on every commit to the `main` branch.

## Contributing

Contributions are welcome! Please follow these steps:
1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes and ensure all tests pass.
4.  Submit a pull request with a clear description of your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
