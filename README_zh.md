[English Documentation](./README.md)

# OpenVPN Admin Go

一个功能完整的 OpenVPN 管理系统，采用 Go 后端和 Next.js 前端构建。本项目为管理 OpenVPN 服务器、用户、证书和连接监控提供了完整的解决方案，通过直观的 Web 界面进行操作。

## 🚀 功能特性

### 核心管理功能

- **🔐 用户认证与授权：** 多角色用户系统（超级管理员、管理员、经理、用户），基于 JWT 的身份验证
- **🏢 部门管理：** 层次化组织结构，便于用户管理和访问控制
- **⚙️ OpenVPN 服务器管理：** 完整的服务器生命周期管理，包括启动、停止、重启和配置更新
- **👥 客户端管理：** 自动化客户端证书生成、配置文件创建和访问控制
- **📊 实时监控：** 实时连接状态、流量监控和会话跟踪
- **📝 全面日志记录：** 服务器和客户端日志访问，支持实时更新

### 高级功能

- **🌐 固定 IP 分配：** 为特定客户端分配静态 IP 地址
- **🔒 客户端访问控制：** 暂停/恢复客户端访问，无需撤销证书
- **📋 证书管理：** 自动化证书生成、续期和撤销
- **🔄 实时同步：** OpenVPN 状态与数据库的自动同步
- **🎯 子网管理：** 配置客户端特定的子网路由
- **📈 使用分析：** 跟踪连接时长、数据传输和使用模式

## 🏗️ 系统架构

本项目由两个主要组件构成：

### 后端 (openvpn-admin-go)

- **开发语言：** Go 1.21+
- **Web 框架：** Gin (HTTP 路由器)
- **数据库：** PostgreSQL 配合 GORM ORM
- **身份验证：** JWT 令牌
- **OpenVPN 集成：** 直接与 OpenVPN 服务系统集成

### 前端 (openvpn-web)

- **开发框架：** Next.js 14 配合 TypeScript
- **UI 库：** Tailwind CSS + Radix UI 组件
- **状态管理：** Zustand
- **国际化：** i18next
- **API 客户端：** Axios

## 🖥️ 管理界面

### 1. 交互式 CLI 菜单

应用程序提供直观的命令行界面和交互式菜单：

```bash
# 启动交互式 CLI
go run main.go

# 或使用特定命令
go run cmd/main.go --help
```

#### 直接运行编译后的二进制（无 Docker）

```bash
# 编译二进制
go build -o bin/openvpn-go .

# 启动交互式菜单
./bin/openvpn-go

# （可选）使用同一个二进制启动 Web/API 服务
# 可在另一终端或后台运行，保证菜单保持可用
./bin/openvpn-go web --port 8085
```

**可用的 CLI 操作：**

- 服务器管理（启动/停止/重启/配置）
- 客户端管理（创建/删除/暂停/恢复）
- Web 服务管理
- 配置查看
- 日志监控

### 2. Web 仪表板

现代化的 Web 界面，默认访问地址：`http://localhost:8085`

**主要特性：**

- 响应式设计，支持桌面和移动设备
- 实时连接监控
- 拖拽式证书管理
- 多语言支持（中文/英文）
- 基于角色的访问控制

## 📦 安装和设置

### 系统要求

- **Go 1.21+** - 后端开发
- **Node.js 18+** - 前端开发
- **OpenVPN** - VPN 服务器软件
- **Linux 系统** - Ubuntu/Debian/CentOS（需要 systemd）
- **Root/Sudo 权限** - 用于 OpenVPN 服务管理

### 快速开始

1. **克隆仓库：**

   ```bash
   git clone <your-repository-url>
   cd openvpn-admin-go
   ```

2. **后端设置：**

   ```bash
   # 安装 Go 依赖
   go mod tidy

   # 创建环境配置文件
   cp .env.example .env
   # 编辑 .env 文件进行配置

   # 运行应用程序（会自动安装依赖）
   go run main.go
   ```

3. **前端设置：**

   ```bash
   cd openvpn-web

   # 安装 Node.js 依赖
   npm install

   # 启动开发服务器
   npm run dev
   ```

4. **访问应用程序：**
   - **CLI 界面：** 运行 `go run main.go` 进入交互式菜单
   - **Web 界面：** 打开 `http://localhost:3000`（前端）或 `http://localhost:8085`（后端 API）
   - **默认登录：** `superadmin@gmail.com` / `superadmin`

### 环境配置

在项目根目录创建 `.env` 文件：

```env
# 数据库配置
DATABASE_URL=postgres://openvpn:openvpn_secret@localhost:5432/openvpn?sslmode=disable

# JWT 配置
JWT_SECRET=your-super-secret-jwt-key

# OpenVPN 配置
OPENVPN_SERVER_HOSTNAME=your-server-ip-or-domain
OPENVPN_PORT=1194
OPENVPN_PROTO=udp
OPENVPN_SERVER_NETWORK=10.8.0.0
OPENVPN_SERVER_NETMASK=255.255.255.0
OPENVPN_CLIENT_CONFIG_DIR=/etc/openvpn/clients
OPENVPN_STATUS_LOG_PATH=/var/log/openvpn/status.log
OPENVPN_LOG_PATH=/var/log/openvpn/openvpn.log

# 可选：LevelDB 额外存储
LEVELDB_PATH=/var/lib/openvpn-manager
```

## 📁 项目结构

### 后端结构

```
openvpn-admin-go/
├── cmd/                    # CLI 命令和入口点
│   ├── main.go            # 主 CLI 界面和交互式菜单
│   ├── server.go          # 服务器管理命令
│   ├── client.go          # 客户端管理命令
│   ├── web.go             # Web 服务管理
│   └── environment.go     # 环境设置和验证
├── controller/            # HTTP 请求处理器
│   ├── auth.go           # 身份验证端点
│   ├── client.go         # 客户端管理 API
│   ├── server.go         # 服务器管理 API
│   ├── department.go     # 部门管理 API
│   └── log.go            # 日志访问 API
├── model/                # 数据模型和数据库架构
│   ├── client.go         # 用户和客户端模型
│   ├── department.go     # 部门模型
│   ├── server.go         # 服务器配置模型
│   └── status.go         # 连接状态和日志模型
├── openvpn/              # OpenVPN 集成层
│   ├── config.go         # 配置管理
│   ├── server.go         # 服务器操作
│   ├── client.go         # 客户端证书和配置生成
│   ├── status_parser.go  # 状态日志解析
│   └── ccd.go            # 客户端特定配置
├── router/               # API 路由定义
├── middleware/           # HTTP 中间件（JWT、RBAC）
├── database/             # 数据库连接和迁移
├── services/             # 后台服务（同步、监控）
├── utils/                # 工具函数
├── constants/            # 应用程序常量
├── template/             # OpenVPN 配置模板
└── main.go              # 应用程序入口点
```

### 前端结构

```
openvpn-web/
├── src/
│   ├── app/              # Next.js App Router 页面
│   │   ├── dashboard/    # 主应用程序仪表板
│   │   ├── auth/         # 身份验证页面
│   │   └── page.tsx      # 首页
│   ├── components/       # 可重用的 React 组件
│   │   ├── ui/           # 基础 UI 组件（按钮、输入等）
│   │   └── layout/       # 布局组件
│   ├── services/         # API 客户端和服务函数
│   ├── types/            # TypeScript 类型定义
│   ├── lib/              # 工具库和配置
│   └── i18n/             # 国际化文件
├── public/               # 静态资源
└── package.json          # 依赖和脚本
```

## 🔌 API 端点

RESTful API 提供全面的管理功能：

### 身份验证和用户管理

- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户身份验证
- `GET /api/user/me` - 获取当前用户资料
- `PATCH /api/user/me` - 更新用户资料
- `POST /api/user/logout` - 用户登出
- `GET /api/user/refresh` - 刷新 JWT 令牌

### 客户端管理

- `GET /api/client` - 列出所有客户端/用户
- `POST /api/client` - 创建新客户端
- `GET /api/client/:id` - 获取客户端详情
- `PUT /api/client/:id` - 更新客户端
- `DELETE /api/client/:id` - 删除客户端
- `GET /api/client/config/:username` - 下载客户端配置
- `POST /api/client/:username/pause` - 暂停客户端访问
- `POST /api/client/:username/resume` - 恢复客户端访问

### 服务器管理

- `GET /api/server/status` - 获取服务器状态
- `POST /api/server/start` - 启动 OpenVPN 服务器
- `POST /api/server/stop` - 停止 OpenVPN 服务器
- `POST /api/server/restart` - 重启 OpenVPN 服务器
- `PUT /api/server/update` - 更新服务器配置

### 部门管理

- `GET /api/departments` - 列出部门
- `POST /api/departments` - 创建部门
- `PUT /api/departments/:id` - 更新部门
- `DELETE /api/departments/:id` - 删除部门

### 监控和日志

- `GET /api/logs/server` - 获取服务器日志
- `GET /api/logs/client` - 获取客户端日志
- `GET /api/client/status/live` - 获取实时连接状态

## 🔐 用户角色和权限

| 角色           | 权限                                   |
| -------------- | -------------------------------------- |
| **超级管理员** | 完整系统访问权限，用户管理，服务器配置 |
| **管理员**     | 用户管理，部门管理，客户端操作         |
| **经理**       | 指定部门内的客户端管理                 |
| **用户**       | 查看个人资料，下载个人 VPN 配置        |

## 🚀 部署

### 生产环境部署

1. **构建应用程序：**

   ```bash
   # 后端
   go build -o openvpn-admin main.go

   # 前端
   cd openvpn-web
   npm run build
   ```

2. **配置 systemd 服务：**

   ```bash
   sudo cp openvpn-admin /usr/local/bin/
   sudo systemctl enable openvpn-admin
   sudo systemctl start openvpn-admin
   ```

3. **设置反向代理（nginx）：**

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

### Docker 部署

```bash
# 使用 Docker Compose 构建和运行
docker-compose up -d
```

## 🛠️ 开发

### 后端开发

```bash
# 热重载运行
go run main.go

# 运行测试
go test ./...

# 格式化代码
go fmt ./...
```

### 前端开发

```bash
cd openvpn-web

# 开发服务器
npm run dev

# 类型检查
npm run type-check

# 代码检查
npm run lint
```

## 📝 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 🤝 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📞 支持

如需支持和问题咨询：

- 在 GitHub 上创建 issue
- 查看 `docs/` 目录中的文档
- 查看 `examples/` 目录中的示例配置
