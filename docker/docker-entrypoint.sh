#!/bin/bash
set -e

APP_BINARY="/app/openvpn-go"
ENABLE_WEB=${ENABLE_WEB:-false}
WEB_PORT=${WEB_PORT:-8085}

if [ $# -eq 0 ] || { [ "$1" = "openvpn-go" ] && [ $# -eq 1 ]; }; then
    if [ "$ENABLE_WEB" = "true" ] || [ "$ENABLE_WEB" = "1" ]; then
        echo "[INFO] ENABLE_WEB is set, starting web service on port ${WEB_PORT}"
        exec "$APP_BINARY" web --port "$WEB_PORT"
    fi
    echo "[INFO] Starting OpenVPN Admin CLI menu"
    exec "$APP_BINARY"
fi

exec "$@"
