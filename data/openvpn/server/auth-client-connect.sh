#!/bin/bash
# client-connect script: called after TLS/cert auth succeeds.
# Uses $common_name (cert CN) to check the blacklist — no password needed from client.
# Exit 0 = allow connection. Exit 1 = deny connection.

USERNAME="$common_name"
BLACKLIST_FILE="${OPENVPN_BLACKLIST_FILE:-/etc/openvpn/server/blacklist.txt}"
LOG_FILE="/var/log/openvpn-auth-blacklist.log"

mkdir -p "$(dirname "$LOG_FILE")" 2>/dev/null
touch "$LOG_FILE" 2>/dev/null

if [ -z "$USERNAME" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S'): ERROR: common_name is empty, denying connection." >> "$LOG_FILE"
    exit 1
fi

if [ ! -f "$BLACKLIST_FILE" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S'): Blacklist not found, allowing $USERNAME." >> "$LOG_FILE"
    exit 0
fi

if grep -qx "$USERNAME" "$BLACKLIST_FILE"; then
    echo "$(date '+%Y-%m-%d %H:%M:%S'): DENIED: $USERNAME is paused (in blacklist)." >> "$LOG_FILE"
    exit 1
fi

echo "$(date '+%Y-%m-%d %H:%M:%S'): ALLOWED: $USERNAME." >> "$LOG_FILE"
exit 0
