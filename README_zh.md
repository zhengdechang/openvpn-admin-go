[English Documentation](./README.md)

# openvpn-admin-go

本项目是一个基于 Web 的 OpenVPN 管理面板的后端。它使用 Go 构建，并为前端 (openvpn-web) 提供 RESTful API。本项目的主要目标是通过用户友好的 Web 界面简化 OpenVPN 服务器、用户及其配置的管理。

## 功能

- **用户认证和管理：** 为管理员和普通用户提供安全的用户注册和登录功能。
- **部门管理：** 将用户组织到部门中，以便更好地进行访问控制和管理。
- **OpenVPN 服务器配置管理：** 轻松创建、更新和删除 OpenVPN 服务器配置。
- **OpenVPN 客户端配置生成和管理：** 生成客户端配置文件（例如 .ovpn）并管理客户端访问。
- **服务器状态监控：** 监控 OpenVPN 服务器的状态，包括活动连接和流量。
- **客户端连接监控：** 查看当前连接的客户端及其会话详细信息。
- **日志查看：** 访问和查看服务器和客户端日志以进行故障排除和监控。

## 管理界面

### 命令行界面 (CLI)

该应用程序可以通过命令行界面 (CLI) 进行管理。CLI 的主要入口点是 `cmd/main.go`。它支持用于服务器和客户端管理的各种命令，包括管理用户、部门和配置。

要查看可用命令和选项的列表，您可以运行：
```bash
go run cmd/main.go --help
```
或者，如果您已构建可执行文件：
```bash
./your-executable-name --help
```

### Web 服务界面

该应用程序还提供了一个基于 Web 的管理界面。`openvpn-admin-go` 项目本身充当后端 API。使用此 API 的前端用户界面是 `openvpn-web` 项目，通常位于 `openvpn-web/` 目录中。

Web 界面允许执行与 CLI 类似的管理任务，但提供了图形用户界面以便于交互。

## 安装和设置

1.  **先决条件：**
    *   Go 1.18 或更高版本。

2.  **克隆存储库：**
    ```bash
    git clone https://github.com/your-username/openvpn-admin-go.git
    # (如果不同，请替换为实际的存储库 URL)
    cd openvpn-admin-go
    ```

3.  **安装依赖项：**
    ```bash
    go mod tidy
    ```

4.  **配置应用程序：**
    *   复制示例环境文件：
        ```bash
        cp .env.example .env
        ```
    *   编辑 `.env` 文件并设置必要的变量，包括：
        *   用于数据库连接的 `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD`、`DB_NAME`。
        *   用于签署 JWT 令牌的 `JWT_SECRET`。
        *   适用于您环境的其他相关设置。

5.  **运行应用程序：**
    ```bash
    go run main.go
    ```
    或者，如果您的主包位于 `cmd` 目录中：
    ```bash
    go run cmd/main.go
    ```

## 项目结构

以下是关键目录及其用途的概述：

*   `cmd/`: 包含主应用程序入口点。
*   `common/`: 项目中共享的实用程序函数和帮助程序。
*   `constants/`: 项目范围的常量，例如配置键或默认值。
*   `controller/`: 处理传入的 HTTP 请求，对其进行处理并与服务交互。
*   `database/`: 数据库连接管理、迁移和查询构建器。
*   `middleware/`: 用于身份验证、日志记录和 CORS 等任务的 HTTP 中间件。
*   `model/`: 表示应用程序实体（例如用户、OpenVPN 服务器）的数据结构和类型。
*   `openvpn/`: 与 OpenVPN 服务器交互的核心逻辑（例如，管理配置、监控状态）。
*   `router/`: 定义 API 路由并将其映射到控制器处理程序。
*   `services/`: 业务逻辑以及控制器和数据层（数据库、OpenVPN）之间的协调。
*   `template/`: HTML 模板或其他模板文件（如有）（例如，用于客户端配置文件）。

## API 端点

RESTful API 允许前端 (openvpn-web) 和其他客户端与后端交互。API 路由在 `router/` 目录中定义。端点通常按功能分组：

*   **身份验证：** 用于用户登录、注册和令牌管理的端点。
*   **用户：** 管理用户帐户和配置文件。
*   **部门：** 管理用于用户组织的部门。
*   **服务器：** OpenVPN 服务器配置的 CRUD 操作。
*   **客户端：** 生成和管理客户端 VPN 配置。
*   **状态和监控：** 用于检查服务器状态和连接客户端的端点。
*   **日志：** 访问服务器和客户端日志。

有关特定端点的详细信息，请参阅 `router/` 目录中的代码和相应的控制器处理程序。
