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
	CommonName     string
	RealAddress    string
	VirtualAddress string // For OpenVPN 2.x, this might be in a separate routing table section or not directly in CLIENT_LIST
	BytesReceived  int64
	BytesSent      int64
	ConnectedSince time.Time
	LastRef        time.Time
	// OpenVPN 2.x status version 2/3 might also include Virtual IPv6, Username, Client ID, Peer ID, Data Channel Cipher
}

// ParseStatusLog reads and parses an OpenVPN status log file.
// The typical format for ConnectedSince and LastRef in OpenVPN status logs is "Day Mon DD HH:MM:SS YYYY".
// Example: "Thu Oct 26 10:00:00 2023"
func ParseStatusLog(logPath string) ([]OpenVPNClientStatus, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open status log file %s: %w", logPath, err)
	}
	defer file.Close()

	var clients []OpenVPNClientStatus
	scanner := bufio.NewScanner(file)
	clientListStarted := false

	// Define the expected time layout for OpenVPN status logs
	// Example: Thu Oct 26 10:00:00 2023
	const timeLayout = "Mon Jan _2 15:04:05 2006" // Note: _2 for day without leading zero, Jan for month abbr.

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "CLIENT_LIST") {
			if !clientListStarted { // First time we see CLIENT_LIST, it's the header
				clientListStarted = true
				// Validate header if necessary, e.g., check number of columns or specific names
				// For now, we assume a known structure for CLIENT_LIST data rows
				continue // Skip header line
			}

			// Subsequent CLIENT_LIST lines are data
			parts := strings.Split(line, ",")
			// Expected format: CLIENT_LIST,Common-Name,Real Address,Virtual Address,Bytes Received,Bytes Sent,Connected Since,Last Ref, ...
			// For OpenVPN 2.x, Virtual Address might not be in this line, or it might be a different column.
			// Status log version 1 (default in older OpenVPN) has fewer fields.
			// Status log version 2/3 (OpenVPN 2.3+) has more fields.
			// This parser assumes a structure common in many setups:
			// Common-Name, Real Address, Bytes Received, Bytes Sent, Connected Since, Last Ref
			// If Virtual Address is present as the 3rd field (index 2) in data rows:
			// CLIENT_LIST,cn,real_ip:port,virt_ip,bytes_recv,bytes_sent,conn_since,last_ref
			// Let's check the number of parts to infer the structure slightly.
			// A common structure has at least 7 parts for client data rows:
			// CLIENT_LIST, CommonName, RealAddress, BytesReceived, BytesSent, ConnectedSince, LastRef
			// (this count includes CLIENT_LIST itself as parts[0])
			// If VirtualAddress is included, it would be 8 parts.

			if len(parts) < 7 { // Minimum expected fields for data row after CLIENT_LIST prefix
				// If parts[0] is "CLIENT_LIST", then len(parts) needs to be at least 1 (for "CLIENT_LIST") + 6 data fields = 7
				// Or if the line itself starts with CommonName (after header was skipped)
				// The problem states "lines starting with CLIENT_LIST," indicating data rows also have this prefix.
				continue // Skip malformed or unexpected lines
			}

			// parts[0] is "CLIENT_LIST"
			// parts[1] is CommonName
			// parts[2] is RealAddress
			// parts[3] is VirtualAddress (IF PRESENT AND EXPECTED AT THIS POSITION)
			// OR parts[3] is BytesReceived if VirtualAddress is not in this line.
			// For simplicity, let's assume a common format where VirtualAddress might be missing or in a different section.
			// We will parse based on a fixed number of fields for now.
			// Let's target the structure:
			// CLIENT_LIST,CommonName,RealAddress,BytesReceived,BytesSent,ConnectedSince,LastRef
			// This means parts[0]=CLIENT_LIST, parts[1]=CN, parts[2]=RealAddr, parts[3]=BytesRecv, parts[4]=BytesSent, parts[5]=ConnSince, parts[6]=LastRef
			// This implies len(parts) should be 7 for this specific structure.

			if len(parts) != 7 && len(parts) != 8 { // Allow for optional VirtualAddress field for now
				// Some logs might have Virtual Address as parts[3], pushing other fields
				// CLIENT_LIST,cn,real_addr,virtual_addr,bytes_r,bytes_s,conn_since,last_ref (8 parts)
				// CLIENT_LIST,cn,real_addr,bytes_r,bytes_s,conn_since,last_ref (7 parts, no virtual_addr in this line)
				continue // Skip if not matching these specific lengths.
			}

			var client OpenVPNClientStatus
			client.CommonName = parts[1]
			client.RealAddress = parts[2]

			idxOffset := 0
			if len(parts) == 8 { // Assuming VirtualAddress is parts[3]
				client.VirtualAddress = parts[3]
				idxOffset = 1
			}

			bytesRecv, err := strconv.ParseInt(parts[3+idxOffset], 10, 64)
			if err != nil {
				// Log parsing error for this field but continue with other fields/clients
				// Or return error: return nil, fmt.Errorf("failed to parse BytesReceived for %s: %w", client.CommonName, err)
				bytesRecv = 0 // Default to 0 on error
			}
			client.BytesReceived = bytesRecv

			bytesSent, err := strconv.ParseInt(parts[4+idxOffset], 10, 64)
			if err != nil {
				bytesSent = 0 // Default to 0 on error
			}
			client.BytesSent = bytesSent

			// ConnectedSince: parts[5+idxOffset]
			connectedSinceStr := parts[5+idxOffset]
			connectedSince, err := time.Parse(timeLayout, connectedSinceStr)
			if err != nil {
				// Handle error, maybe log or set a zero time value
				// return nil, fmt.Errorf("failed to parse ConnectedSince for %s ('%s'): %w", client.CommonName, connectedSinceStr, err)
				connectedSince = time.Time{} // Zero time
			}
			client.ConnectedSince = connectedSince

			// LastRef: parts[6+idxOffset]
			// Ensure part exists before trying to access it, especially if len(parts) was 7 and idxOffset is 0.
			// For len(parts) == 7, max index is 6. parts[6+0] is valid.
			// For len(parts) == 8, max index is 7. parts[6+1] is valid.
			lastRefStr := parts[6+idxOffset]
			lastRef, err := time.Parse(timeLayout, lastRefStr)
			if err != nil {
				lastRef = time.Time{} // Zero time
			}
			client.LastRef = lastRef

			clients = append(clients, client)
		}
		// Other sections to parse (optional, based on requirements):
		// ROUTING_TABLE for virtual IP to common name mapping (if VirtualAddress not in CLIENT_LIST)
		// GLOBAL_STATS for global server statistics
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading status log file: %w", err)
	}

	if !clientListStarted && len(clients) == 0 {
		// This means no "CLIENT_LIST" lines were found at all.
		// It could be an empty log, a log without active clients, or a malformed log.
		// Depending on requirements, this might be an error or just an empty result.
		// For now, returning empty slice and no error is acceptable.
	}

	return clients, nil
}
