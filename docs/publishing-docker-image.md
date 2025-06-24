# Publishing and Distributing the Tower of Song Docker Image

This guide explains how to make the Tower of Song Docker image available to anyone via a public Docker registry (such as Docker Hub), and how users can install and run it from the cloud.

---

## 1. Build the Docker Image Locally

If you haven't already, build the image:

```bash
docker build -t yourdockerhubusername/tower-of-song:latest .
```

---

## 2. Create a Docker Hub Account

1. Go to [https://hub.docker.com/](https://hub.docker.com/) and sign up for a free account.
2. (Optional) Create a new repository named `tower-of-song`.

---

## 3. Log In to Docker Hub from Your Terminal

```bash
docker login
```
- Enter your Docker Hub username and password when prompted.

---

## 4. Tag and Push the Image

Tag your image (if you didn't already):

```bash
docker tag tower-of-song yourdockerhubusername/tower-of-song:latest
```

Push the image to Docker Hub:

```bash
docker push yourdockerhubusername/tower-of-song:latest
```

- The image will now be available publicly at `docker.io/yourdockerhubusername/tower-of-song:latest`.

---

## 5. How Users Can Install and Run the Image

Any user can now pull and run the image from anywhere with Docker installed:

```bash
docker pull yourdockerhubusername/tower-of-song:latest
```

Run the container (mounting your music and config):

```bash
docker run -d \
  --name tower-of-song \
  -p 8080:8080 \
  -v /path/to/your/music:/app/music:ro \
  -v /path/to/config.json:/app/config.json:ro \
  yourdockerhubusername/tower-of-song:latest
```

---

## 6. Keeping the Image Updated

- After making changes, rebuild and push a new version:
  ```bash
  docker build -t yourdockerhubusername/tower-of-song:latest .
  docker push yourdockerhubusername/tower-of-song:latest
  ```
- Consider using tags for releases (e.g., `:v1.0.0`).

---

## 7. (Optional) Automate Builds with GitHub Actions

You can automate Docker image builds and pushes using GitHub Actions. See `.github/workflows/docker-ci.yml` in this repo for an example.

---

## 8. Share the Image

- Share the pull/run instructions in your README and documentation.
- Users can now install Tower of Song from the cloud with a single command!

---

**Example for users:**

```bash
docker run -d \
  --name tower-of-song \
  -p 8080:8080 \
  -v /path/to/your/music:/app/music:ro \
  -v /path/to/config.json:/app/config.json:ro \
  yourdockerhubusername/tower-of-song:latest
```

---

Enjoy sharing your music server with the world! 