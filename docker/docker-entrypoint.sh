#!/bin/bash
set -e

APP_BINARY="/app/openvpn-go"
ENABLE_WEB=${ENABLE_WEB:-false}
WEB_PORT=${WEB_PORT:-8085}

# When no explicit command is provided, start the CLI menu. If web is enabled,
# start it in the background so the menu remains available.
if [ $# -eq 0 ] || { [ "$1" = "openvpn-go" ] && [ $# -eq 1 ]; }; then
    if [ "$ENABLE_WEB" = "true" ] || [ "$ENABLE_WEB" = "1" ]; then
        echo "[INFO] ENABLE_WEB is set, starting web service on port ${WEB_PORT} in the background"
        "$APP_BINARY" web --port "$WEB_PORT" &
    fi

    echo "[INFO] Starting OpenVPN Admin CLI menu"
    exec "$APP_BINARY"
fi

exec "$@"
