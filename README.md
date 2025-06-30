# Tower of Song

Tower of Song is a lightweight, self-hosted music streaming server written in Go. It scans your local music folders, extracts metadata, and provides a secure RESTful API for browsing, searching, and streaming your music collection. Includes a simple HTML test client for API exploration.

## Features

- Scan and index music files (MP3, FLAC, WAV) from configurable folders
- Extract metadata (title, artist, album) from ID3v2 tags
- Secure RESTful API for:
  - Listing all music
  - Searching by title, artist, album, or path
  - Streaming audio files
  - Viewing library statistics
- Token-based authentication for all endpoints
- Dockerized for easy deployment
- Simple HTML/JavaScript test client

## Quick Start

### With Docker

```bash
git clone https://github.com/yourusername/tower-of-song.git
cd tower-of-song
docker build -t tower-of-song .
docker run -p 8080:8080 -v /path/to/music:/app/music tower-of-song
```

To stop the server:
```bash
docker ps                # Find the container ID or name
docker stop <container>  # Replace <container> with the ID or name
```

### Local Development

Requirements:
- Go 1.22+
- SQLite3 (for development)

```bash
go mod download
go build -o tower-of-song
./tower-of-song
```

## Configuration

Edit `config.json` to specify your music folders:
```json
{
    "music_folders": ["/app/music"]
}
```

## Authentication

All API endpoints (except `/login`) require authentication.
- Obtain a token by POSTing to `/login`:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"password"}' http://localhost:8080/login
  ```
- Use the returned token in the `Authorization` header for all other requests:
  ```bash
  curl -H "Authorization: token-..." http://localhost:8080/stats
  ```
- Default credentials: `admin` / `password` (change in production!)

## API Endpoints

- `POST /login` — Authenticate and receive a token
- `GET /stats` — Get library statistics (requires token)
- `GET /list` — List all music files (requires token)
- `GET /search?q=...` — Search music (requires token)
- `GET /stream?id=...` — Stream audio file (requires token)

See [docs/tech-specification.md](docs/tech-specification.md) for full API details.

## Test Client

A simple HTML client is provided in `templates/index.html`:
- Login with your credentials
- Test all API endpoints
- Stream audio directly in your browser

## Security

- All endpoints require authentication (except `/login`)
- CORS enabled for cross-origin requests
- In-memory token store (tokens lost on restart)
- No user management or persistent database yet

## Roadmap

- Persistent database storage
- User management and roles
- Playlist support
- Audio transcoding
- Improved search and metadata editing

## Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, features, or documentation improvements.

## License

MIT License 