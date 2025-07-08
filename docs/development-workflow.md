# Development Workflow

## Branch Strategy

### Main Branch (`main`)
- **Protected branch** - requires pull request reviews
- **Automated builds** - triggers CI/CD pipeline on every push
- **Production ready** - only contains tested, reviewed code
- **Docker images** - automatically built and pushed to Docker Hub

### Development Branch (`dev`)
- **Working branch** - for active development
- **Feature branches** - create from `dev` for new features
- **Testing** - run tests locally before merging to `main`

## Workflow Process

1. **Start Development**
   ```bash
   git checkout dev
   git pull origin dev
   ```

2. **Create Feature Branch** (for new features)
   ```bash
   git checkout -b feature/your-feature-name
   # Make your changes
   git add .
   git commit -m "Add feature description"
   git push origin feature/your-feature-name
   # Create PR to merge into dev
   ```

3. **Direct Development** (for small changes)
   ```bash
   git checkout dev
   # Make your changes
   git add .
   git commit -m "Description of changes"
   git push origin dev
   ```

4. **Merge to Main** (when ready for production)
   ```bash
   # Create PR from dev to main
   # CI/CD will automatically:
   # - Run tests
   # - Security scan
   # - Build Docker image
   # - Push to Docker Hub
   ```

## CI/CD Pipeline

### What Happens on Main Branch Push:
1. **Testing Phase**
   - Go module download and caching
   - Unit tests execution
   - Security vulnerability scanning

2. **Build Phase** (only if tests pass)
   - Multi-platform Docker image build (ARM64, AMD64)
   - Image tagging with version metadata
   - Push to Docker Hub registry

### Docker Image Tags:
- `latest` - Latest successful build
- `main` - Latest main branch build
- `sha-{commit}` - Specific commit builds
- `v1.0.0` - Semantic version tags

## Raspberry Pi Deployment

Once the CI/CD pipeline completes successfully, you can update your Raspberry Pi:

```bash
# On your Raspberry Pi
docker pull punkscience/tower-of-song:latest
docker stop tower-of-song
docker rm tower-of-song
docker run -d --name tower-of-song -p 8080:8080 -v /path/to/music:/app/music -v /path/to/data:/app/data punkscience/tower-of-song:latest
```

## Required GitHub Secrets

Configure these secrets in your GitHub repository settings:

- `DOCKER_USERNAME` - Your Docker Hub username
- `DOCKER_PASSWORD` - Your Docker Hub access token

## Local Development

### Prerequisites
- Go 1.21+
- Docker
- Git

### Running Locally
```bash
go mod download
go run main.go
```

### Testing
```bash
go test -v ./...
```

### Building Docker Image Locally
```bash
docker build -t tower-of-song .
docker run -p 8080:8080 -v /path/to/music:/app/music -v /path/to/data:/app/data tower-of-song
``` 