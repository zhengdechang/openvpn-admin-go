#!/bin/bash
set -e

WEB_PORT=${WEB_PORT:-8085}
export WEB_PORT

# 后端容器入口：由 supervisord 托管 openvpn-server（VPN 守护进程）与
# openvpn-go-api（Go API），不再直接 exec `openvpn-go web`。
# 这样 OpenVPN 会随容器一起启动，且 API 通过 supervisorctl 管理 VPN 进程时
# 能找到 /etc/supervisor/supervisord.conf 的 socket。
# openvpn-go-api 启动时 main() 仍会执行 InitCore（goose 迁移 + seed + 同步）。
# 前端已是独立容器，本容器不再启动任何前端进程。
if [ $# -eq 0 ] || { [ "$1" = "openvpn-go" ] && [ $# -eq 1 ]; }; then
    echo "[INFO] Starting supervisord (manages openvpn-server + openvpn-go-api on port ${WEB_PORT})"
    mkdir -p /var/log/supervisor /var/run
    exec supervisord -n -c /etc/supervisor/supervisord.conf
fi

# 显式命令（如 migrate / 交互菜单）透传执行
exec "$@"
