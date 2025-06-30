# Troubleshooting Guide

## Common Issues and Solutions

### 1. Connection Reset / Failed to Fetch

**Symptoms:**
- `net::ERR_CONNECTION_RESET` errors
- `Failed to fetch` JavaScript errors
- Login fails with connection issues

**Possible Causes and Solutions:**

#### A. Docker Container Not Running
```bash
# Check if container is running
docker ps

# If not running, check logs
docker logs tower-of-song

# Restart the container
docker restart tower-of-song
```

#### B. Port Not Exposed Correctly
```bash
# Check if port 8080 is exposed
docker port tower-of-song

# Should show: 0.0.0.0:8080->8080/tcp
```

#### C. Firewall Issues
```bash
# Check if port 8080 is open on Raspberry Pi
sudo netstat -tlnp | grep 8080

# If using ufw firewall, allow port 8080
sudo ufw allow 8080
```

#### D. Network Configuration
```bash
# Test if server responds locally on Pi
curl http://localhost:8080/stats

# Test from another machine on network
curl http://<raspberry-pi-ip>:8080/stats
```

### 2. 404 Errors on Stream Endpoint

**Symptoms:**
- Stream requests return 404 Not Found
- Audio player shows "no supported source found"

**Possible Causes and Solutions:**

#### A. Music Files Not Mounted Correctly
```bash
# Check if music volume is mounted
docker exec tower-of-song ls -la /app/music

# Check if config.json is correct
docker exec tower-of-song cat /app/config.json
```

#### B. File Permissions
```bash
# Check file permissions on music directory
ls -la /path/to/your/music

# Fix permissions if needed
chmod -R 755 /path/to/your/music
```

#### C. Music Files Not Scanned
```bash
# Check if files are in database
docker exec tower-of-song sqlite3 :memory: "SELECT COUNT(*) FROM music;"

# Force rescan by restarting container
docker restart tower-of-song
```

### 3. Authentication Issues

**Symptoms:**
- Login fails even with correct credentials
- 401 Unauthorized errors

**Solutions:**

#### A. Check Config File
```bash
# Verify config.json has correct credentials
docker exec tower-of-song cat /app/config.json
```

#### B. Rebuild with Updated Config
```bash
# Stop container
docker stop tower-of-song

# Update config.json with correct credentials
# Then restart
docker start tower-of-song
```

### 4. CORS Issues

**Symptoms:**
- Browser blocks requests due to CORS policy
- Preflight requests fail

**Solutions:**

#### A. Check CORS Headers
```bash
# Test CORS headers
curl -H "Origin: http://localhost" -H "Access-Control-Request-Method: GET" \
  -H "Access-Control-Request-Headers: Authorization" \
  -X OPTIONS http://<raspberry-pi-ip>:8080/stats
```

#### B. Verify CORS Configuration
The server should return:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, OPTIONS, POST`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

### 5. Debugging Steps

#### Step 1: Verify Container Status
```bash
docker ps -a
docker logs tower-of-song
```

#### Step 2: Test Local Connectivity
```bash
# On Raspberry Pi
curl http://localhost:8080/stats
```

#### Step 3: Test Network Connectivity
```bash
# From another machine
curl http://<raspberry-pi-ip>:8080/stats
```

#### Step 4: Check Music Files
```bash
# Verify music files are accessible
docker exec tower-of-song find /app/music -name "*.mp3" | head -5
```

#### Step 5: Test Authentication
```bash
# Test login endpoint
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' \
  http://<raspberry-pi-ip>:8080/login
```

### 6. Common Docker Commands

```bash
# View container logs
docker logs tower-of-song

# Execute commands in container
docker exec -it tower-of-song /bin/bash

# Restart container
docker restart tower-of-song

# Remove and recreate container
docker rm -f tower-of-song
docker run -d --name tower-of-song -p 8080:8080 \
  -v /path/to/music:/app/music:ro \
  -v $(pwd)/config.json:/app/config.json:ro \
  tower-of-song
```

### 7. Network Debugging

```bash
# Check if port is listening
netstat -tlnp | grep 8080

# Check firewall status
sudo ufw status

# Test port connectivity
telnet <raspberry-pi-ip> 8080
```

### 8. Input/Output Errors and Volume Mount Issues

**Symptoms:**
- `find: '/app/music': Input/output error`
- Connection reset errors
- 404 errors on stream endpoint
- File scanning fails mid-process

**Possible Causes and Solutions:**

#### A. Corrupted Volume Mount
```bash
# Check if the volume mount is working
docker exec tower-of-song ls -la /app/

# If this fails, the volume mount is corrupted
```

#### B. Filesystem Issues on Host
```bash
# Check filesystem health on Raspberry Pi
df -h
sudo fsck /dev/mmcblk0p2  # Only if filesystem is unmounted

# Check disk space
df -h /path/to/your/music
```

#### C. USB Drive Issues (if using external storage)
```bash
# Check USB drive health
sudo dmesg | grep -i usb
sudo dmesg | grep -i error

# Check if USB drive is properly mounted
mount | grep usb
```

#### D. File Permissions and Ownership
```bash
# Check permissions on host
ls -la /path/to/your/music

# Fix permissions if needed
sudo chown -R $USER:$USER /path/to/your/music
chmod -R 755 /path/to/your/music
```

#### E. Docker Volume Corruption
```bash
# Stop the container
docker stop tower-of-song

# Remove the container
docker rm tower-of-song

# Recreate with fresh volume mount
docker run -d --name tower-of-song -p 8080:8080 \
  -v /path/to/your/music:/app/music:ro \
  -v $(pwd)/config.json:/app/config.json:ro \
  tower-of-song
```

#### F. Alternative: Use Docker Volume Instead of Bind Mount
```bash
# Create a named volume
docker volume create music-data

# Copy music files to volume
docker run --rm -v music-data:/music -v /path/to/your/music:/source alpine cp -r /source/* /music/

# Run container with named volume
docker run -d --name tower-of-song -p 8080:8080 \
  -v music-data:/app/music:ro \
  -v $(pwd)/config.json:/app/config.json:ro \
  tower-of-song
```

#### G. Check for Hardware Issues
```bash
# Check system logs for hardware errors
sudo journalctl -f

# Check temperature (overheating can cause I/O errors)
vcgencmd measure_temp

# Check memory usage
free -h
```

#### H. Emergency Recovery Steps
```bash
# 1. Stop all containers
docker stop $(docker ps -q)

# 2. Restart Docker service
sudo systemctl restart docker

# 3. Check if filesystem is accessible
ls -la /path/to/your/music

# 4. If still failing, reboot the Pi
sudo reboot

# 5. After reboot, recreate container
docker run -d --name tower-of-song -p 8080:8080 \
  -v /path/to/your/music:/app/music:ro \
  -v $(pwd)/config.json:/app/config.json:ro \
  tower-of-song
```

---

**Immediate Action Plan:**

1. **Stop the container immediately:**
   ```bash
   docker stop tower-of-song
   ```

2. **Check your filesystem:**
   ```bash
   df -h
   ls -la /path/to/your/music
   ```

3. **If filesystem is accessible, try recreating the container:**
   ```bash
   docker rm tower-of-song
   docker run -d --name tower-of-song -p 8080:8080 \
     -v /path/to/your/music:/app/music:ro \
     -v $(pwd)/config.json:/app/config.json:ro \
     tower-of-song
   ```

4. **If still failing, consider using a different storage location or USB drive.**

---

**Prevention:**
- Use high-quality microSD cards for Raspberry Pi
- Avoid filling storage to >90% capacity
- Use external USB drives for large music collections
- Regularly backup your music files
- Monitor system temperature and power supply

---

If you're still experiencing issues, please provide:
1. Output of `docker logs tower-of-song`
2. Output of `docker ps`
3. Your `config.json` contents (without sensitive data)
4. Network configuration details 