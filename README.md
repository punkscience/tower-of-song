# Tower of Song

Tower of Song is a lightweight, self-hosted music streaming server written in Go. It scans your local music folders, extracts metadata, and provides a secure RESTful API for browsing, searching, and streaming your music collection. Features a beautiful newspaper-themed web interface for easy music management.

## Features

- **🎵 Music Management**: Scan and index music files (MP3, FLAC, WAV) from configurable folders
- **📊 Metadata Extraction**: Extract metadata (title, artist, album) from ID3v2 tags
- **🔐 Secure API**: Token-based authentication for all endpoints
- **🌐 Beautiful UI**: Newspaper-themed web interface with modern typography
- **🐳 Docker Ready**: Multi-platform Docker images (ARM64 for Raspberry Pi, AMD64 for servers)
- **🔄 CI/CD Pipeline**: Automated testing, security scanning, and Docker image builds
- **📱 Responsive Design**: Works on desktop, tablet, and mobile devices

## Quick Start

### With Docker (Recommended)

```bash
# Pull the latest image
docker pull punkscience/tower-of-song:latest

# Run with your music library
docker run -d --name tower-of-song \
  -p 8080:8080 \
  -v /path/to/your/music:/app/music:ro \
  -v /path/to/data:/app/data \
  punkscience/tower-of-song:latest
```

### Local Development

**Prerequisites:**
- Go 1.23.10+
- Docker (for testing)

**Quick Test:**
```bash
# Use the provided test scripts
./scripts/quick-test.sh
```

**Manual Setup:**
```bash
go mod download
go build -o tower-of-song
./tower-of-song
```

## Configuration

Edit `config.json` to specify your music folders and credentials:
```json
{
    "music_folders": ["/app/music"],
    "username": "admin",
    "password": "your-secure-password"
}
```

## Authentication

All API endpoints (except `/login`) require authentication:
- **Login**: POST to `/login` with username/password
- **Token**: Use returned token in `Authorization` header
- **Default**: `admin` / `password` (change in production!)

**Example:**
```bash
# Login
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' \
  http://localhost:8080/login

# Use token
curl -H "Authorization: token-..." http://localhost:8080/stats
```

## API Endpoints

- `POST /login` — Authenticate and receive a token
- `GET /stats` — Get library statistics (requires token)
- `GET /list` — List all music files (requires token)
- `GET /search?q=...` — Search music (requires token)
- `GET /stream?id=...` — Stream audio file (requires token)
- `GET /trackinfo?id=...` — Get track details (requires token)

See [docs/api-specification.md](docs/api-specification.md) for complete API documentation and examples.

## Web Interface

Access the beautiful newspaper-themed interface at `http://localhost:8080`:
- **Modern Design**: Clean, responsive layout with authentic typography
- **Easy Navigation**: Search, browse, and stream your music collection
- **Real-time Updates**: Live statistics and now-playing information
- **Mobile Friendly**: Works perfectly on all devices

## Development Workflow

### Branch Strategy
- **`main`**: Production-ready code with automated CI/CD
- **`dev`**: Active development branch
- **Feature branches**: Create from `dev` for new features

### Local Testing
```bash
# Full test with logs and cleanup
./scripts/test-local.sh

# Quick test
./scripts/quick-test.sh

# Cleanup
./scripts/cleanup.sh
```

### CI/CD Pipeline
- **Automated Testing**: Go tests and security scanning
- **Multi-platform Builds**: ARM64 (Raspberry Pi) + AMD64 (servers)
- **Docker Hub**: Automatic image publishing on successful builds
- **Security**: Vulnerability scanning with `govulncheck`

## Deployment

### Raspberry Pi
See [docs/raspberry-pi-setup.md](docs/raspberry-pi-setup.md) for detailed setup instructions.

### Production
```bash
# Pull latest image
docker pull punkscience/tower-of-song:latest

# Run with persistent data
docker run -d --name tower-of-song \
  --restart unless-stopped \
  -p 8080:8080 \
  -v /path/to/music:/app/music:ro \
  -v /path/to/data:/app/data \
  punkscience/tower-of-song:latest
```

## Security

- **Authentication Required**: All endpoints protected (except login)
- **Token-based**: Secure token authentication system
- **CORS Enabled**: Cross-origin request support
- **Persistent Database**: SQLite with WAL mode for data integrity
- **Security Scanning**: Automated vulnerability checks in CI/CD

## Documentation

- [API Specification](docs/api-specification.md) - Complete API reference
- [Raspberry Pi Setup](docs/raspberry-pi-setup.md) - Pi deployment guide
- [Development Workflow](docs/development-workflow.md) - Development process
- [Technical Specification](docs/tech-specification.md) - Architecture details
- [Troubleshooting](docs/troubleshooting.md) - Common issues and solutions

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch from `dev`
3. Make your changes
4. Test with the provided scripts
5. Submit a pull request

## License

MIT License 