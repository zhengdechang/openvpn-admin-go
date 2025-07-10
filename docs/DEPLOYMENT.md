# 🚀 Deployment Guide

This document describes the automated CI/CD pipeline and deployment options for OpenVPN Admin Go.

## 📦 Automated Build & Release

### GitHub Actions Workflow

The project includes an automated GitHub Actions workflow that:

1. **Triggers on**: PR merge to `main` branch
2. **Builds**: Go binaries for multiple platforms (Linux, Windows, macOS)
3. **Creates**: Docker image with both frontend and backend
4. **Publishes**: Docker image to Docker Hub
5. **Releases**: Binaries and release notes on GitHub

### Required Secrets

Configure these secrets in your GitHub repository settings:

```bash
DOCKERHUB_USERNAME=your-dockerhub-username
DOCKERHUB_TOKEN=your-dockerhub-access-token
```

### Workflow Features

- ✅ Multi-platform Go binary builds (Linux, Windows, macOS)
- ✅ Combined Docker image with frontend and backend
- ✅ Automatic versioning based on timestamp
- ✅ GitHub release with download links
- ✅ Docker Hub publishing with latest and versioned tags

## 🐳 Docker Deployment

### Quick Start

```bash
# Pull and run the latest image
docker run -d \
  --name openvpn-admin \
  --cap-add NET_ADMIN \
  --device /dev/net/tun \
  -p 8085:8085 \
  -p 3000:3000 \
  -p 1194:1194/udp \
  -v openvpn_data:/app/data \
  -v openvpn_logs:/app/logs \
  -e JWT_SECRET=your-secret-key \
  -e OPENVPN_SERVER_HOSTNAME=your-server-ip \
  your-dockerhub-username/openvpn-admin-go:latest
```

### Service Modes

The Docker image supports different service modes via the `SERVICE_MODE` environment variable:

#### 1. Full Stack Mode (Default)
```bash
docker run -e SERVICE_MODE=all your-dockerhub-username/openvpn-admin-go:latest
```
- Runs both backend (port 8085) and frontend (port 3000)
- Complete web interface with API

#### 2. Backend Only Mode
```bash
docker run -e SERVICE_MODE=backend your-dockerhub-username/openvpn-admin-go:latest
```
- Runs only the Go backend API
- Useful for microservice deployments

#### 3. Frontend Only Mode
```bash
docker run -e SERVICE_MODE=frontend your-dockerhub-username/openvpn-admin-go:latest
```
- Runs only the Next.js frontend
- Requires separate backend instance

### Docker Compose Deployment

Use the provided `docker-compose.production.yml`:

```bash
# Navigate to docker directory
cd docker

# Copy and configure environment
cp .env.docker.example .env
# Edit .env with your settings

# Start services
docker-compose -f docker-compose.production.yml up -d

# View logs
docker-compose -f docker-compose.production.yml logs -f

# Stop services
docker-compose -f docker-compose.production.yml down
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVICE_MODE` | Service mode: `all`, `backend`, `frontend` | `all` |
| `BACKEND_PORT` | Backend API port | `8085` |
| `FRONTEND_PORT` | Frontend UI port | `3000` |
| `JWT_SECRET` | JWT signing secret | Required |
| `OPENVPN_SERVER_HOSTNAME` | Server IP/domain | Required |
| `OPENVPN_PORT` | OpenVPN port | `1194` |
| `OPENVPN_PROTO` | Protocol (udp/tcp) | `udp` |
| `DB_PATH` | Database file path | `/app/data/db.sqlite3` |

## 🔧 Manual Binary Deployment

### Download Binaries

Download from GitHub releases:
- Linux: `openvpn-admin-linux-amd64.tar.gz`
- Windows: `openvpn-admin-windows-amd64.zip`
- macOS: `openvpn-admin-darwin-amd64.tar.gz`

### Installation

```bash
# Linux/macOS
tar -xzf openvpn-admin-linux-amd64.tar.gz
chmod +x openvpn-admin-linux-amd64
./openvpn-admin-linux-amd64

# Windows
# Extract openvpn-admin-windows-amd64.zip
# Run openvpn-admin-windows-amd64.exe
```

## 🌐 Production Setup with Nginx

### Nginx Configuration

Use the provided `nginx.conf` for reverse proxy setup:

```bash
# Copy nginx configuration from frontend directory
cp openvpn-web/nginx.conf /etc/nginx/sites-available/openvpn-admin
ln -s /etc/nginx/sites-available/openvpn-admin /etc/nginx/sites-enabled/

# Test and reload nginx
nginx -t
systemctl reload nginx
```

### SSL/HTTPS Setup

1. Obtain SSL certificates (Let's Encrypt recommended)
2. Update nginx.conf with SSL configuration
3. Redirect HTTP to HTTPS

## 📊 Monitoring & Logging

### Health Checks

- Backend: `http://localhost:8085/api/health`
- Frontend: `http://localhost:3000`
- Combined: Both endpoints available

### Log Locations

- Application logs: `/app/logs/`
- OpenVPN logs: `/var/log/openvpn/`
- Nginx logs: `/var/log/nginx/`

### Docker Logs

```bash
# View all logs
docker-compose logs -f

# View specific service
docker-compose logs -f openvpn-admin

# View last 100 lines
docker-compose logs --tail=100 openvpn-admin
```

## 🔒 Security Considerations

1. **Change default JWT secret**
2. **Use strong passwords**
3. **Enable firewall rules**
4. **Regular security updates**
5. **Monitor access logs**
6. **Use HTTPS in production**

## 🚨 Troubleshooting

### Common Issues

1. **Permission denied for /dev/net/tun**
   - Ensure `--cap-add NET_ADMIN` and `--device /dev/net/tun`

2. **Database locked errors**
   - Check file permissions on data volume

3. **Frontend can't connect to backend**
   - Verify `NEXT_PUBLIC_API_BASE_URL` environment variable

4. **OpenVPN service fails to start**
   - Check privileged mode and capabilities
   - Verify OpenVPN configuration

### Debug Commands

```bash
# Check container status
docker ps -a

# Enter container for debugging
docker exec -it openvpn-admin /bin/bash

# Check service logs
docker logs openvpn-admin

# Test network connectivity
docker exec openvpn-admin curl -f http://localhost:8085/api/health
```
