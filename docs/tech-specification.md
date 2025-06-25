# Tower of Song - Technical Specification

## Overview

Tower of Song is a lightweight, self-hosted music streaming server written in Go. It provides a RESTful API for browsing, searching, and streaming music files from local directories. The application automatically scans configured music folders, extracts metadata from audio files, and serves them through a web interface.

## Architecture

### Technology Stack

- **Backend**: Go 1.22.4
- **Database**: SQLite (in-memory with WAL mode)
- **Audio Metadata**: ID3v2 tags via `github.com/bogem/id3v2`
- **Containerization**: Docker with multi-stage build
- **Frontend**: HTML/JavaScript test client

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Music Files   │    │   Tower of      │    │   Web Client    │
│   (Local FS)    │◄──►│   Song Server   │◄──►│   (Browser)     │
│                 │    │   (Go)          │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   SQLite DB     │
                       │   (In-Memory)   │
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
- **Response**: JSON with total file count
- **Example**: `{"total_files": 1250}`

#### `/list` - List All Music
- **Method**: GET
- **Response**: JSON array of music files
- **Sorting**: By artist ASC, then title ASC
- **Fields**: id, path, title, artist, album

#### `/search` - Search Music Library
- **Method**: GET
- **Parameters**: `q` (search query)
- **Search Scope**: title, artist, album, file path
- **Matching**: Case-insensitive LIKE queries
- **Response**: JSON array of matching files

#### `/stream` - Stream Audio File
- **Method**: GET
- **Parameters**: `id` (file ID from database)
- **Response**: Audio stream with `audio/mpeg` content type
- **Features**: Direct file streaming, no transcoding

#### `/trackinfo` - Get Track Details
- **Method**: GET
- **Parameters**: `id` (track id)
- **Authentication**: Required
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

### 3. Cross-Origin Resource Sharing (CORS)

All API endpoints support CORS with the following headers:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type`

### 4. Configuration Management

**Configuration File**: `config.json`
```json
{
    "music_folders": ["/app/music"],
    "username": "yourusername",
    "password": "yourpassword"
}
```

**Features**:
- JSON-based configuration
- Support for multiple music directories
- Runtime configuration loading
- Docker volume mounting support

## Authentication and Security

### Authentication Model

- All API endpoints (except `/login`) require authentication via a token.
- Credentials (username and password) are now read from `config.json`:
  ```json
  {
    "music_folders": ["/app/music"],
    "username": "yourusername",
    "password": "yourpassword"
  }
  ```
- Clients must POST to `/login` with a JSON body containing `username` and `password`.
- On successful authentication, the server returns a token.
- The client must include this token in the `Authorization` header for all subsequent API requests.
- If the token is missing or invalid, the server responds with `401 Unauthorized`.
- The server uses an in-memory token store (tokens are lost on restart).

### API Changes

- All endpoints (`/stats`, `/list`, `/search`, `/stream`) now require a valid token.
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
- All requests to `/stats`, `/list`, `/search`, `/stream` must include:
  - Header: `Authorization: <token>`

#### Error Responses
- `401 Unauthorized` if token is missing or invalid

## Test Client Changes

- The HTML test client now includes a login form.
- All controls are disabled until the user logs in successfully.
- The token is stored in JavaScript and sent with all API requests in the `Authorization` header.
- Streaming audio also requires authentication.

## Implementation Details

### Database Management

- **Type**: In-memory SQLite database
- **Mode**: WAL (Write-Ahead Logging) enabled for better concurrency
- **Persistence**: Data is lost on server restart
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

## Deployment

### Docker Deployment

**Multi-stage Build**:
1. **Builder Stage**: Go compilation
2. **Runtime Stage**: Minimal Ubuntu image

**Dockerfile Features**:
- Optimized image size
- Security certificates included
- SQLite3 runtime dependencies
- Port 8080 exposed

**Usage**:
```bash
docker build -t tower-of-song .
docker run -p 8080:8080 -v /path/to/music:/app/music tower-of-song
```

**To stop the Docker server:**
```bash
docker ps                # Find the container ID or name
docker stop <container>  # Replace <container> with the ID or name
```

### Local Development

**Prerequisites**:
- Go 1.22.4+
- SQLite3 (for development)

**Build and Run**:
```bash
go mod download
go build -o tower-of-song
./tower-of-song
```

## Testing

### Test Client

Located in `tests/test.html`, provides a simple web interface for:
- Testing all API endpoints
- Audio streaming functionality
- JSON response visualization
- Interactive search capabilities

### API Testing

**Manual Testing**:
```bash
# Get statistics
curl http://localhost:8080/stats

# List all files
curl http://localhost:8080/list

# Search for music
curl "http://localhost:8080/search?q=artist_name"

# Stream audio
curl http://localhost:8080/stream?id=123
```

## Performance Characteristics

### Scalability
- **Memory Usage**: Minimal (in-memory database)
- **Concurrent Requests**: Limited by Go's HTTP server
- **File I/O**: Direct streaming, no buffering
- **Database**: Fast queries with indexed fields

### Limitations
- **Persistence**: Data lost on restart
- **Transcoding**: No audio format conversion
- **Caching**: No response caching
- **Authentication**: No user management

## Security Considerations

### Current Security Model
- **No Authentication**: Open access to all endpoints
- **File Access**: Direct file system access
- **CORS**: Permissive (allows all origins)
- **Input Validation**: Basic SQL injection protection

### Security Recommendations
- Implement authentication/authorization
- Add input sanitization
- Restrict CORS origins
- Implement rate limiting
- Add HTTPS support

## Future Enhancement Opportunities

### High Priority
1. **Persistent Database**: File-based SQLite storage
2. **Authentication**: User management system
3. **Playlists**: Create and manage music playlists

### Medium Priority
1. **Caching**: Response and metadata caching
2. **Search Improvements**: Full-text search, fuzzy matching
3. **Metadata Editing**: Web interface for tag editing
4. **Statistics**: Detailed usage analytics

### Low Priority
1. **Mobile App**: Native mobile client
2. **Social Features**: Sharing and recommendations
3. **Cloud Integration**: Remote music sources
4. **Advanced Audio**: Equalizer, effects

## Dependencies

### Core Dependencies
- `github.com/bogem/id3v2 v1.2.0`: Audio metadata parsing
- `modernc.org/sqlite v1.35.0`: SQLite database driver

### Build Dependencies
- `golang.org/x/sys v0.28.0`: System calls
- `golang.org/x/text v0.3.2`: Text processing
- `modernc.org/libc v1.61.13`: C library bindings

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

### HTTP Status Codes

- `200 OK`: Successful request
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

### Recommended Monitoring
- Request metrics
- Error rates
- File processing statistics
- Performance metrics

## Operating Instructions

See [docs/instructions.md](instructions.md) for setup, deployment, and usage instructions.

---

*This technical specification documents the current state of the Tower of Song application as of the latest implementation. Use this document as a foundation for future feature development and system enhancements.*
