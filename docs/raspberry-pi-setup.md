# Tower of Song: Raspberry Pi Home Streaming Setup Guide

This guide will walk you through installing, configuring, and exposing the Tower of Song music streaming server on a Raspberry Pi, making it accessible from anywhere.

---

## Prerequisites

- **Raspberry Pi 3/4** (recommended, running Raspberry Pi OS or similar)
- **MicroSD card** (16GB+ recommended)
- **Home network with router access**
- **Music files** on a USB drive or local storage
- **A computer for SSH/remote access**

---

## 1. Prepare Your Raspberry Pi

1. **Flash Raspberry Pi OS** to your SD card (use [Raspberry Pi Imager](https://www.raspberrypi.com/software/)).
2. **Boot and connect** your Pi to your home network (Ethernet recommended for stability).
3. **Update your system:**
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```
4. **Install Docker:**
   ```bash
   curl -sSL https://get.docker.com | sh
   sudo usermod -aG docker $USER
   # Log out and back in for group changes to take effect
   ```
5. **(Optional) Install Docker Compose:**
   ```bash
   sudo apt install -y docker-compose
   ```

---

## 2. Deploy Tower of Song

### Option A: Use Pre-built Docker Image (Recommended)

The easiest way to get started is using the pre-built Docker image from Docker Hub:

```bash
# Pull the latest image
docker pull punkscience/tower-of-song:latest

# Create directories for your music and data
mkdir -p ~/music ~/tower-data

# Run the server
docker run -d \
  --name tower-of-song \
  --restart unless-stopped \
  -p 8080:8080 \
  -v ~/music:/app/music:ro \
  -v ~/tower-data:/app/data \
  punkscience/tower-of-song:latest
```

### Option B: Build from Source

If you want to build the image locally or use a specific version:

```bash
# Clone the repository
git clone https://github.com/punkscience/tower-of-song.git
cd tower-of-song

# Build the Docker image
docker build -t tower-of-song .

# Run the server
docker run -d \
  --name tower-of-song \
  --restart unless-stopped \
  -p 8080:8080 \
  -v ~/music:/app/music:ro \
  -v ~/tower-data:/app/data \
  tower-of-song
```

---

## 3. Configure Your Music Library

1. **Copy your music files** to the `~/music` directory on your Pi
2. **Check the server logs** to see the scanning progress:
   ```bash
   docker logs -f tower-of-song
   ```
3. **Test the interface** by opening `http://<raspberry-pi-ip>:8080` in your browser
4. **Login** with the default credentials: `admin` / `password`

---

## 4. Make It Accessible from the Internet

### **A. Secure HTTPS Access with Caddy (Recommended for Custom Domains)**

1. **Get a domain name** (from any registrar).
2. **Point your domain's DNS A record to your home IP address.**
3. **Install Caddy:**
   ```bash
   sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
   curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo apt-key add -
   curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
   sudo apt update
   sudo apt install caddy
   ```
4. **Edit `/etc/caddy/Caddyfile`:**
   ```
   yourdomain.com {
       reverse_proxy localhost:8080
   }
   ```
5. **Restart Caddy:**
   ```bash
   sudo systemctl restart caddy
   ```
6. **Forward port 80/443 on your router to your Pi for HTTPS.**
7. **Access your server at `https://yourdomain.com` with a valid certificate.**

### **B. Secure Public Access with ngrok**

If you want a quick, secure tunnel without DNS or router config, use [ngrok](https://ngrok.com/):

1. **Sign up for a free account** at [ngrok.com](https://ngrok.com/).
2. **Download and install ngrok** on your Raspberry Pi:
   ```bash
   wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-arm.zip
   unzip ngrok-stable-linux-arm.zip
   sudo mv ngrok /usr/local/bin/
   ```
3. **Authenticate ngrok** with your authtoken (from your ngrok dashboard):
   ```bash
   ngrok config add-authtoken <YOUR_NGROK_AUTHTOKEN>
   ```
4. **Start your Docker container** as usual (see above).
5. **In a new terminal, run:**
   ```bash
   ngrok http 8080
   ```
6. **ngrok will display a public HTTPS URL** (e.g., `https://abcd1234.ngrok.io`).
7. **Open this URL in your browser** from anywhere in the world. It will proxy securely to your Pi with a valid certificate.

---

## 5. Management Commands

### **Update to Latest Version**
```bash
# Stop the current container
docker stop tower-of-song
docker rm tower-of-song

# Pull the latest image
docker pull punkscience/tower-of-song:latest

# Start with the new image
docker run -d \
  --name tower-of-song \
  --restart unless-stopped \
  -p 8080:8080 \
  -v ~/music:/app/music:ro \
  -v ~/tower-data:/app/data \
  punkscience/tower-of-song:latest
```

### **View Logs**
```bash
# Follow logs in real-time
docker logs -f tower-of-song

# View recent logs
docker logs --tail 50 tower-of-song
```

### **Stop the Service**
```bash
docker stop tower-of-song
```

### **Start the Service**
```bash
docker start tower-of-song
```

### **Restart the Service**
```bash
docker restart tower-of-song
```

---

## 6. Using the Service

- **Local Access**: Open `http://<raspberry-pi-ip>:8080` in your browser
- **Remote Access**: Use your domain (with Caddy) or ngrok URL
- **Login**: Use `admin` / `password` (change in production!)
- **Features**: Browse, search, and stream your music from anywhere

---

## 7. Security Recommendations

- **Change the default password** by editing the `config.json` file and rebuilding
- **Use HTTPS** for remote access (Caddy or ngrok)
- **Keep your Pi updated** regularly
- **Consider a VPN** for private access
- **Monitor access logs** for unusual activity

---

## 8. Troubleshooting

### **Can't access from outside?**
- Double-check port forwarding and your Pi's IP
- Make sure your ISP doesn't block incoming ports
- Verify the Docker container is running: `docker ps`

### **Music not found?**
- Check your music directory path and Docker volume mount
- Verify file permissions: `ls -la ~/music`
- Check container logs: `docker logs tower-of-song`

### **Container won't start?**
- Check available disk space: `df -h`
- Verify Docker is running: `sudo systemctl status docker`
- Check for port conflicts: `netstat -tlnp | grep 8080`

### **Performance issues?**
- Use Ethernet instead of WiFi for better stability
- Ensure adequate power supply for your Pi
- Monitor system resources: `htop`

---

## 9. Advanced Configuration

### **Custom Configuration**
Create a custom `config.json` file:
```json
{
    "music_folders": ["/app/music"],
    "username": "your-username",
    "password": "your-secure-password"
}
```

Mount it with the container:
```bash
docker run -d \
  --name tower-of-song \
  --restart unless-stopped \
  -p 8080:8080 \
  -v ~/music:/app/music:ro \
  -v ~/tower-data:/app/data \
  -v ~/config.json:/app/config.json:ro \
  punkscience/tower-of-song:latest
```

### **Multiple Music Directories**
Mount multiple directories:
```bash
docker run -d \
  --name tower-of-song \
  --restart unless-stopped \
  -p 8080:8080 \
  -v ~/music1:/app/music/music1:ro \
  -v ~/music2:/app/music/music2:ro \
  -v ~/tower-data:/app/data \
  punkscience/tower-of-song:latest
```

---

## 10. Support and Updates

- **GitHub Repository**: [https://github.com/punkscience/tower-of-song](https://github.com/punkscience/tower-of-song)
- **Docker Hub**: [https://hub.docker.com/r/punkscience/tower-of-song](https://hub.docker.com/r/punkscience/tower-of-song)
- **Issues**: Report bugs and feature requests on GitHub
- **Updates**: Pull the latest image regularly for security updates

---

*This guide covers the basic setup for Tower of Song on Raspberry Pi. For advanced configuration and troubleshooting, see the main documentation.* 