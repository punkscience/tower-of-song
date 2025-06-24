# Tower of Song: Operating Instructions

This document explains how to set up, run, and use Tower of Song.

---

## 1. Docker Deployment

**Build the Docker image:**
```bash
docker build -t tower-of-song .
```

**Run the server:**
```bash
docker run -p 8080:8080 -v /path/to/music:/app/music -v $(pwd)/config.json:/app/config.json tower-of-song
```

**Stop the server:**
```bash
docker ps                # Find the container ID or name
docker stop <container>  # Replace <container> with the ID or name
```

---

## 2. Local Development

**Requirements:**
- Go 1.22+
- SQLite3 (for development)

**Build and run:**
```bash
go mod download
go build -o tower-of-song
./tower-of-song
```

---

## 3. Configuration

Edit `config.json` to specify your music folders and credentials:
```json
{
    "music_folders": ["/app/music"],
    "username": "yourusername",
    "password": "yourpassword"
}
```

---

## 4. Using the Service

- Open the test client at `http://localhost:8080` (or your server's IP/domain).
- Log in with your credentials.
- Browse, search, and stream your music from anywhere!

---

## 5. More Information

- For API details, architecture, and security, see [docs/tech-specification.md](tech-specification.md).
- For Raspberry Pi setup, see [docs/raspberry-pi-setup.md](raspberry-pi-setup.md).
- For publishing the Docker image, see [docs/publishing-docker-image.md](publishing-docker-image.md). 