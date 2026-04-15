#!/bin/bash
set -e

APP_BINARY="/app/openvpn-go"
ENABLE_WEB=${ENABLE_WEB:-true}
WEB_PORT=${WEB_PORT:-8085}

# Start nginx (frontend) whenever web mode is enabled
if [ "$ENABLE_WEB" = "true" ] || [ "$ENABLE_WEB" = "1" ]; then
    echo "[INFO] Starting nginx frontend..."
    nginx -g "daemon off;" &
fi

# Direct web-server invocation: openvpn-go web --port <n>
# Run only the web server — single process, no CLI race condition
if [ "$1" = "openvpn-go" ] && [ "$2" = "web" ]; then
    echo "[INFO] Running web server on port ${WEB_PORT}..."
    exec "$@"
fi

# Default (no args or bare "openvpn-go"): start web in background + interactive CLI
if [ $# -eq 0 ] || { [ "$1" = "openvpn-go" ] && [ $# -eq 1 ]; }; then
    if [ "$ENABLE_WEB" = "true" ] || [ "$ENABLE_WEB" = "1" ]; then
        echo "[INFO] ENABLE_WEB is set, starting web service on port ${WEB_PORT} in the background"
        "$APP_BINARY" web --port "$WEB_PORT" &
    fi
    echo "[INFO] Starting OpenVPN Admin CLI menu"
    exec "$APP_BINARY"
fi

exec "$@"
