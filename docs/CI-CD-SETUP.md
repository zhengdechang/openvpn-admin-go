# 🤖 CI/CD Setup Guide

This document provides a complete guide for setting up the automated CI/CD pipeline for OpenVPN Admin Go.

## 📋 Overview

The CI/CD pipeline automatically:
1. **Builds** multi-platform Go binaries (Linux, Windows, macOS)
2. **Creates** combined Docker image with frontend and backend
3. **Publishes** to Docker Hub with versioned tags
4. **Releases** binaries on GitHub with release notes

## 🔧 Prerequisites

### GitHub Repository Setup

1. **Fork or clone** the repository to your GitHub account
2. **Configure secrets** in repository settings (`Settings > Secrets and variables > Actions`)

### Required Secrets

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `DOCKERHUB_USERNAME` | Your Docker Hub username | `your-username` |
| `DOCKERHUB_TOKEN` | Docker Hub access token | `dckr_pat_...` |

### Docker Hub Token Setup

1. Go to [Docker Hub Account Settings](https://hub.docker.com/settings/security)
2. Click "New Access Token"
3. Name: `GitHub Actions`
4. Permissions: `Read, Write, Delete`
5. Copy the generated token to GitHub secrets

## 🚀 Workflow Triggers

The CI/CD pipeline triggers on:
- **Direct push** to `main` branch
- **Pull Request merge** to `main` branch

**Note**: Only merged PRs trigger the build, not open PRs.

## 📦 Build Process

### 1. Environment Setup
- Go 1.21 installation
- Node.js 18 with npm cache
- Docker Buildx for multi-platform builds

### 2. Frontend Build
```bash
cd openvpn-web
npm ci
npm run build
```

### 3. Backend Build
Multi-platform Go binaries:
- `openvpn-admin-linux-amd64`
- `openvpn-admin-windows-amd64.exe`
- `openvpn-admin-darwin-amd64`

### 4. Docker Image Build
- **Base**: Alpine Linux with runtime dependencies
- **Frontend**: Next.js production build
- **Backend**: Go binary
- **Services**: Configurable startup modes
- **Platforms**: linux/amd64, linux/arm64

### 5. Release Creation
- **Version**: Timestamp-based (e.g., `v20240101120000`)
- **Binaries**: Compressed archives for each platform
- **Docker**: Tagged with version and `latest`
- **Notes**: Auto-generated with usage instructions

## 🐳 Docker Image Features

### Service Modes
Control via `SERVICE_MODE` environment variable:

```bash
# Full stack (default)
SERVICE_MODE=all

# Backend only
SERVICE_MODE=backend

# Frontend only  
SERVICE_MODE=frontend
```

### Multi-Architecture Support
- **linux/amd64**: Standard x86_64 servers
- **linux/arm64**: ARM-based servers (e.g., Apple Silicon, Raspberry Pi)

### Health Checks
Built-in health monitoring:
- Backend: `http://localhost:8085/api/health`
- Frontend: `http://localhost:3000`
- Database connectivity check

## 📁 File Organization

```
openvpn-admin-go/
├── .github/workflows/
│   └── build-and-release.yml    # CI/CD pipeline
├── docker/
│   ├── Dockerfile.combined      # Multi-service Docker image
│   ├── docker-compose.production.yml
│   ├── .env.docker.example      # Environment template
│   └── README.md               # Docker deployment guide
├── openvpn-web/
│   ├── nginx.conf              # Reverse proxy configuration
│   └── ...                     # Frontend source code
├── scripts/
│   └── docker-entrypoint.sh    # Container startup script
└── docs/
    ├── DEPLOYMENT.md           # Deployment guide
    └── CI-CD-SETUP.md         # This file
```

## 🔄 Deployment Workflow

### 1. Development
```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes and commit
git add .
git commit -m "Add new feature"

# Push to GitHub
git push origin feature/new-feature
```

### 2. Pull Request
```bash
# Create PR on GitHub
# Review and approve
# Merge to main branch
```

### 3. Automatic Build
- GitHub Actions triggers automatically
- Builds and tests all components
- Creates Docker image and binaries
- Publishes to Docker Hub and GitHub Releases

### 4. Deployment
```bash
# Pull latest image
docker pull your-username/openvpn-admin-go:latest

# Or use specific version
docker pull your-username/openvpn-admin-go:v20240101120000

# Deploy with Docker Compose
cd docker
cp .env.docker.example .env
# Edit .env with your configuration
docker-compose -f docker-compose.production.yml up -d
```

## 🔍 Monitoring and Troubleshooting

### GitHub Actions Logs
1. Go to repository **Actions** tab
2. Click on the workflow run
3. Expand job steps to view detailed logs

### Common Issues

#### Build Failures
- **Go build errors**: Check Go version and dependencies
- **Frontend build errors**: Verify Node.js version and npm install
- **Docker build errors**: Check Dockerfile syntax and base images

#### Docker Hub Push Failures
- **Authentication**: Verify DOCKERHUB_USERNAME and DOCKERHUB_TOKEN
- **Permissions**: Ensure token has write permissions
- **Repository**: Check if repository exists on Docker Hub

#### Release Creation Failures
- **Permissions**: Ensure GitHub Actions has write permissions
- **Token**: Verify GITHUB_TOKEN is available (automatic)

### Debug Commands

```bash
# Check workflow status
gh workflow list
gh run list

# View specific run logs
gh run view <run-id>

# Test Docker image locally
docker run --rm your-username/openvpn-admin-go:latest

# Check image layers
docker history your-username/openvpn-admin-go:latest
```

## 🔒 Security Best Practices

1. **Secrets Management**
   - Never commit secrets to repository
   - Use GitHub encrypted secrets
   - Rotate tokens regularly

2. **Image Security**
   - Use minimal base images (Alpine)
   - Regular security updates
   - Scan images for vulnerabilities

3. **Access Control**
   - Limit Docker Hub token permissions
   - Use branch protection rules
   - Require PR reviews

## 📈 Optimization Tips

1. **Build Speed**
   - Use GitHub Actions cache
   - Optimize Docker layer caching
   - Parallel builds where possible

2. **Image Size**
   - Multi-stage builds
   - Remove unnecessary files
   - Use .dockerignore

3. **Reliability**
   - Health checks
   - Retry mechanisms
   - Proper error handling

## 🆘 Support

If you encounter issues:

1. **Check logs** in GitHub Actions
2. **Review documentation** in this repository
3. **Create an issue** with detailed error information
4. **Check Docker Hub** for image availability
