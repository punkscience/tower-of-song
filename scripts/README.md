# Tower of Song - Test Scripts

This directory contains scripts for testing Tower of Song locally using Docker.

## Configuring Test Folders

You can specify your own music and data directories by creating a `test-folders.conf` file in this directory:

```
MUSIC_DIR=/absolute/path/to/your/music
DATA_DIR=/absolute/path/to/your/data
```

- Copy `test-folders.conf.example` to `test-folders.conf` and edit as needed.
- This file is ignored by git and will not be checked in.
- If not set, the scripts default to `./test-music` and `./test-data`.

## Scripts

### `test-local.sh` - Full Test Script
The most comprehensive test script with detailed output and automatic cleanup.

**Features:**
- Colored output for better readability
- Automatic cleanup on exit (Ctrl+C)
- Health checks to ensure server is running
- Detailed logging
- Error handling
- Reads test folders from `test-folders.conf` if present

**Usage:**
```bash
./scripts/test-local.sh
```

### `quick-test.sh` - Quick Test Script
A simpler script for faster testing without the full setup.

**Features:**
- Minimal output
- Fast execution
- Manual cleanup required
- Reads test folders from `test-folders.conf` if present

**Usage:**
```bash
./scripts/quick-test.sh
```

### `cleanup.sh` - Cleanup Script
Clean up test containers, images, and optionally test data.

**Usage:**
```bash
./scripts/cleanup.sh
```

## Test Setup

The scripts will create the following directory structure (unless overridden):

```
tower-of-song/
├── test-music/     # Mounted as /app/music (read-only)
├── test-data/      # Mounted as /app/data (persistent)
└── scripts/
    ├── test-local.sh
    ├── quick-test.sh
    └── cleanup.sh
```

## Adding Test Music

To test with actual music files:

1. Copy some MP3 files to the `test-music/` directory (or your configured `MUSIC_DIR`)
2. Run one of the test scripts
3. The server will scan and index the music files

## Default Login

- **Username:** `admin`
- **Password:** Check `config.json` for the configured password

## Accessing the Application

Once the server is running, open your browser and go to:
- **URL:** `http://localhost:8080`
- **API Documentation:** Available in the `/docs` directory

## Troubleshooting

### Port Already in Use
If port 8080 is already in use, you can modify the `PORT` variable in the script or stop the existing service.

### Permission Issues
The scripts use Docker volumes to avoid permission issues that occur when running the Go application directly.

### Container Won't Start
Check the container logs:
```bash
docker logs tower-of-song-test
```

### Clean Slate
To start completely fresh:
```bash
./scripts/cleanup.sh
# Then run your test script again
```

## Production vs Test

These scripts are for **testing only**. For production deployment:

1. Use the CI/CD pipeline to build and push to Docker Hub
2. Pull the image on your Raspberry Pi
3. Use proper volume mounts for your music library and data persistence 