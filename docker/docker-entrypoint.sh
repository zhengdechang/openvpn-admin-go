#!/bin/bash
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 初始化函数
init_directories() {
    log_info "创建必要的目录..."
    
    # 创建 supervisor 相关目录
    mkdir -p /etc/supervisor/conf.d
    mkdir -p /var/log/supervisor
    mkdir -p /var/run
    
    # 创建 OpenVPN 相关目录
    mkdir -p /etc/openvpn/server
    mkdir -p /etc/openvpn/client
    mkdir -p /var/log/openvpn
    mkdir -p /var/lib/openvpn
    
    # 创建应用相关目录
    mkdir -p /app/data
    mkdir -p /app/logs
    mkdir -p /app/config
    
    # 设置权限
    chown -R nobody:nogroup /var/log/supervisor
    chown -R nobody:nogroup /etc/openvpn
    chown -R nobody:nogroup /var/log/openvpn
    chown -R appuser:appgroup /app
    
    log_info "目录创建完成"
}

# 初始化 supervisor 配置
init_supervisor_config() {
    log_info "初始化 supervisor 配置..."
    
    # 检查是否存在主配置文件
    if [ ! -f "/etc/supervisor/supervisord.conf" ]; then
        log_info "生成 supervisor 主配置文件..."
        cd /app
    /app/openvpn-go supervisor-config --main-only
    fi
    
    log_info "Supervisor 配置初始化完成"
}

# 启动 supervisord
start_supervisord() {
    log_info "启动 supervisord..."
    
    # 检查 supervisord 是否已经运行
    if pgrep supervisord > /dev/null; then
        log_warn "supervisord 已经在运行"
        return 0
    fi
    
    # 启动 supervisord
    supervisord -c /etc/supervisor/supervisord.conf
    
    # 等待 supervisord 启动
    sleep 2
    
    # 验证启动状态
    if pgrep supervisord > /dev/null; then
        log_info "supervisord 启动成功"
    else
        log_error "supervisord 启动失败"
        return 1
    fi
}

# 根据环境变量配置服务
configure_services() {
    log_info "配置服务..."
    
    cd /app
    
    # 配置 Web 服务
    WEB_PORT=${WEB_PORT:-8085}
    WEB_AUTOSTART=${WEB_AUTOSTART:-false}
    
    log_info "配置 OpenVPN-Admin Web 服务 (端口: $WEB_PORT, 自动启动: $WEB_AUTOSTART)"
    /app/openvpn-go supervisor-config --service web --port $WEB_PORT --autostart $WEB_AUTOSTART

    # 配置 OpenVPN 服务
    OPENVPN_AUTOSTART=${OPENVPN_AUTOSTART:-false}

    log_info "配置 OpenVPN 服务 (自动启动: $OPENVPN_AUTOSTART)"
    /app/openvpn-go supervisor-config --service openvpn --autostart $OPENVPN_AUTOSTART

    # 配置前端服务
    FRONTEND_AUTOSTART=${FRONTEND_AUTOSTART:-false}

    log_info "配置前端服务 (自动启动: $FRONTEND_AUTOSTART)"
    /app/openvpn-go supervisor-config --service frontend --autostart $FRONTEND_AUTOSTART
    
    # 重新加载配置
    supervisorctl reread
    supervisorctl update
    
    log_info "服务配置完成"
}

# 启动指定的服务
start_services() {
    log_info "启动服务..."

    # 根据 SERVICE_MODE 环境变量决定启动哪些服务
    SERVICE_MODE=${SERVICE_MODE:-all}

    # 等待一下确保配置完全加载
    sleep 2

    case $SERVICE_MODE in
        "web"|"api")
            log_info "启动 API 服务模式"
            supervisorctl start openvpn-go-api
            ;;
        "openvpn"|"backend")
            log_info "启动 OpenVPN 服务模式"
            supervisorctl start openvpn-server
            ;;
        "frontend")
            log_info "启动前端模式"
            supervisorctl start openvpn-frontend
            ;;
        "all"|*)
            log_info "启动所有服务"
            supervisorctl start openvpn-go-api || true
            supervisorctl start openvpn-frontend || true
            if [ "$OPENVPN_AUTOSTART" = "true" ]; then
                supervisorctl start openvpn-server || true
            fi
            ;;
    esac

    log_info "服务启动完成"
}

# 显示服务状态
show_status() {
    log_info "服务状态:"
    supervisorctl status
}

# 主函数
main() {
    log_info "OpenVPN Admin Docker 容器启动..."
    
    # 初始化
    init_directories
    init_supervisor_config
    
    # 启动 supervisord
    start_supervisord
    
    # 配置服务
    configure_services
    
    # 启动服务
    start_services
    
    # 显示状态
    show_status
    
    log_info "容器初始化完成"
    
    # 如果有参数传入，执行指定命令
    if [ $# -gt 0 ]; then
        log_info "执行命令: $@"
        exec "$@"
    else
        # 否则保持容器运行
        log_info "容器保持运行状态，使用 supervisorctl 管理服务"
        tail -f /var/log/supervisor/supervisord.log
    fi
}

# 信号处理
trap 'log_info "收到停止信号，正在关闭服务..."; supervisorctl shutdown; exit 0' SIGTERM SIGINT

# 执行主函数
main "$@"
