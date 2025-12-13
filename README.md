[中文文档](./README_zh.md)

# OpenVPN Admin Go

A comprehensive OpenVPN management system built with Go backend and Next.js frontend. This project provides a complete solution for managing OpenVPN servers, users, certificates, and monitoring connections through an intuitive web interface.

## 🚀 Features

### Core Management

- **🔐 User Authentication & Authorization:** Multi-role user system (SuperAdmin, Admin, Manager, User) with JWT-based authentication
- **🏢 Department Management:** Hierarchical organization structure for better user management and access control
- **⚙️ OpenVPN Server Management:** Complete server lifecycle management including start, stop, restart, and configuration updates
- **👥 Client Management:** Automated client certificate generation, configuration file creation, and access control
- **📊 Real-time Monitoring:** Live connection status, traffic monitoring, and session tracking
- **📝 Comprehensive Logging:** Server and client log access with real-time updates

### Advanced Features

- **🌐 Fixed IP Assignment:** Assign static IP addresses to specific clients
- **🔒 Client Access Control:** Pause/resume client access without certificate revocation
- **📋 Certificate Management:** Automated certificate generation, renewal, and revocation
- **🔄 Real-time Synchronization:** Automatic synchronization of OpenVPN status with database
- **🎯 Subnet Management:** Configure client-specific subnet routing
- **📈 Usage Analytics:** Track connection duration, data transfer, and usage patterns

## 🏗️ Architecture

This project consists of two main components:

### Backend (openvpn-admin-go)

- **Language:** Go 1.21+
- **Framework:** Gin (HTTP router)
- **Database:** SQLite with GORM ORM
- **Authentication:** JWT tokens
- **OpenVPN Integration:** Direct system integration with OpenVPN service

### Frontend (openvpn-web)

- **Framework:** Next.js 14 with TypeScript
- **UI Library:** Tailwind CSS + Radix UI components
- **State Management:** Zustand
- **Internationalization:** i18next
- **API Client:** Axios

## 🖥️ Management Interfaces

### 1. Interactive CLI Menu

The application provides an intuitive command-line interface with interactive menus:

```bash
# Start the interactive CLI
go run main.go

# Or use specific commands
go run cmd/main.go --help
```

#### Run the compiled binary (no Docker)

```bash
# Build the binary
go build -o bin/openvpn-go .

# Launch the interactive menu
./bin/openvpn-go

# (Optional) Start the web/API service from the same binary
# Run this in another terminal or background it so the menu stays interactive
./bin/openvpn-go web --port 8085
```

**Available CLI Operations:**

- Server management (start/stop/restart/configure)
- Client management (create/delete/pause/resume)
- Web service management
- Configuration viewing
- Log monitoring

### 2. Web Dashboard

Modern web interface accessible at `http://localhost:8085` (default):

**Key Features:**

- Responsive design for desktop and mobile
- Real-time connection monitoring
- Drag-and-drop certificate management
- Multi-language support (English/Chinese)
- Role-based access control

## 📦 Installation and Setup

### Prerequisites

- **Go 1.21+** - Backend development
- **Node.js 18+** - Frontend development
- **OpenVPN** - VPN server software
- **Linux System** - Ubuntu/Debian/CentOS (systemd required)
- **Root/Sudo Access** - For OpenVPN service management

### Quick Start

1. **Clone the Repository:**

   ```bash
   git clone <your-repository-url>
   cd openvpn-admin-go
   ```

2. **Backend Setup:**

   ```bash
   # Install Go dependencies
   go mod tidy

   # Create environment file
   cp .env.example .env
   # Edit .env with your configuration

   # Run the application (will auto-install dependencies)
   go run main.go
   ```

3. **Frontend Setup:**

   ```bash
   cd openvpn-web

   # Install Node.js dependencies
   npm install

   # Start development server
   npm run dev
   ```

4. **Access the Application:**
   - **CLI Interface:** Run `go run main.go` for interactive menu
   - **Web Interface:** Open `http://localhost:3000` (frontend) or `http://localhost:8085` (backend API)
   - **Default Login:** `superadmin@gmail.com` / `superadmin`

### Environment Configuration

Create a `.env` file in the project root:

```env
# Database Configuration
DB_PATH=data/db.sqlite3

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key

# OpenVPN Configuration
OPENVPN_SERVER_HOSTNAME=your-server-ip-or-domain
OPENVPN_PORT=1194
OPENVPN_PROTO=udp
OPENVPN_SERVER_NETWORK=10.8.0.0
OPENVPN_SERVER_NETMASK=255.255.255.0
OPENVPN_CLIENT_CONFIG_DIR=/etc/openvpn/clients
OPENVPN_STATUS_LOG_PATH=/var/log/openvpn/status.log
OPENVPN_LOG_PATH=/var/log/openvpn/openvpn.log

# Optional: LevelDB for additional storage
LEVELDB_PATH=/var/lib/openvpn-manager
```

## 📁 Project Structure

### Backend Structure

```
openvpn-admin-go/
├── cmd/                    # CLI commands and entry points
│   ├── main.go            # Main CLI interface with interactive menu
│   ├── server.go          # Server management commands
│   ├── client.go          # Client management commands
│   ├── web.go             # Web service management
│   └── environment.go     # Environment setup and validation
├── controller/            # HTTP request handlers
│   ├── auth.go           # Authentication endpoints
│   ├── client.go         # Client management API
│   ├── server.go         # Server management API
│   ├── department.go     # Department management API
│   └── log.go            # Log access API
├── model/                # Data models and database schemas
│   ├── client.go         # User and client models
│   ├── department.go     # Department model
│   ├── server.go         # Server configuration model
│   └── status.go         # Connection status and logging models
├── openvpn/              # OpenVPN integration layer
│   ├── config.go         # Configuration management
│   ├── server.go         # Server operations
│   ├── client.go         # Client certificate and config generation
│   ├── status_parser.go  # Status log parsing
│   └── ccd.go            # Client-specific configurations
├── docker/               # Docker deployment files
│   ├── Dockerfile.combined   # Multi-service Docker image
│   ├── docker-compose.production.yml
│   ├── .env.docker.example   # Environment template
│   └── README.md            # Docker deployment guide
├── .github/workflows/    # CI/CD pipeline configuration
├── docs/                 # Documentation
├── scripts/              # Utility scripts
├── router/               # API route definitions
├── middleware/           # HTTP middleware (JWT, RBAC)
├── database/             # Database connection and migrations
├── services/             # Background services (sync, monitoring)
├── utils/                # Utility functions
├── constants/            # Application constants
├── template/             # OpenVPN configuration templates
└── main.go              # Application entry point
```

### Frontend Structure

```
openvpn-web/
├── src/
│   ├── app/              # Next.js App Router pages
│   │   ├── dashboard/    # Main application dashboard
│   │   ├── auth/         # Authentication pages
│   │   └── page.tsx      # Landing page
│   ├── components/       # Reusable React components
│   │   ├── ui/           # Base UI components (buttons, inputs, etc.)
│   │   └── layout/       # Layout components
│   ├── services/         # API client and service functions
│   ├── types/            # TypeScript type definitions
│   ├── lib/              # Utility libraries and configurations
│   └── i18n/             # Internationalization files
├── public/               # Static assets
├── nginx.conf            # Nginx reverse proxy configuration
└── package.json          # Dependencies and scripts
```

## 🔌 API Endpoints

The RESTful API provides comprehensive management capabilities:

### Authentication & User Management

- `POST /api/user/register` - User registration
- `POST /api/user/login` - User authentication
- `GET /api/user/me` - Get current user profile
- `PATCH /api/user/me` - Update user profile
- `POST /api/user/logout` - User logout
- `GET /api/user/refresh` - Refresh JWT token

### Client Management

- `GET /api/client` - List all clients/users
- `POST /api/client` - Create new client
- `GET /api/client/:id` - Get client details
- `PUT /api/client/:id` - Update client
- `DELETE /api/client/:id` - Delete client
- `GET /api/client/config/:username` - Download client configuration
- `POST /api/client/:username/pause` - Pause client access
- `POST /api/client/:username/resume` - Resume client access

### Server Management

- `GET /api/server/status` - Get server status
- `POST /api/server/start` - Start OpenVPN server
- `POST /api/server/stop` - Stop OpenVPN server
- `POST /api/server/restart` - Restart OpenVPN server
- `PUT /api/server/update` - Update server configuration

### Department Management

- `GET /api/departments` - List departments
- `POST /api/departments` - Create department
- `PUT /api/departments/:id` - Update department
- `DELETE /api/departments/:id` - Delete department

### Monitoring & Logs

- `GET /api/logs/server` - Get server logs
- `GET /api/logs/client` - Get client logs
- `GET /api/client/status/live` - Get live connection status

## 🔐 User Roles & Permissions

| Role           | Permissions                                               |
| -------------- | --------------------------------------------------------- |
| **SuperAdmin** | Full system access, user management, server configuration |
| **Admin**      | User management, department management, client operations |
| **Manager**    | Client management within assigned department              |
| **User**       | View own profile, download own VPN configuration          |

## 🤖 Automated CI/CD

### GitHub Actions Workflow

The project includes automated build and release pipeline:

- **Trigger**: PR merge to `main` branch
- **Builds**: Multi-platform Go binaries (Linux, Windows, macOS)
- **Docker**: Combined frontend+backend image
- **Publishes**: Docker Hub + GitHub Releases

#### Required Secrets

Configure in GitHub repository settings:
```
DOCKERHUB_USERNAME=your-dockerhub-username
DOCKERHUB_TOKEN=your-dockerhub-access-token
```

#### Automated Releases

Every merge to `main` automatically:
1. 🔨 Builds Go binaries for all platforms
2. 🐳 Creates and pushes Docker image
3. 📦 Creates GitHub release with binaries
4. 🏷️ Tags with timestamp version

## 🚀 Deployment

### Production Deployment

1. **Build the applications:**

   ```bash
   # Backend
   go build -o openvpn-go main.go

   # Frontend
   cd openvpn-web
   npm run build
   ```

2. **Configure systemd service:**

   ```bash
   sudo cp openvpn-go /usr/local/bin/
   sudo systemctl enable openvpn-go
   sudo systemctl start openvpn-go
   ```

3. **Setup reverse proxy (nginx):**

   ```nginx
   server {
       listen 80;
       server_name your-domain.com;

       location /api/ {
           proxy_pass http://localhost:8085;
       }

       location / {
           proxy_pass http://localhost:3000;
       }
   }
   ```

### Docker Deployment

#### Quick Start with Pre-built Image

```bash
# Pull and run the latest image
docker run -d \
  --name openvpn-admin \
  --cap-add NET_ADMIN \
  --device /dev/net/tun \
  -p 8085:8085 \
  -p 3000:3000 \
  -p 80:80 \
  -p 1194:1194/udp \
  -v openvpn_data:/app/data \
  -v openvpn_logs:/app/logs \
  -e JWT_SECRET=your-secret-key \
  -e OPENVPN_SERVER_HOSTNAME=your-server-ip \
  zhengdechang/openvpn-admin-go:latest
```

```bash
# Enter the container and launch the interactive provisioning menu
docker exec -it openvpn-admin openvpn-go
```

> **提示：** `openvpn-go` 提供菜单化流程：
> 1. 选择“运行环境检查/安装”即可一键拉起 OpenVPN、OpenSSL、Supervisor 及证书模板。
> 2. “Web 服务管理”会同时管理前后端：后端端口固定为 **8085**，前端端口可通过修改 Nginx 监听端口自定义。
> 3. 同一菜单可查看日志、状态并启动/停止 Web 栈，满足“新机器 → 拉取镜像 → 进入容器执行菜单 → 一键启动前后端”的流程。

#### Ubuntu 22.04 Docker Deployment

For Ubuntu-based deployment with manual service control:

**Build Ubuntu Image:**
```bash
docker build -f docker/Dockerfile.combined -t openvpn-admin:ubuntu .
```

**Interactive Run (Recommended):**
```bash
docker run -it --privileged \
  -p 8085:8085 \
  -p 3000:3000 \
  -p 80:80 \
  -p 1194:1194/udp \
  -v ./data:/app/data \
  -v ./logs:/app/logs \
  -v ./config:/app/config \
  -v ./openvpn:/etc/openvpn \
  openvpn-admin:ubuntu
```

**Background Run:**
```bash
docker run -d --privileged \
  --name openvpn-admin \
  -p 8085:8085 \
  -p 3000:3000 \
  -p 80:80 \
  -p 1194:1194/udp \
  -v ./data:/app/data \
  -v ./logs:/app/logs \
  -v ./config:/app/config \
  -v ./openvpn:/etc/openvpn \
  openvpn-admin:ubuntu

# Enter container
docker exec -it openvpn-admin /bin/bash
```

**Service Management with Supervisor:**

The container now uses supervisor for service management instead of systemd. Services are automatically started via the docker-entrypoint.sh script.

Manual service management:
```bash
# Check service status
supervisorctl status

# Start/stop/restart services
supervisorctl start openvpn-go-api
supervisorctl stop openvpn-go-api
supervisorctl restart openvpn-go-api

supervisorctl start openvpn-server
supervisorctl stop openvpn-server
supervisorctl restart openvpn-server

# View service logs
supervisorctl tail openvpn-go-api
supervisorctl tail openvpn-server

# Follow logs in real-time
supervisorctl tail -f openvpn-go-api
```

#### Docker Compose Deployment

```bash
# Navigate to docker directory
cd docker

# Copy environment configuration
cp .env.docker.example .env
# Edit .env with your settings

# Start services
docker-compose -f docker-compose.production.yml up -d
```

#### Service Modes

The Docker image supports different deployment modes:

- **Full Stack** (`SERVICE_MODE=all`): Both frontend and backend (default)
- **Backend Only** (`SERVICE_MODE=backend`): API server only
- **Frontend Only** (`SERVICE_MODE=frontend`): Web interface only

#### Build from Source

```bash
# Build and run with Docker Compose
docker-compose up -d
```

#### Docker Environment Variables

Container pre-configured environment variables:

- `GIN_MODE=release`
- `NODE_ENV=production`
- `TZ=UTC`
- `DB_PATH=/app/data/db.sqlite3`
- `OPENVPN_CONFIG_DIR=/etc/openvpn`

**Supervisor Service Control Variables:**

- `SERVICE_MODE=all|api|backend|frontend` - Controls which services to start (default: all)
- `WEB_PORT=8085` - API service port (default: 8085)
- `WEB_AUTOSTART=true|false` - Auto-start API service (default: true)
- `FRONTEND_AUTOSTART=true|false` - Auto-start Nginx frontend (default: false)
- `OPENVPN_AUTOSTART=true|false` - Auto-start OpenVPN service (default: false)

#### Docker Troubleshooting

**View Logs:**
```bash
# Application logs
tail -f /app/logs/web.log

# Nginx logs
tail -f /var/log/nginx/error.log
```

**Check Service Status:**
```bash
# Check supervisor status
supervisorctl status

# Check specific service
supervisorctl status openvpn-go-api
supervisorctl status openvpn-server

# Check processes
ps aux | grep supervisord
ps aux | grep openvpn

# Check ports
netstat -tlnp | grep 8085
```

## 🔧 Supervisor Service Management

The application uses supervisor for process management in Docker containers, providing better service control and logging compared to systemd.

### Service Architecture

- **supervisord**: Main process manager
- **openvpn-server**: OpenVPN service process
- **openvpn-go-api**: Web interface process

### Configuration Files

- Main config: `/etc/supervisor/supervisord.conf`
- Service configs: `/etc/supervisor/conf.d/`
  - `openvpn-server.conf` - OpenVPN service configuration
  - `openvpn-go-api.conf` - Web service configuration

### Service Management Commands

```bash
# View all services status
supervisorctl status

# Start services
supervisorctl start openvpn-go-api
supervisorctl start openvpn-server
supervisorctl start all

# Stop services
supervisorctl stop openvpn-go-api
supervisorctl stop openvpn-server
supervisorctl stop all

# Restart services
supervisorctl restart openvpn-go-api
supervisorctl restart openvpn-server

# Reload configuration
supervisorctl reread
supervisorctl update

# View logs
supervisorctl tail openvpn-go-api
supervisorctl tail openvpn-server
supervisorctl tail -f openvpn-go-api  # Follow logs

# Shutdown supervisor
supervisorctl shutdown
```

### Using openvpn-go CLI for Supervisor Management

```bash
# Configure supervisor services
./openvpn-go supervisor-config --main-only                    # Install main config only
./openvpn-go supervisor-config --service api --port 8085      # Configure API service
./openvpn-go supervisor-config --service frontend             # Configure Nginx frontend
./openvpn-go supervisor-config --service openvpn --autostart  # Configure OpenVPN service
```

### Log Files

Supervisor manages centralized logging:

- Supervisor main log: `/var/log/supervisor/supervisord.log`
- Web service logs: `/var/log/supervisor/openvpn-go-api.log`
- OpenVPN service logs: `/var/log/supervisor/openvpn-server.log`
- Error logs: `/var/log/supervisor/*-error.log`

**Restart Services:**
```bash
# Stop service
pkill openvpn-go

# Restart API via supervisor
supervisorctl restart openvpn-go-api
```

## 🛠️ Development

### Backend Development

```bash
# Run with hot reload
go run main.go

# Run tests
go test ./...

# Format code
go fmt ./...
```

### Frontend Development

```bash
cd openvpn-web

# Development server
npm run dev

# Type checking
npm run type-check

# Linting
npm run lint
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

For support and questions:

- Create an issue on GitHub
- Check the documentation in the `docs/` directory
- Review the example configurations in the `examples/` directory
