#!/bin/bash
set -e

APP_BINARY="/app/openvpn-go"
ENABLE_WEB=${ENABLE_WEB:-true}
WEB_PORT=${WEB_PORT:-8085}
OPENVPN_CONFIG_FILE="/etc/openvpn/server/config.json"

# ── Pre-seed config.json with Docker-appropriate values ─────────────────────
# This runs before the binary starts so InstallEnvironment() picks up the right
# server hostname and protocol when it auto-generates certs and server.conf.
if [ -n "$OPENVPN_SERVER_HOSTNAME" ] || [ -n "$OPENVPN_PROTO" ]; then
    mkdir -p /etc/openvpn/server

    _HOSTNAME="${OPENVPN_SERVER_HOSTNAME:-192.168.2.1}"
    _PROTO="${OPENVPN_PROTO:-tcp6}"

    if [ -f "$OPENVPN_CONFIG_FILE" ]; then
        # Update existing config in-place
        python3 - <<PYEOF
import json
with open("$OPENVPN_CONFIG_FILE") as f:
    cfg = json.load(f)
cfg["openvpn_server_hostname"] = "$_HOSTNAME"
cfg["openvpn_proto"] = "$_PROTO"
with open("$OPENVPN_CONFIG_FILE", "w") as f:
    json.dump(cfg, f, indent=2)
print("[entrypoint] config.json updated: hostname=$_HOSTNAME proto=$_PROTO")
PYEOF
    else
        # Write a minimal seed so the binary merges our values with its defaults
        cat > "$OPENVPN_CONFIG_FILE" << JSONEOF
{
  "openvpn_port": 4500,
  "openvpn_proto": "$_PROTO",
  "openvpn_sync_certs": true,
  "openvpn_use_crl": true,
  "openvpn_server_hostname": "$_HOSTNAME",
  "openvpn_server_network": "10.8.0.0",
  "openvpn_server_netmask": "255.255.255.0",
  "openvpn_routes": [],
  "openvpn_client_config_dir": "/etc/openvpn/client",
  "openvpn_tls_version": "1.2",
  "openvpn_tls_key": "ta.key",
  "openvpn_tls_key_path": "/etc/openvpn/server/tls-auth.key",
  "openvpn_client_to_client": false,
  "dns_server_ip": "",
  "dns_server_domain": "",
  "openvpn_status_log_path": "/etc/openvpn/status.log",
  "openvpn_log_path": "/etc/openvpn/openvpn.log",
  "openvpn_management_port": 7505,
  "openvpn_blacklist_file": "/etc/openvpn/server/blacklist.txt"
}
JSONEOF
        echo "[entrypoint] config.json created: hostname=$_HOSTNAME proto=$_PROTO"
    fi
fi

# ── Start nginx (frontend) ───────────────────────────────────────────────────
if [ "$ENABLE_WEB" = "true" ] || [ "$ENABLE_WEB" = "1" ]; then
    echo "[INFO] Starting nginx frontend..."
    nginx -g "daemon off;" &
fi

# ── Fix auth script permissions (must be readable+executable by nobody) ─────
[ -f "/etc/openvpn/server/auth-blacklist.sh" ] && chmod 755 /etc/openvpn/server/auth-blacklist.sh
# Create auth log file writable by nobody (OpenVPN runs as nobody)
touch /var/log/openvpn-auth-blacklist.log 2>/dev/null && chmod 666 /var/log/openvpn-auth-blacklist.log 2>/dev/null || true

# ── Start OpenVPN server if certs exist (Docker: bypass systemctl) ───────────
OPENVPN_SERVER_CONF="/etc/openvpn/server/server.conf"
if [ -f "$OPENVPN_SERVER_CONF" ] && ! pgrep -x openvpn > /dev/null 2>&1; then
    echo "[INFO] Starting OpenVPN server daemon..."
    mkdir -p /etc/openvpn /var/log/openvpn
    openvpn --config "$OPENVPN_SERVER_CONF" \
            --daemon \
            --log /etc/openvpn/openvpn.log \
            --status /etc/openvpn/status.log 1 \
        && echo "[INFO] OpenVPN daemon started" \
        || echo "[WARN] OpenVPN daemon failed to start (will retry after install)"
fi

# ── Direct web-server invocation: openvpn-go web --port <n> ─────────────────
if [ "$1" = "openvpn-go" ] && [ "$2" = "web" ]; then
    echo "[INFO] Running web server on port ${WEB_PORT}..."
    exec "$@"
fi

# ── Default (no args): web in background + interactive CLI ──────────────────
if [ $# -eq 0 ] || { [ "$1" = "openvpn-go" ] && [ $# -eq 1 ]; }; then
    if [ "$ENABLE_WEB" = "true" ] || [ "$ENABLE_WEB" = "1" ]; then
        echo "[INFO] ENABLE_WEB is set, starting web service on port ${WEB_PORT} in the background"
        "$APP_BINARY" web --port "$WEB_PORT" &
    fi
    echo "[INFO] Starting OpenVPN Admin CLI menu"
    exec "$APP_BINARY"
fi

exec "$@"
