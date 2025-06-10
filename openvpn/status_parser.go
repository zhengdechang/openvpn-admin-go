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
	logHeaderTimeLayout   = "2006-01-02 15:04:05"
	onlineThresholdDuration = 5 * time.Minute
	jancsitechPrefix      = "Jancsitech-"
)

// OpenVPNClientStatus holds information for each client parsed from the status log.
type OpenVPNClientStatus struct {
	CommonName            string    `json:"commonName"`     // Often used as UserID
	Username              string    `json:"username"`       // Added: For derived username
	RealAddress           string    `json:"realAddress"`    // Connection IP
	VirtualAddress        string    `json:"virtualAddress"` // Ensure this is populated from ROUTING TABLE
	BytesReceived         int64     `json:"bytesReceived"`
	BytesSent             int64     `json:"bytesSent"`
	ConnectedSince        time.Time `json:"connectedSince"` // To be parsed from CLIENT_LIST
	LastRef               time.Time `json:"lastRef"`        // To be parsed from ROUTING TABLE
	IsOnline              bool      `json:"isOnline"`       // Added: Calculated online status
	OnlineDurationSeconds int64     `json:"onlineDurationSeconds"` // Will be calculated based on ConnectedSince
}

// ParseStatusLog reads and parses an OpenVPN status log file.
// It returns a slice of client statuses, the log update time, and an error if any.
func ParseStatusLog(logPath string) ([]OpenVPNClientStatus, time.Time, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to open status log file %s: %w", logPath, err)
	}
	defer file.Close()

	var clients []OpenVPNClientStatus
	var logUpdateTime time.Time

	clientDataMap := make(map[string]OpenVPNClientStatus)
	routingDataMap := make(map[string]struct {
		VirtualAddress string
		LastRefTime    time.Time
	})

	scanner := bufio.NewScanner(file)
	var parsingClientList, parsingRoutingTable bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "Updated,") {
			parts := strings.Split(line, ",")
			if len(parts) == 2 {
				logUpdateTime, _ = time.Parse(logHeaderTimeLayout, strings.TrimSpace(parts[1]))
				// If parsing fails, logUpdateTime remains zero, which is handled.
			}
			continue
		}

		if strings.HasPrefix(line, "OpenVPN CLIENT LIST") { // Marks the beginning of the client list section header
			parsingClientList = false // Ensure we don't start parsing from "OpenVPN CLIENT LIST" itself
			parsingRoutingTable = false
			continue
		}

		if line == "Common Name,Real Address,Bytes Received,Bytes Sent,Connected Since" {
			parsingClientList = true
			parsingRoutingTable = false
			continue // Skip header line
		}

		if line == "ROUTING TABLE" { // Marks the end of CLIENT_LIST and start of ROUTING_TABLE header block
			parsingClientList = false
			parsingRoutingTable = false // Wait for actual header
			continue
		}

		if line == "Virtual Address,Common Name,Real Address,Last Ref" {
			parsingClientList = false
			parsingRoutingTable = true
			continue // Skip header line
		}

		if line == "GLOBAL STATS" || line == "END" { // GLOBAL STATS or END marks the end of ROUTING_TABLE
			parsingClientList = false
			parsingRoutingTable = false
			continue
		}

		if parsingClientList {
			parts := strings.Split(line, ",")
			if len(parts) == 5 {
				commonName := strings.TrimSpace(parts[0])
				realAddress := strings.TrimSpace(parts[1])
				bytesReceived, _ := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
				bytesSent, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
				connectedSinceStr := strings.TrimSpace(parts[4])
				connectedSince, _ := time.Parse(logHeaderTimeLayout, connectedSinceStr)

				username := commonName
				if strings.HasPrefix(commonName, jancsitechPrefix) {
					username = strings.TrimPrefix(commonName, jancsitechPrefix)
				}

				clientDataMap[commonName] = OpenVPNClientStatus{
					CommonName:     commonName,
					Username:       username,
					RealAddress:    realAddress,
					BytesReceived:  bytesReceived,
					BytesSent:      bytesSent,
					ConnectedSince: connectedSince,
				}
			}
		} else if parsingRoutingTable {
			parts := strings.Split(line, ",")
			if len(parts) == 4 {
				virtualAddress := strings.TrimSpace(parts[0])
				commonName := strings.TrimSpace(parts[1])
				// RealAddress from routing table (parts[2]) is ignored for now, taking from CLIENT_LIST
				lastRefStr := strings.TrimSpace(parts[3])
				lastRef, _ := time.Parse(logHeaderTimeLayout, lastRefStr)

				if existing, ok := routingDataMap[commonName]; !ok || lastRef.After(existing.LastRefTime) {
					routingDataMap[commonName] = struct {
						VirtualAddress string
						LastRefTime    time.Time
					}{
						VirtualAddress: virtualAddress,
						LastRefTime:    lastRef,
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
			client.VirtualAddress = routeInfo.VirtualAddress
			client.LastRef = routeInfo.LastRefTime
		}

		if !client.LastRef.IsZero() && !logUpdateTime.IsZero() {
			client.IsOnline = logUpdateTime.Sub(client.LastRef) <= onlineThresholdDuration
		} else if !client.LastRef.IsZero() && logUpdateTime.IsZero() { // Log update time not parsed, fallback to time.Now()
			client.IsOnline = time.Now().Sub(client.LastRef) <= onlineThresholdDuration
		} else {
			client.IsOnline = false
		}


		if !client.ConnectedSince.IsZero() {
			// The old code used time.Since(cs) which is time.Now().Sub(cs). We'll stick to that.
			client.OnlineDurationSeconds = int64(time.Since(client.ConnectedSince).Seconds())
		} else {
			client.OnlineDurationSeconds = 0
		}
		clients = append(clients, client)
	}

	return clients, logUpdateTime, nil
}

// GetStatusFilePath returns the path to the OpenVPN status file.
// This should ideally come from configuration.
func GetStatusFilePath() string {
	// Placeholder: In a real application, this path would come from a configuration file
	// or an environment variable (e.g., via openvpn.LoadConfig()).
	// Common locations: /var/log/openvpn/status.log
	statusPath := os.Getenv("OPENVPN_STATUS_PATH")
	if statusPath != "" {
		return statusPath
	}
	return "/var/log/openvpn/status.log" // 正确的默认路径
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
