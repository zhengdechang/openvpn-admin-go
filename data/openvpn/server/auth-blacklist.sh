#!/bin/bash
set -e

# Script receives path to a temp file with username and password
CRED_FILE="$1"
username=$(head -n 1 "$CRED_FILE")

# Configurable via environment variables or defaults
BLACKLIST_FILE="${OPENVPN_BLACKLIST_FILE:-/etc/openvpn/blacklist.txt}"
MANAGEMENT_PORT="${OPENVPN_MANAGEMENT_PORT:-7505}"
MANAGEMENT_HOST="127.0.0.1"
LOG_FILE="/var/log/openvpn-auth-blacklist.log"

# Ensure log directory and file exist and are writable (basic setup)
# Note: Running mkdir and touch here might fail if the script itself doesn't have permissions.
# It's better to ensure the log directory/file is prepared by a deployment script or OpenVPN's setup.
# For simplicity in this script, we'll attempt it, but it's not guaranteed to succeed.
if [ ! -d "$(dirname "$LOG_FILE")" ]; then
  mkdir -p "$(dirname "$LOG_FILE")"
fi
touch "$LOG_FILE" # Ensure file exists, created if not

echo "$(date '+%Y-%m-%d %H:%M:%S'): Processing auth for user: $username" >> "$LOG_FILE"

if [ ! -f "$BLACKLIST_FILE" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S'): Blacklist file $BLACKLIST_FILE not found. Allowing user $username." >> "$LOG_FILE"
    exit 0
fi

# Escape username for sed (handles simple cases like '/')
# Complex usernames (e.g. with '&', '\', or other special sed characters) might need more robust escaping.
# Using a temporary variable for clarity.
USERNAME_FOR_SED=$(echo "$username" | sed 's/[\/&]/\\&/g')

if grep -qx "$username" "$BLACKLIST_FILE"; then
    echo "$(date '+%Y-%m-%d %H:%M:%S'): User $username is blacklisted. Attempting removal and disconnection." >> "$LOG_FILE"

    # Remove from blacklist - use the escaped version for sed pattern
    # The pattern `^${USERNAME_FOR_SED}$` ensures matching the whole line exactly.
    if sed -i "/^${USERNAME_FOR_SED}$/d" "$BLACKLIST_FILE"; then
        echo "$(date '+%Y-%m-%d %H:%M:%S'): User $username removed from blacklist $BLACKLIST_FILE." >> "$LOG_FILE"
    else
        echo "$(date '+%Y-%m-%d %H:%M:%S'): Failed to remove user $username from blacklist $BLACKLIST_FILE. User might not exist or sed error." >> "$LOG_FILE"
        # Decide if to proceed with kill attempt or exit. For now, proceed.
    fi

    # Disconnect via management interface
    # Ensure nc is available
    if ! command -v nc &> /dev/null; then
        echo "$(date '+%Y-%m-%d %H:%M:%S'): nc (netcat) command not found. Cannot disconnect user $username." >> "$LOG_FILE"
        exit 1 # Exit with failure, as disconnection is not possible
    fi

    # Send kill command. Using printf for better control over newlines.
    # The management interface might take a moment to process.
    # Adding a small timeout to nc if it hangs, e.g. `nc -w 5` for 5 seconds.
    if printf "kill %s\nexit\n" "$username" | nc -w 5 "$MANAGEMENT_HOST" "$MANAGEMENT_PORT"; then
        echo "$(date '+%Y-%m-%d %H:%M:%S'): Kill command sent for user $username to $MANAGEMENT_HOST:$MANAGEMENT_PORT." >> "$LOG_FILE"
    else
        echo "$(date '+%Y-%m-%d %H:%M:%S'): Failed to send kill command for user $username to $MANAGEMENT_HOST:$MANAGEMENT_PORT. nc exit code: $?" >> "$LOG_FILE"
        # Even if nc fails, we still want to deny auth as user was in blacklist
    fi

    exit 1 # Authentication failure
else
    echo "$(date '+%Y-%m-%d %H:%M:%S'): User $username is not blacklisted. Allowing." >> "$LOG_FILE"
    exit 0 # Authentication success
fi
