package openvpn

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

)

const (
	logHeaderTimeLayout      = "Mon Jan 02 15:04:05 2006" // Adjusted for "Thu Sep 14 10:00:00 2023"
	onlineThresholdDuration    = 5 * time.Minute
	jancsitechPrefix         = "Jancsitech-"
	clientListHeader         = "HEADER,CLIENT_LIST,Common Name,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID,Data Channel Cipher"
	routingTableHeader       = "HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)"
	// Expected field counts for CSV data lines (e.g., "CLIENT_LIST,..."), not including the initial "HEADER" field if it were on the same line.
	clientListFieldsCount   = 13 // CLIENT_LIST + 12 data fields
	routingTableFieldsCount = 6  // ROUTING_TABLE + 5 data fields
)

// OpenVPNClientStatus holds information for each client parsed from the status log.
type OpenVPNClientStatus struct {
	CommonName            string    `json:"commonName"`
	Username              string    `json:"username"`
	RealAddress           string    `json:"realAddress"`
	VirtualAddress        string    `json:"virtualAddress"` // Populated from CLIENT_LIST
	VirtualIPv6Address    string    `json:"virtualIPv6Address"` // Added
	BytesReceived         int64     `json:"bytesReceived"`
	BytesSent             int64     `json:"bytesSent"`
	BytesReceivedFormatted string   `json:"bytesReceivedFormatted"` // New field
	BytesSentFormatted     string   `json:"bytesSentFormatted"`     // New field
	ConnectedSince        time.Time `json:"connectedSince"`      // Parsed from "Connected Since (time_t)"
	ConnectedSinceTimeT   int64     `json:"connectedSinceTimeT"` // Added
	LastRef               time.Time `json:"lastRef"`             // Parsed from "Last Ref (time_t)"
	LastRefTimeT          int64     `json:"lastRefTimeT"`        // Added
	ClientID              string    `json:"clientID"`            // Added
	PeerID                string    `json:"peerID"`              // Added
	DataChannelCipher     string    `json:"dataChannelCipher"`   // Added
	IsOnline              bool      `json:"isOnline"`
	OnlineDurationSeconds int64     `json:"onlineDurationSeconds"`
}

// ParseStatusLog reads and parses an OpenVPN status log file (new format).
// It returns a slice of client statuses, the log update time, and an error if any.
func ParseStatusLog(logPath string) ([]OpenVPNClientStatus, time.Time, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to open status log file %s: %w", logPath, err)
	}
	defer file.Close()

	var clients []OpenVPNClientStatus
	var logUpdateTime time.Time
	var logUpdateTimeEpoch int64

	// clientDataMap stores data primarily from CLIENT_LIST
	clientDataMap := make(map[string]OpenVPNClientStatus)
	// routingDataMap stores data from ROUTING_TABLE, keyed by Common Name for merging
	routingDataMap := make(map[string]struct {
		VirtualAddressFromRoute string    // To differentiate from CLIENT_LIST's VirtualAddress if necessary
		LastRefTime             time.Time
		LastRefTimeT            int64
	})

	scanner := bufio.NewScanner(file)
	var parsingClientList, parsingRoutingTable bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "TITLE,") {
			// Ignore TITLE line or process if needed in the future
			continue
		}

		if strings.HasPrefix(line, "TIME,") {
			parts := strings.Split(line, ",")
			if len(parts) == 3 {
				// parts[0] is "TIME"
				// parts[1] is human-readable time, e.g., "Thu Sep 14 10:00:00 2023"
				// parts[2] is epoch time, e.g., "1694685600"
				logUpdateTime, _ = time.Parse(logHeaderTimeLayout, strings.TrimSpace(parts[1]))
				logUpdateTimeEpoch, _ = strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
				// Use epoch if available and parsing human-readable failed
				if logUpdateTime.IsZero() && logUpdateTimeEpoch > 0 {
					logUpdateTime = time.Unix(logUpdateTimeEpoch, 0)
				}
			}
			continue
		}

		if line == clientListHeader {
			parsingClientList = true
			parsingRoutingTable = false
			continue // Skip header line
		}

		if line == routingTableHeader {
			parsingClientList = false
			parsingRoutingTable = true
			continue // Skip header line
		}

		if strings.HasPrefix(line, "GLOBAL_STATS") || strings.HasPrefix(line, "END") {
			parsingClientList = false
			parsingRoutingTable = false
			continue
		}

		if parsingClientList && strings.HasPrefix(line, "CLIENT_LIST,") {
			parts := strings.Split(line, ",")
			if len(parts) == clientListFieldsCount {
				// CLIENT_LIST line: CLIENT_LIST,Common Name,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID,Data Channel Cipher
				// parts index:       0            1            2             3                4                   5               6            7                 8                       9         10         11       12
				commonName := strings.TrimSpace(parts[1])
				realAddress := strings.TrimSpace(parts[2])
				virtualAddress := strings.TrimSpace(parts[3])
				virtualIPv6Address := strings.TrimSpace(parts[4])
				bytesReceived, _ := strconv.ParseInt(strings.TrimSpace(parts[5]), 10, 64)
				bytesSent, _ := strconv.ParseInt(strings.TrimSpace(parts[6]), 10, 64)
				// "Connected Since" (human-readable parts[7]) is ignored, "(time_t)" (parts[8]) is epoch
				connectedSinceEpoch, _ := strconv.ParseInt(strings.TrimSpace(parts[8]), 10, 64)
				connectedSince := time.Unix(connectedSinceEpoch, 0)
				username := strings.TrimSpace(parts[9])
				clientIDStr := strings.TrimSpace(parts[10]) // Client ID
				peerIDStr := strings.TrimSpace(parts[11])     // Peer ID
				dataChannelCipher := strings.TrimSpace(parts[12])

				clientDataMap[commonName] = OpenVPNClientStatus{
					CommonName:         commonName,
					Username:           username, // Directly from log
					RealAddress:        realAddress,
					VirtualAddress:     virtualAddress, // Directly from CLIENT_LIST
					VirtualIPv6Address: virtualIPv6Address,
					BytesReceived:      bytesReceived,
					BytesSent:          bytesSent,
					BytesReceivedFormatted: formatBytes(bytesReceived),
					BytesSentFormatted:    formatBytes(bytesSent),
					ConnectedSince:     connectedSince,
					ConnectedSinceTimeT:connectedSinceEpoch,
					ClientID:           clientIDStr,
					PeerID:             peerIDStr,
					DataChannelCipher:  dataChannelCipher,
				}
			}
		} else if parsingRoutingTable && strings.HasPrefix(line, "ROUTING_TABLE,") {
			parts := strings.Split(line, ",")
			if len(parts) == routingTableFieldsCount {
				// ROUTING_TABLE line: ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
				// parts index:        0              1                2            3             4         5
				virtualAddressRoute := strings.TrimSpace(parts[1])
				commonName := strings.TrimSpace(parts[2])
				// RealAddress from routing table (parts[3]) might differ, CLIENT_LIST is primary
				// "Last Ref" (human-readable parts[4]) is ignored, "(time_t)" (parts[5]) is epoch
				lastRefEpoch, _ := strconv.ParseInt(strings.TrimSpace(parts[5]), 10, 64)
				lastRefTime := time.Unix(lastRefEpoch, 0)

				// Store route info, preferring the latest LastRef if multiple entries for a common name (unlikely for unique CNs)
				if existing, ok := routingDataMap[commonName]; !ok || lastRefTime.After(existing.LastRefTime) {
					routingDataMap[commonName] = struct {
						VirtualAddressFromRoute string
						LastRefTime             time.Time
						LastRefTimeT            int64
					}{
						VirtualAddressFromRoute: virtualAddressRoute,
						LastRefTime:             lastRefTime,
						LastRefTimeT:            lastRefEpoch,
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, logUpdateTime, fmt.Errorf("error reading status log file: %w", err)
	}

	for cn, client := range clientDataMap {
		if routeInfo, ok := routingDataMap[cn]; ok {
			// client.VirtualAddress is already populated from CLIENT_LIST.
			// We can use routeInfo.VirtualAddressFromRoute if needed for verification or fallback.
			client.LastRef = routeInfo.LastRefTime
			client.LastRefTimeT = routeInfo.LastRefTimeT
		}

		// Calculate IsOnline status
		// Use logUpdateTime (from TIME line). If not available, fallback to time.Now()
		// Check against LastRef from routing table.
		effectiveLogTime := logUpdateTime
		if effectiveLogTime.IsZero() {
			effectiveLogTime = time.Now() // Fallback if TIME line was missing/unparseable
		}

		if !client.LastRef.IsZero() {
			client.IsOnline = effectiveLogTime.Sub(client.LastRef) <= onlineThresholdDuration
		} else {
			// If LastRef is not available, we might infer online status based on ConnectedSince,
			// but typically OpenVPN provides LastRef for active connections.
			// For now, if no LastRef, consider not definitively online by this metric.
			client.IsOnline = false
		}

		// Calculate OnlineDurationSeconds based on ConnectedSince
		if !client.ConnectedSince.IsZero() {
			// Duration should be until now (or log update time)
			client.OnlineDurationSeconds = int64(effectiveLogTime.Sub(client.ConnectedSince).Seconds())
			if client.OnlineDurationSeconds < 0 { // Should not happen if clocks are sane
				client.OnlineDurationSeconds = 0
			}
		} else {
			client.OnlineDurationSeconds = 0
		}
		clients = append(clients, client)
	}

	return clients, logUpdateTime, nil
}

// GetStatusFilePath returns the path to the OpenVPN status file.
// This should ideally come from configuration.
// NOTE: LoadConfig and its related types (Config) are not defined in this snippet.
// Assuming LoadConfig() exists elsewhere and returns a struct with OpenVPNStatusLogPath.
// For this refactoring, we'll keep its usage but acknowledge its external dependency.
func GetStatusFilePath() string {
	// Placeholder for actual config loading logic if not available
	// return "/var/log/openvpn/status.log" // Original fallback

	// Assuming LoadConfig() might be in a different file/package or needs to be defined.
	// For now, let's use a simple default to avoid compilation errors if LoadConfig is missing.
	// This part might need adjustment based on the actual project structure.
	type TempConfig struct {
		OpenVPNStatusLogPath string
	}
	var cfg TempConfig // Simplified
	// cfg, err := LoadConfig() // This line would be used if LoadConfig() is available
	// if err != nil {
	//  return "/var/log/openvpn/status.log"
	// }
	// Simulate config loaded for now, replace with actual LoadConfig() if it exists.
	if cfg.OpenVPNStatusLogPath == "" { // If not loaded or empty
		return "/var/log/openvpn/status.log" // Default path
	}
	return cfg.OpenVPNStatusLogPath
}

// ParseAllClientStatuses retrieves all client statuses from the OpenVPN status log.
// It uses GetStatusFilePath to determine the log file location.
func ParseAllClientStatuses() ([]OpenVPNClientStatus, error) {
	logPath := GetStatusFilePath()
	clients, _, err := ParseStatusLog(logPath) // Discard logUpdateTime for now
	if err != nil {
		return nil, fmt.Errorf("failed to parse status log: %w", err)
	}
	return clients, nil
}

// ParseClientStatus retrieves the status for a specific client by commonName.
func ParseClientStatus(commonName string) (*OpenVPNClientStatus, error) {
	statuses, err := ParseAllClientStatuses()
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

// formatBytes converts bytes to a human-readable string (e.g., "1.5 MB")
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
