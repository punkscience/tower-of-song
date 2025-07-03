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

## 2. Get Tower of Song

1. **Clone the repository:**
   ```bash
   git clone https://github.com/punkscience/tower-of-song.git
   cd tower-of-song
   ```
2. **Plug in your USB drive** (if using external storage) and note its mount path (e.g., `/media/pi/MUSIC`).

---

## 3. Build and Run with Docker

1. **Build the Docker image:**
   ```bash
   docker build -t tower-of-song .
   ```
2. **Run the server:**
   ```bash
   docker run -d \
     --name tower-of-song \
     -p 8080:8080 \
     -v /path/to/your/music:/app/music:ro \
     -v $(pwd)/config.json:/app/config.json:ro \
     tower-of-song
   ```
   - The `:ro` makes the music and config read-only inside the container.

3. **Test locally:**
   - Open `http://<raspberry-pi-ip>:8080` in your browser.
   - Use the credentials `admin` / `password` to log in.

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

### **B. (Alternative) Secure Public Access with ngrok**

If you want a quick, secure tunnel without DNS or router config, use [ngrok](https://ngrok.com/):

1. Sign up for a free account at [ngrok.com](https://ngrok.com/).
2. Download and install ngrok on your Raspberry Pi:
   ```bash
   wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-arm.zip
   unzip ngrok-stable-linux-arm.zip
   sudo mv ngrok /usr/local/bin/
   ```
3. Authenticate ngrok with your authtoken (from your ngrok dashboard):
   ```bash
   ngrok config add-authtoken <YOUR_NGROK_AUTHTOKEN>
   ```
4. Start your Docker container as usual (see above).
5. In a new terminal, run:
   ```bash
   ngrok http 8080
   ```
6. ngrok will display a public HTTPS URL (e.g., `https://abcd1234.ngrok.io`).
7. Open this URL in your browser from anywhere in the world. It will proxy securely to your Pi with a valid certificate.

---

## 5. Run on Boot (Optional)

Docker containers started with `--restart unless-stopped` will auto-start on boot:
```bash
# Stop and remove if already running
sudo docker rm -f tower-of-song
# Start with restart policy
sudo docker run -d \
  --name tower-of-song \
  --restart unless-stopped \
  -p 8080:8080 \
  -v /path/to/your/music:/app/music:ro \
  -v $(pwd)/config.json:/app/config.json:ro \
  tower-of-song
```

---

## 6. (Optional) Secure Public Access with ngrok

If you want to expose your music server securely to the internet with HTTPS and a valid certificate, you can use [ngrok](https://ngrok.com/):

### **A. Install ngrok**

1. Sign up for a free account at [ngrok.com](https://ngrok.com/).
2. Download and install ngrok on your Raspberry Pi:
   ```bash
   wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-arm.zip
   unzip ngrok-stable-linux-arm.zip
   sudo mv ngrok /usr/local/bin/
   ```
3. Authenticate ngrok with your authtoken (from your ngrok dashboard):
   ```bash
   ngrok config add-authtoken <YOUR_NGROK_AUTHTOKEN>
   ```

### **B. Start a Secure Tunnel**

1. Start your Docker container as usual (see above).
2. In a new terminal, run:
   ```bash
   ngrok http 8080
   ```
3. ngrok will display a public HTTPS URL (e.g., `https://abcd1234.ngrok.io`).
4. Open this URL in your browser from anywhere in the world. It will proxy securely to your Pi with a valid certificate.

### **C. Notes and Tips**
- The HTTPS URL will change each time unless you pay for a reserved domain.
- You can use the ngrok dashboard to monitor traffic and connections.
- This is a great way to share your music server securely without configuring your router or exposing your home IP.

---

## 7. Using the Service

- Open the test client at `http://<your-public-ip>:8080` or your DuckDNS domain.
- Log in with your credentials.
- Browse, search, and stream your music from anywhere!

---

## 8. Security Recommendations

- **Change the default password** in the source code and rebuild.
- **Do not expose without HTTPS** for sensitive use.
- **Consider a VPN** for private access.
- **Keep your Pi and Docker updated.**

---

## 9. Troubleshooting

- **Can't access from outside?**
  - Double-check port forwarding and your Pi's IP.
  - Make sure your ISP doesn't block incoming ports.
- **Music not found?**
  - Check your `config.json` path and Docker volume mount.
- **Performance issues?**
  - Use Ethernet instead of WiFi for best results.

---

## 10. Uninstalling

```bash
sudo docker rm -f tower-of-song
sudo docker rmi tower-of-song
```

---

Enjoy your private home music streaming service! 