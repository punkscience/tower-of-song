# Tower of Song - Technical Specification

## Overview

Tower of Song is a lightweight, self-hosted music streaming server written in Go. It provides a RESTful API for browsing, searching, and streaming music files from local directories. The application automatically scans configured music folders, extracts metadata from audio files, and serves them through a beautiful newspaper-themed web interface.

## Architecture

### Technology Stack

- **Backend**: Go 1.23.10
- **Database**: SQLite (persistent on-disk at `/app/data/towerofsong.db`, WAL mode)
- **Audio Metadata**: ID3v2 tags via `github.com/bogem/id3v2`
- **Containerization**: Docker with multi-stage build
- **Frontend**: HTML/JavaScript with newspaper-themed design
- **Styling**: Tailwind CSS with Google Fonts (Playfair Display, Source Sans Pro)
- **CI/CD**: GitHub Actions with automated testing and security scanning

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Music Files   │    │   Tower of      │    │   Web Client    │
│   (Local FS)    │◄──►│   Song Server   │◄──►│   (Browser)     │
│                 │    │   (Go 1.23.10)  │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   SQLite DB     │
                       │   (Persistent)  │
                       └─────────────────┘
```

## Core Features

### 1. Music Library Management

**Automatic Scanning**
- Recursively scans configured music folders
- Supports multiple audio formats: MP3, FLAC, WAV
- Runs every 24 hours in background goroutine
- Thread-safe scanning with mutex protection

**Metadata Extraction**
- Reads ID3v2 tags from audio files
- Extracts title, artist, and album information
- Falls back to filename if metadata is unavailable
- Handles missing or corrupted metadata gracefully

**Database Schema**
```sql
CREATE TABLE music (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT UNIQUE,
    title TEXT,
    artist TEXT,
    album TEXT
)
```

### 2. RESTful API Endpoints

#### `/stats` - Library Statistics
- **Method**: GET
- **Authentication**: Required
- **Response**: JSON with total file count
- **Example**: `{"total_files": 1250}`

#### `/list` - List All Music
- **Method**: GET
- **Authentication**: Required
- **Response**: JSON array of music files
- **Sorting**: By artist ASC, then title ASC
- **Fields**: id, path, title, artist, album

#### `/search` - Search Music Library
- **Method**: GET
- **Authentication**: Required
- **Parameters**: `q` (search query)
- **Search Scope**: title, artist, album, file path
- **Matching**: Case-insensitive LIKE queries
- **Response**: JSON array of matching files

#### `/stream` - Stream Audio File
- **Method**: GET
- **Authentication**: Required
- **Parameters**: `id` (file ID from database)
- **Response**: Audio stream with `audio/mpeg` content type
- **Features**: Direct file streaming, no transcoding

#### `/trackinfo` - Get Track Details
- **Method**: GET
- **Authentication**: Required
- **Parameters**: `id` (track id)
- **Response**: JSON object with id, path, title, artist, album
- **Example**:
  ```json
  {
    "id": "123",
    "path": "/app/music/artist/album/song.mp3",
    "title": "Song Title",
    "artist": "Artist Name",
    "album": "Album Name"
  }
  ```
- **Errors**:
  - 400 if id is missing
  - 404 if track not found

### 3. Web Interface

**Design Philosophy**
- Newspaper-themed interface with authentic typography
- Responsive design for desktop, tablet, and mobile
- Clean, professional appearance with modern UX

**Technical Implementation**
- **Fonts**: Google Fonts (Playfair Display for headlines, Source Sans Pro for body)
- **Styling**: Tailwind CSS with custom newspaper theme
- **Layout**: CSS Grid for responsive column layout
- **Interactions**: Smooth hover effects and transitions

**Features**
- **Authentication Form**: Secure login with token storage
- **Statistics Display**: Real-time library information
- **Search Interface**: Easy music discovery
- **Audio Player**: Integrated streaming with metadata display
- **API Response Viewer**: Debug and development tool

### 4. Cross-Origin Resource Sharing (CORS)

All API endpoints support CORS with the following headers:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

### 5. Configuration Management

**Configuration File**: `config.json`
```json
{
    "music_folders": ["/app/music"],
    "username": "admin",
    "password": "your-secure-password"
}
```

**Features**:
- JSON-based configuration
- Support for multiple music directories
- Runtime configuration loading
- Docker volume mounting support
- Secure credential storage

## Authentication and Security

### Authentication Model

- All API endpoints (except `/login`) require authentication via a token
- Credentials (username and password) are read from `config.json`
- Clients must POST to `/login` with a JSON body containing `username` and `password`
- On successful authentication, the server returns a token
- The client must include this token in the `Authorization` header for all subsequent API requests
- If the token is missing or invalid, the server responds with `401 Unauthorized`
- The server uses an in-memory token store (tokens are lost on restart)

### Security Features

- **Token-based Authentication**: Secure token system for API access
- **CORS Protection**: Configurable cross-origin request handling
- **Input Validation**: SQL injection protection and parameter validation
- **Security Scanning**: Automated vulnerability checks in CI/CD pipeline
- **HTTPS Ready**: Designed for secure deployment with reverse proxies

### API Changes

- All endpoints (`/stats`, `/list`, `/search`, `/stream`, `/trackinfo`) require a valid token
- New endpoint: `/login` (POST)

#### `/login` - Authenticate User
- **Method**: POST
- **Request Body**:
  ```json
  { "username": "admin", "password": "password" }
  ```
- **Response**:
  ```json
  { "token": "token-..." }
  ```
- **Error**: 401 Unauthorized if credentials are invalid

#### Authenticated Requests
- All requests to `/stats`, `/list`, `/search`, `/stream`, `/trackinfo` must include:
  - Header: `Authorization: <token>`

#### Error Responses
- `401 Unauthorized` if token is missing or invalid

## Implementation Details

### Database Management

- **Type**: Persistent SQLite database file at `/app/data/towerofsong.db`
- **Mode**: WAL (Write-Ahead Logging) enabled for better concurrency
- **Persistence**: Data is preserved across server restarts and container reboots
- **Thread Safety**: Mutex-protected database operations

### File Processing

**Supported Formats**:
- MP3 (`.mp3`)
- FLAC (`.flac`)
- WAV (`.wav`)

**Metadata Handling**:
- Uses `github.com/bogem/id3v2` library
- Graceful fallback for missing tags
- Automatic cleanup of tag resources
- Error logging for corrupted files

### Concurrency Model

- **Main Thread**: HTTP server handling requests
- **Background Thread**: Music folder scanning (24-hour intervals)
- **Synchronization**: Mutex for database operations
- **Non-blocking**: File operations don't block HTTP requests

### Error Handling

- **File Not Found**: Returns 404 status
- **Database Errors**: Logged to console
- **Metadata Errors**: Logged with file path
- **Configuration Errors**: Application exits gracefully
- **Authentication Errors**: Returns 401 Unauthorized

## Deployment

### Docker Deployment

**Multi-stage Build**:
1. **Builder Stage**: Go 1.23.10 compilation
2. **Runtime Stage**: Minimal Ubuntu image

**Dockerfile Features**:
- Optimized image size
- Security certificates included
- SQLite3 runtime dependencies
- Port 8080 exposed
- Static assets (CSS, templates) included

**Multi-platform Support**:
- **ARM64**: Raspberry Pi and ARM servers
- **AMD64**: x86_64 servers and desktops

**Usage**:
```bash
docker pull punkscience/tower-of-song:latest
docker run -d --name tower-of-song \
  -p 8080:8080 \
  -v /path/to/music:/app/music:ro \
  -v /path/to/data:/app/data \
  punkscience/tower-of-song:latest
```

### Local Development

**Prerequisites**:
- Go 1.23.10+
- Docker (for testing)

**Build and Run**:
```bash
go mod download
go build -o tower-of-song
./tower-of-song
```

**Testing**:
```bash
# Use provided test scripts
./scripts/quick-test.sh
./scripts/test-local.sh
```

## CI/CD Pipeline

### GitHub Actions Workflow

**Trigger**: Push to `main` branch or pull requests
**Jobs**:
1. **Testing Phase**:
   - Go module download and caching
   - Unit tests execution
   - Security vulnerability scanning with `govulncheck`

2. **Build Phase** (only if tests pass):
   - Multi-platform Docker image build (ARM64, AMD64)
   - Image tagging with version metadata
   - Push to Docker Hub registry

**Docker Image Tags**:
- `latest` - Latest successful build
- `main` - Latest main branch build
- `sha-{commit}` - Specific commit builds
- `v1.0.0` - Semantic version tags

### Security Scanning

- **Automated**: Runs on every build
- **Tool**: `govulncheck` for Go vulnerability detection
- **Scope**: Standard library and dependencies
- **Action**: Build fails if vulnerabilities detected

## Testing

### Test Client

Located in `templates/index.html`, provides a comprehensive web interface for:
- Testing all API endpoints
- Audio streaming functionality
- JSON response visualization
- Interactive search capabilities
- Beautiful newspaper-themed design

### API Testing

**Manual Testing**:
```bash
# Login
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' \
  http://localhost:8080/login

# Get statistics (with token)
curl -H "Authorization: token-..." http://localhost:8080/stats

# List all files
curl -H "Authorization: token-..." http://localhost:8080/list

# Search for music
curl -H "Authorization: token-..." "http://localhost:8080/search?q=artist_name"

# Stream audio
curl -H "Authorization: token-..." http://localhost:8080/stream?id=123
```

### Local Testing Scripts

**Quick Test**:
```bash
./scripts/quick-test.sh
```

**Full Test with Logs**:
```bash
./scripts/test-local.sh
```

**Cleanup**:
```bash
./scripts/cleanup.sh
```

## Performance Characteristics

### Scalability
- **Memory Usage**: Minimal (persistent database)
- **Concurrent Requests**: Limited by Go's HTTP server
- **File I/O**: Direct streaming, no buffering
- **Database**: Fast queries with indexed fields

### Limitations
- **Transcoding**: No audio format conversion
- **Caching**: No response caching
- **User Management**: Single user authentication
- **Playlists**: No playlist support yet

## Security Considerations

### Current Security Model
- **Token Authentication**: Secure token-based access control
- **File Access**: Direct file system access with read-only mounts
- **CORS**: Configurable cross-origin request handling
- **Input Validation**: SQL injection protection and parameter validation
- **Security Scanning**: Automated vulnerability detection

### Security Recommendations
- Use HTTPS in production (reverse proxy)
- Implement rate limiting
- Add user management for multi-user environments
- Regular security updates
- Monitor access logs

## Future Enhancement Opportunities

### High Priority
1. **User Management**: Multi-user support with roles
2. **Playlists**: Create and manage music playlists
3. **Caching**: Response and metadata caching
4. **Search Improvements**: Full-text search, fuzzy matching

### Medium Priority
1. **Metadata Editing**: Web interface for tag editing
2. **Statistics**: Detailed usage analytics
3. **Audio Transcoding**: Format conversion support
4. **Mobile App**: Native mobile client

### Low Priority
1. **Social Features**: Sharing and recommendations
2. **Cloud Integration**: Remote music sources
3. **Advanced Audio**: Equalizer, effects
4. **API Rate Limiting**: Request throttling

## Dependencies

### Core Dependencies
- `github.com/bogem/id3v2 v1.2.0`: Audio metadata parsing
- `github.com/gin-gonic/gin v1.10.1`: HTTP web framework
- `modernc.org/sqlite v1.38.0`: SQLite database driver

### Build Dependencies
- `golang.org/x/sys v0.33.0`: System calls
- `golang.org/x/text v0.15.0`: Text processing
- `modernc.org/libc v1.65.10`: C library bindings

## API Reference

### Response Formats

**Error Response**:
```json
{
    "error": "Error message"
}
```

**Music File Object**:
```json
{
    "id": "123",
    "path": "/app/music/artist/album/song.mp3",
    "title": "Song Title",
    "artist": "Artist Name",
    "album": "Album Name"
}
```

**Statistics Response**:
```json
{
    "total_files": 1250
}
```

**Login Response**:
```json
{
    "token": "token-..."
}
```

### HTTP Status Codes

- `200 OK`: Successful request
- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Authentication required or failed
- `404 Not Found`: File or resource not found
- `500 Internal Server Error`: Server error

### /login (POST)
- **Request**: `{ "username": "admin", "password": "password" }`
- **Response**: `{ "token": "..." }`
- **Error**: `401 Unauthorized` if credentials are invalid

### Authenticated Requests
- All other endpoints require:
  - Header: `Authorization: <token>`
- **Example**:
  ```bash
  curl -H "Authorization: token-..." http://localhost:8080/stats
  ```

## Monitoring and Logging

### Current Logging
- File scanning progress
- Metadata extraction errors
- Streaming requests
- Database errors
- Authentication attempts

### Recommended Monitoring
- Request metrics
- Error rates
- File processing statistics
- Performance metrics
- Security events

## Operating Instructions

See [docs/instructions.md](instructions.md) for setup, deployment, and usage instructions.

---

*This technical specification documents the current state of the Tower of Song application as of the latest implementation. Use this document as a foundation for future feature development and system enhancements.*
