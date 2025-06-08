package openvpn

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// OpenVPNClientStatus holds information for each client parsed from the status log.
type OpenVPNClientStatus struct {
	CommonName            string    `json:"commonName"` // Often used as UserID
	RealAddress           string    `json:"realAddress"` // Connection IP
	VirtualAddress        string    `json:"virtualAddress"` // Allocated VPN IP
	BytesReceived         int64     `json:"bytesReceived"`
	BytesSent             int64     `json:"bytesSent"`
	ConnectedSince        time.Time `json:"connectedSince"`
	LastRef               time.Time `json:"lastRef"`
	OnlineDurationSeconds int64     `json:"onlineDurationSeconds"` // Online duration
}

// ParseStatusLog reads and parses an OpenVPN status log file.
func ParseStatusLog(logPath string) ([]OpenVPNClientStatus, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open status log file %s: %w", logPath, err)
	}
	defer file.Close()

	var clients []OpenVPNClientStatus
	scanner := bufio.NewScanner(file)
	clientListStarted := false
	// OpenVPN status log time format e.g. "Mon Jan _2 15:04:05 2006" or "Thu Oct 26 10:00:00 2023"
	const timeLayout = "Mon Jan _2 15:04:05 2006"

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "CLIENT_LIST") {
			if !clientListStarted {
				clientListStarted = true
				// TODO: Potentially parse header to determine column order dynamically
				continue // Skip header line
			}

			parts := strings.Split(line, ",")
			// Example formats:
			// CLIENT_LIST,CommonName,RealAddress,VirtualAddress,BytesReceived,BytesSent,ConnectedSince,LastRef (8 parts)
			// CLIENT_LIST,CommonName,RealAddress,BytesReceived,BytesSent,ConnectedSince,LastRef (7 parts, no VirtualAddress inline)

			if len(parts) < 7 { // Minimum parts for a valid client line (7 if no VirtualAddress, 8 with it)
				continue // Skip malformed lines
			}

			var client OpenVPNClientStatus
			client.CommonName = parts[1]
			client.RealAddress = parts[2]

			idxOffset := 0 // Offset for BytesReceived, BytesSent, etc. if VirtualAddress is present

			// Check if VirtualAddress is likely present in parts[3]
			// A common indicator is if parts[3] contains a "." (as in an IP) AND there are enough subsequent parts
			if len(parts) >= 8 && (strings.Contains(parts[3], ".") || strings.Contains(parts[3], ":")) { // IPv4 or IPv6
				client.VirtualAddress = parts[3]
				idxOffset = 1
			} else if len(parts) == 7 { // No VirtualAddress in this line, fields start directly after RealAddress
				client.VirtualAddress = "" // Explicitly empty
				idxOffset = 0
			} else {
                // Line doesn't match 7 or 8 parts structure or VirtualAddress heuristics
                continue
            }

            // Ensure there are enough parts for the remaining fields based on idxOffset
            // Need 4 more fields: BytesReceived, BytesSent, ConnectedSince, LastRef
            // So, 3 (for CN, RealAddr, OptVirtAddr) + idxOffset + 4 <= len(parts)
            // Simplified: parts[0]=CL, parts[1]=CN, parts[2]=RealAddr.
            // If idxOffset=1 (VirtAddr present), next is parts[3+1]=parts[4] for BytesRcv. We need up to parts[6+1]=parts[7] for LastRef.
            // If idxOffset=0 (No VirtAddr), next is parts[3+0]=parts[3] for BytesRcv. We need up to parts[6+0]=parts[6] for LastRef.
            // So, required length is 7+idxOffset.
            if len(parts) < (7 + idxOffset) {
                continue // Not enough fields for the assumed structure
            }

			client.BytesReceived, _ = strconv.ParseInt(parts[3+idxOffset], 10, 64) // Error handling for ParseInt can be added
			client.BytesSent, _ = strconv.ParseInt(parts[4+idxOffset], 10, 64)

			connectedSinceStr := parts[5+idxOffset]
			cs, err := time.Parse(timeLayout, connectedSinceStr)
			if err == nil {
				client.ConnectedSince = cs
				if !cs.IsZero() { // Calculate duration if ConnectedSince is valid
					client.OnlineDurationSeconds = int64(time.Since(cs).Seconds())
				} else {
					client.OnlineDurationSeconds = 0
				}
			} else {
				client.ConnectedSince = time.Time{} // Ensure zero value on parse error
				client.OnlineDurationSeconds = 0
			}

			// LastRef is often the last field, parts[6+idxOffset]
			lastRefStr := parts[6+idxOffset]
			lr, err := time.Parse(timeLayout, lastRefStr)
			if err == nil {
				client.LastRef = lr
			} else {
				client.LastRef = time.Time{} // Ensure zero value on parse error
			}
			clients = append(clients, client)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading status log file: %w", err)
	}
	return clients, nil
}

// GetStatusFilePath returns the path to the OpenVPN status file.
// This should ideally come from configuration.
func GetStatusFilePath() string {
    // Placeholder: In a real application, this path would come from a configuration file
    // or an environment variable (e.g., via openvpn.LoadConfig()).
    // Common locations: /var/log/openvpn-status.log, /run/openvpn/status.log, or custom.
    // Ensure the OpenVPN server is configured to write to this path.
    statusPath := os.Getenv("OPENVPN_STATUS_PATH")
    if statusPath != "" {
        return statusPath
    }
    return "/tmp/openvpn-status.log" // Default placeholder if not set by env var
}

// GetAllClientStatuses retrieves all client statuses from the OpenVPN status log.
// It uses GetStatusFilePath to determine the log file location.
func GetAllClientStatuses() ([]OpenVPNClientStatus, error) {
    logPath := GetStatusFilePath()
    // ParseStatusLog will handle errors related to file opening.
    return ParseStatusLog(logPath)
}

// GetClientStatus retrieves the status for a specific client by commonName.
func GetClientStatus(commonName string) (*OpenVPNClientStatus, error) {
	statuses, err := GetAllClientStatuses()
	if err != nil {
		return nil, fmt.Errorf("failed to get all client statuses: %w", err)
	}

	for _, status := range statuses {
		if status.CommonName == commonName {
			return &status, nil
		}
	}
	// If no client is found, return nil and no error, or a specific "not found" error.
	// Returning nil, nil is common to indicate "not found" without it being a processing error.
	return nil, nil
}
