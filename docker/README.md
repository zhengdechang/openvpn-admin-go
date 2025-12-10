# 🐳 Docker Deployment

This directory contains all Docker-related files for OpenVPN Admin Go deployment.

## 📁 Files Overview

- `Dockerfile.combined` - Multi-stage Dockerfile that builds both frontend and backend into a single image
- `docker-compose.production.yml` - Production-ready Docker Compose configuration
- `.env.docker.example` - Environment variables template for Docker deployment

## 🚀 Quick Start

### 1. Using Pre-built Image

```bash
# Pull and run the latest image (CLI menu by default)
docker run -it --rm \
  --name openvpn-admin \
  --cap-add NET_ADMIN \
  --device /dev/net/tun \
  -p 8085:8085 \
  -p 1194:1194/udp \
  -v openvpn_data:/app/data \
  -v openvpn_logs:/app/logs \
  -e JWT_SECRET=your-secret-key \
  -e OPENVPN_SERVER_HOSTNAME=your-server-ip \
  zhengdechang/openvpn-admin-go:latest

# Start directly into the built-in web service (same image)
docker run -d \
  --name openvpn-admin-web \
  --cap-add NET_ADMIN \
  --device /dev/net/tun \
  -p 8085:8085 \
  -v openvpn_data:/app/data \
  -v openvpn_logs:/app/logs \
  -e ENABLE_WEB=true \
  -e WEB_PORT=8085 \
  -e JWT_SECRET=your-secret-key \
  -e OPENVPN_SERVER_HOSTNAME=your-server-ip \
  zhengdechang/openvpn-admin-go:latest
```

### 2. Using Docker Compose

```bash
# Copy environment configuration
cp .env.docker.example .env
# Edit .env with your settings

# Start services
docker-compose -f docker-compose.production.yml up -d

# View logs
docker-compose -f docker-compose.production.yml logs -f

# Stop services
docker-compose -f docker-compose.production.yml down
```

## 🔧 Service Architecture

The application is deployed as two separate services using the same Docker image:

### Backend Service (openvpn-backend)
- Runs the Go API server on port 8085
- Handles OpenVPN management operations
- Requires privileged mode for OpenVPN functionality
- Uses host networking for proper OpenVPN operation

### Frontend Service (openvpn-frontend)
- Runs nginx web server on port 80
- Serves static frontend files from `/app/frontend/out`
- Proxies API requests to the backend service
- Provides the web interface for users

### Service Communication
```bash
docker run -e SERVICE_MODE=backend zhengdechang/openvpn-admin-go:latest
```
- Runs only the Go backend API
- Useful for microservice deployments

### Frontend Only Mode
```bash
docker run -e SERVICE_MODE=frontend zhengdechang/openvpn-admin-go:latest
```
- Runs only the Next.js frontend
- Requires separate backend instance

## 🌐 Nginx Reverse Proxy

The frontend service uses nginx to:

- Serve static frontend files from `/app/frontend/out`
- Proxy `/api/*` requests to the backend service
- Provide rate limiting and security headers
- Handle static file caching

Nginx configuration is located at `docker/nginx/nginx.conf`.

## 📊 Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `OPENVPN_SERVER_HOSTNAME` | Server IP/domain | - | Yes |
| `OPENVPN_PORT` | OpenVPN port | `1194` | No |
| `OPENVPN_PROTO` | Protocol (udp/tcp) | `udp` | No |
| `OPENVPN_NETWORK` | VPN network | `10.8.0.0` | No |
| `OPENVPN_NETMASK` | VPN netmask | `255.255.255.0` | No |
| `DB_PATH` | Database file path | `/app/data/db.sqlite3` | No |
| `TZ` | Timezone | `UTC` | No |

## 🌐 Ports

- **Port 80** - Frontend web interface (nginx)
- **Port 8085** - Backend API (direct access, optional)
- **Port 1194/udp** - OpenVPN server

## 🔒 Security Considerations

1. **Change default JWT secret** - Use a strong, random secret
2. **Use strong passwords** - Change default admin password
3. **Network security** - Use host network mode or proper port mapping
4. **File permissions** - Ensure proper volume permissions
5. **Regular updates** - Keep images updated

## 📝 Local Directory Mapping

All data is mapped to local directories for easy access and backup:

- `./data` → `/app/data` - Application database and data files
- `./logs` → `/app/logs` - Application log files
- `./config` → `/app/config` - Configuration files
- `./openvpn/etc` → `/etc/openvpn` - OpenVPN configuration directory
- `./openvpn/logs` → `/var/log/openvpn` - OpenVPN log files

These directories will be created automatically when you start the services.

## 🚨 Troubleshooting

### Common Issues

1. **Permission denied for /dev/net/tun**
   ```bash
   # Ensure proper capabilities and device access
   --cap-add NET_ADMIN --device /dev/net/tun
   ```

2. **Database locked errors**
   ```bash
   # Check volume permissions
   docker exec -it openvpn-admin ls -la /app/data/
   ```

3. **Frontend can't connect to backend**
   ```bash
   # Verify API URL environment variable
   docker exec -it openvpn-admin env | grep API_BASE_URL
   ```

### Debug Commands

```bash
# Check container status
docker ps -a

# View container logs
docker logs openvpn-admin

# Enter container for debugging
docker exec -it openvpn-admin /bin/bash

# Test health endpoints (only when ENABLE_WEB=true)
curl http://localhost:8085/api/health
curl http://localhost:3000
```

## 🔄 Building from Source

```bash
# Build the image locally
docker build -f Dockerfile.combined -t openvpn-admin-go:local ..

# Run the locally built image
docker run -d --name openvpn-admin-local openvpn-admin-go:local
```

## 📈 Monitoring

### Health Checks

- Backend: `http://localhost:8085/api/health` (only when `ENABLE_WEB=true`)
- Frontend: `http://localhost:3000` (served from the same image; start the web UI explicitly)
- Combined: Both endpoints available when the web UI is enabled

### Log Locations

- Application logs: `/app/logs/`
- OpenVPN logs: `/var/log/openvpn/`
- Container logs: `docker logs <container-name>`
