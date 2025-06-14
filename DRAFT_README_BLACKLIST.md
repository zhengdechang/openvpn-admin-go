## OpenVPN Client Blacklist Feature

### Overview

The OpenVPN Client Blacklist feature provides a mechanism to automatically disconnect and temporarily block specific users from connecting to the VPN. When a user listed in the blacklist file attempts to authenticate, their connection is denied, their existing session (if any) is terminated, and they are removed from the blacklist, effectively giving them a "one-time" block.

### How it Works

The feature leverages OpenVPN's `auth-user-pass-verify` directive, which delegates the authentication decision to an external script.

1.  **Script Execution:** OpenVPN is configured to use the custom script `/app/scripts/auth-blacklist.sh` for user authentication.
2.  **Username Check:** When a client attempts to connect, OpenVPN executes this script, passing the client's username (and password, though the password is not used by this script for blacklist logic).
3.  **Blacklist Lookup:** The script reads the username and checks for its presence in the configured blacklist file.
4.  **User Blacklisted:** If the username is found in the blacklist file:
    *   The script removes the username from the blacklist file. This means the user will be able to connect on their *next* attempt unless they are re-added.
    *   A command is sent to the OpenVPN management interface (running on `127.0.0.1` at a configured port) to find and `kill` any active session for that username. This ensures the user is disconnected if they were already connected prior to being blacklisted.
    *   The script exits with a failure status, causing OpenVPN to deny the current connection attempt.
5.  **User Not Blacklisted:** If the username is not found in the blacklist file, the script exits with a success status, and OpenVPN proceeds with its normal authentication process (e.g., certificate verification).

### Configuration

The feature is active by default when OpenVPN is configured to use the `auth-blacklist.sh` script as part of its authentication mechanism. The behavior of the script is controlled by environment variables, which are set within the OpenVPN server configuration (`server.conf`) using the `setenv` directive. These `setenv` directives, in turn, source their values from the main application's configuration (`config.json` or corresponding environment variables).

*   **Blacklist File Path:**
    *   **OpenVPN Script Variable:** `OPENVPN_BLACKLIST_FILE` (used by `auth-blacklist.sh`)
    *   **Set by OpenVPN `server.conf` via:** `setenv OPENVPN_BLACKLIST_FILE {{ .OpenVPNBlacklistFile }}`
    *   **Sourced from Go Application Config:** The `{{ .OpenVPNBlacklistFile }}` template variable corresponds to the `openvpn_blacklist_file` field in the application's `config.json` or the `OPENVPN_BLACKLIST_FILE` environment variable supplied to the Go application.
    *   **Default Value in Script:** `/etc/openvpn/blacklist.txt`
    *   **Format:** A plain text file where each line contains a single OpenVPN username to be blacklisted.

*   **OpenVPN Management Port:**
    *   **OpenVPN Script Variable:** `OPENVPN_MANAGEMENT_PORT` (used by `auth-blacklist.sh`)
    *   **Set by OpenVPN `server.conf` via:** `setenv OPENVPN_MANAGEMENT_PORT {{ .OpenVPNManagementPort }}`
    *   **Sourced from Go Application Config:** The `{{ .OpenVPNManagementPort }}` template variable corresponds to the `openvpn_management_port` field in the application's `config.json` or the `OPENVPN_MANAGEMENT_PORT` environment variable supplied to the Go application.
    *   **Default Value in Script:** `7505`
    *   **Note:** The management interface used by the script is hardcoded to connect to `127.0.0.1`. The OpenVPN server configuration should also be set to listen on `127.0.0.1` for the management interface.

### Script Logging

The `auth-blacklist.sh` script logs its actions, including which users are processed, blacklisted, or allowed, to the following file:
*   `/var/log/openvpn-auth-blacklist.log`

Reviewing this log can be helpful for troubleshooting or monitoring the feature's operation. Ensure the directory and file are writable by the user OpenVPN runs as, or adjust permissions/ownership as needed. The script attempts to create the log file and its directory if they don't exist, but this may fail depending on system permissions.

### Manual Blacklisting

To blacklist an OpenVPN user:
1.  Access the server where OpenVPN is running.
2.  Open the blacklist file (default: `/etc/openvpn/blacklist.txt`) with a text editor.
3.  Add the OpenVPN username of the client you wish to blacklist to a new line in the file.
4.  Save the file.

The next time this user attempts to connect, they will be denied access for that attempt, disconnected if already connected, and then removed from the blacklist. To permanently block a user, other measures (like revoking their certificate) should be used. This feature provides a temporary, automated "kick."

---

## config.json Documentation Update Notes

When documenting the main application's `config.json` (or its equivalent configuration sources like environment variables), the following new optional fields should be mentioned:

*   `openvpn_blacklist_file`
    *   **Type:** `string`
    *   **Description:** Specifies the full path to the OpenVPN client blacklist file. This file contains a list of usernames (one per line) that should be temporarily denied access and disconnected. The `auth-blacklist.sh` script reads this path from the `OPENVPN_BLACKLIST_FILE` environment variable, which is set by OpenVPN via the `setenv` directive in its configuration.
    *   **Default:** If not specified, the `auth-blacklist.sh` script will use its internal default of `/etc/openvpn/blacklist.txt`.
    *   **Example:** `"/etc/openvpn/user_blacklist.txt"`

*   `openvpn_management_port`
    *   **Type:** `int`
    *   **Description:** Defines the TCP port on which the OpenVPN server's management interface should listen (and to which the `auth-blacklist.sh` script will connect on `127.0.0.1`). This port is used by the script to send commands to disconnect users. The script reads this port from the `OPENVPN_MANAGEMENT_PORT` environment variable, set by OpenVPN via the `setenv` directive.
    *   **Default:** If not specified, the `auth-blacklist.sh` script will use its internal default of `7505`.
    *   **Example:** `7506`

These settings allow administrators to customize the location of the blacklist file and the OpenVPN management port to suit their environment.
---
