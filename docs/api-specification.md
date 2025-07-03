# Tower of Song - API Specification

## Overview

The Tower of Song API provides secure access to your music collection through a RESTful interface. All endpoints (except `/login`) require authentication via a token-based system.

## Quick Start

### 1. Authentication

First, obtain an access token by authenticating with your credentials:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

**Response:**
```json
{
  "token": "token-1703123456789012345"
}
```

### 2. Using the Token

Include the token in the `Authorization` header for all subsequent requests:

```bash
curl -H "Authorization: token-1703123456789012345" \
  http://localhost:8080/stats
```

## API Endpoints

### Authentication

#### POST /login

Authenticate user and obtain access token.

**Request Body:**
```json
{
  "username": "admin",
  "password": "password"
}
```

**Response:**
```json
{
  "token": "token-1703123456789012345"
}
```

**Status Codes:**
- `200` - Authentication successful
- `400` - Invalid request format
- `401` - Invalid credentials
- `405` - Method not allowed

### Library Management

#### GET /stats

Get music library statistics.

**Headers:**
```
Authorization: token-1703123456789012345
```

**Response:**
```json
{
  "total_files": 1250
}
```

**Status Codes:**
- `200` - Statistics retrieved successfully
- `401` - Unauthorized

#### GET /list

List all music files in the library, sorted by artist and title.

**Headers:**
```
Authorization: token-1703123456789012345
```

**Response:**
```json
[
  {
    "id": "1",
    "path": "/app/music/Artist1/Album1/song1.mp3",
    "title": "Song Title 1",
    "artist": "Artist 1",
    "album": "Album 1"
  },
  {
    "id": "2",
    "path": "/app/music/Artist2/Album2/song2.mp3",
    "title": "Song Title 2",
    "artist": "Artist 2",
    "album": "Album 2"
  }
]
```

**Status Codes:**
- `200` - Files retrieved successfully
- `401` - Unauthorized

### Search

#### GET /search

Search for music files by title, artist, album, or file path.

**Headers:**
```
Authorization: token-1703123456789012345
```

**Query Parameters:**
- `q` (required) - Search query string

**Example:**
```bash
curl -H "Authorization: token-1703123456789012345" \
  "http://localhost:8080/search?q=artist%20name"
```

**Response:**
```json
[
  {
    "id": "1",
    "path": "/app/music/Artist1/Album1/song1.mp3",
    "title": "Song Title 1",
    "artist": "Artist 1",
    "album": "Album 1"
  }
]
```

**Status Codes:**
- `200` - Search results retrieved successfully
- `400` - Missing search query parameter
- `401` - Unauthorized

### Track Information

#### GET /trackinfo

Get detailed information about a specific track by its ID.

**Headers:**
```
Authorization: token-1703123456789012345
```

**Query Parameters:**
- `id` (required) - Track ID

**Example:**
```bash
curl -H "Authorization: token-1703123456789012345" \
  "http://localhost:8080/trackinfo?id=123"
```

**Response:**
```json
{
  "id": "123",
  "path": "/app/music/Artist1/Album1/song1.mp3",
  "title": "Song Title",
  "artist": "Artist Name",
  "album": "Album Name"
}
```

**Status Codes:**
- `200` - Track information retrieved successfully
- `400` - Missing or invalid track ID
- `401` - Unauthorized
- `404` - Track not found

### Streaming

#### GET /stream

Stream an audio file by its ID. Returns the audio data directly.

**Headers:**
```
Authorization: token-1703123456789012345
```

**Query Parameters:**
- `id` (required) - Track ID to stream
- `token` (optional) - Authentication token (alternative to Authorization header)

**Example:**
```bash
# Using Authorization header
curl -H "Authorization: token-1703123456789012345" \
  "http://localhost:8080/stream?id=123" \
  --output song.mp3

# Using token query parameter
curl "http://localhost:8080/stream?id=123&token=token-1703123456789012345" \
  --output song.mp3
```

**Response:**
- `200` - Audio file data (binary)
- Content-Type: `audio/mpeg`

**Status Codes:**
- `200` - Audio file streamed successfully
- `401` - Unauthorized
- `404` - Track not found

## Data Models

### MusicFile

```json
{
  "id": "string",
  "path": "string",
  "title": "string",
  "artist": "string",
  "album": "string"
}
```

**Properties:**
- `id` - Unique identifier for the music file
- `path` - File system path to the music file
- `title` - Song title (from ID3v2 tags or filename)
- `artist` - Artist name (from ID3v2 tags)
- `album` - Album name (from ID3v2 tags)

### Error

```json
{
  "error": "string"
}
```

**Properties:**
- `error` - Error message

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters or request format |
| 401 | Unauthorized - Invalid or missing authentication token |
| 404 | Not Found - Resource not found |
| 405 | Method Not Allowed - HTTP method not supported |
| 500 | Internal Server Error - Server error |

## Authentication

### Token Format

Tokens are generated with the format: `token-<timestamp>`

Example: `token-1703123456789012345`

### Token Usage

1. **Authorization Header (Recommended):**
   ```
   Authorization: token-1703123456789012345
   ```

2. **Query Parameter (for streaming):**
   ```
   ?token=token-1703123456789012345
   ```

### Token Security

- Tokens are stored in memory and lost on server restart
- No expiration time (tokens remain valid until server restart)
- Use HTTPS in production for secure token transmission

## CORS Support

All endpoints support Cross-Origin Resource Sharing (CORS) with the following headers:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, OPTIONS, POST
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Supported Audio Formats

- MP3 (`.mp3`)
- FLAC (`.flac`)
- WAV (`.wav`)

## Complete Examples

### JavaScript/Node.js

```javascript
const axios = require('axios');

class TowerOfSongAPI {
  constructor(baseURL = 'http://localhost:8080') {
    this.baseURL = baseURL;
    this.token = null;
  }

  async login(username, password) {
    const response = await axios.post(`${this.baseURL}/login`, {
      username,
      password
    });
    this.token = response.data.token;
    return this.token;
  }

  async getStats() {
    const response = await axios.get(`${this.baseURL}/stats`, {
      headers: { Authorization: this.token }
    });
    return response.data;
  }

  async listFiles() {
    const response = await axios.get(`${this.baseURL}/list`, {
      headers: { Authorization: this.token }
    });
    return response.data;
  }

  async search(query) {
    const response = await axios.get(`${this.baseURL}/search`, {
      headers: { Authorization: this.token },
      params: { q: query }
    });
    return response.data;
  }

  async getTrackInfo(id) {
    const response = await axios.get(`${this.baseURL}/trackinfo`, {
      headers: { Authorization: this.token },
      params: { id }
    });
    return response.data;
  }

  getStreamURL(id) {
    return `${this.baseURL}/stream?id=${id}&token=${this.token}`;
  }
}

// Usage
const api = new TowerOfSongAPI();
await api.login('admin', 'password');
const stats = await api.getStats();
console.log(`Total files: ${stats.total_files}`);
```

### Python

```python
import requests

class TowerOfSongAPI:
    def __init__(self, base_url='http://localhost:8080'):
        self.base_url = base_url
        self.token = None

    def login(self, username, password):
        response = requests.post(f'{self.base_url}/login', json={
            'username': username,
            'password': password
        })
        self.token = response.json()['token']
        return self.token

    def get_stats(self):
        response = requests.get(f'{self.base_url}/stats', 
                              headers={'Authorization': self.token})
        return response.json()

    def list_files(self):
        response = requests.get(f'{self.baseURL}/list', 
                              headers={'Authorization': self.token})
        return response.json()

    def search(self, query):
        response = requests.get(f'{self.baseURL}/search', 
                              headers={'Authorization': self.token},
                              params={'q': query})
        return response.json()

    def get_track_info(self, track_id):
        response = requests.get(f'{self.baseURL}/trackinfo', 
                              headers={'Authorization': self.token},
                              params={'id': track_id})
        return response.json()

    def get_stream_url(self, track_id):
        return f'{self.baseURL}/stream?id={track_id}&token={self.token}'

# Usage
api = TowerOfSongAPI()
api.login('admin', 'password')
stats = api.get_stats()
print(f"Total files: {stats['total_files']}")
```

### cURL Examples

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}' | \
  jq -r '.token')

# Get statistics
curl -H "Authorization: $TOKEN" http://localhost:8080/stats

# List all files
curl -H "Authorization: $TOKEN" http://localhost:8080/list

# Search for music
curl -H "Authorization: $TOKEN" "http://localhost:8080/search?q=artist%20name"

# Get track info
curl -H "Authorization: $TOKEN" "http://localhost:8080/trackinfo?id=123"

# Stream audio file
curl -H "Authorization: $TOKEN" "http://localhost:8080/stream?id=123" --output song.mp3
```

## Rate Limiting

Currently, no rate limiting is implemented. Consider implementing rate limiting for production deployments.

## Security Considerations

1. **Token Security**: Tokens are stored in memory and lost on server restart
2. **File Access**: Direct file system access to music files
3. **CORS**: Permissive CORS policy (allows all origins)
4. **Input Validation**: Basic SQL injection protection via parameterized queries

## Development Notes

- The API uses SQLite with WAL mode for better concurrency
- Music files are scanned every 24 hours automatically
- Missing files are automatically removed from the database during scans
- All database operations are protected by mutex for thread safety

## OpenAPI Specification

For machine-readable API documentation, see [api-specification.yaml](api-specification.yaml).

You can view this specification in Swagger UI or other OpenAPI tools by importing the YAML file.

---

*This API specification is based on Tower of Song version 1.0.0. For the latest updates, refer to the project documentation.* 