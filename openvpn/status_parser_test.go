package openvpn

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const sampleLogDataNewFormat = `TITLE,OpenVPN 2.5.11 x86_64-pc-linux-gnu [SSL (OpenSSL)] [LZO] [LZ4] [EPOLL] [PKCS11] [MH/PKTINFO] [AEAD] built on Sep 17 2024
TIME,2025-06-11 13:19:29,1749647969
HEADER,CLIENT_LIST,Common Name,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID,Data Channel Cipher
CLIENT_LIST,ff577df0-b18e-43f0-a11f-ac020be3fa38,10.0.0.1:17704,10.8.0.2,,828967,51576,2025-06-11 10:29:16,1749637756,UNDEF,1,0,AES-256-GCM
HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
ROUTING_TABLE,10.8.0.2,ff577df0-b18e-43f0-a11f-ac020be3fa38,10.0.0.1:17704,2025-06-11 13:19:27,1749647967
GLOBAL_STATS,Max bcast/mcast queue length,1
END
`

func TestParseStatusLog_NewFormat(t *testing.T) {
	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "status.log")
	err := os.WriteFile(logFilePath, []byte(sampleLogDataNewFormat), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp log file: %v", err)
	}

	clients, logTime, err := ParseStatusLog(logFilePath)
	if err != nil {
		t.Fatalf("ParseStatusLog failed: %v", err)
	}

	expectedLogTime := time.Unix(1749647969, 0)
	if !logTime.Equal(expectedLogTime) {
		t.Errorf("Expected log time %v, got %v", expectedLogTime, logTime)
	}

	if len(clients) != 1 {
		t.Fatalf("Expected 1 client, got %d", len(clients))
	}

	client := clients[0]
	expectedCommonName := "ff577df0-b18e-43f0-a11f-ac020be3fa38"
	if client.CommonName != expectedCommonName {
		t.Errorf("Expected CommonName %s, got %s", expectedCommonName, client.CommonName)
	}

	expectedRealAddress := "10.0.0.1:17704"
	if client.RealAddress != expectedRealAddress {
		t.Errorf("Expected RealAddress %s, got %s", expectedRealAddress, client.RealAddress)
	}

	expectedVirtualAddress := "10.8.0.2"
	if client.VirtualAddress != expectedVirtualAddress {
		t.Errorf("Expected VirtualAddress %s, got %s", expectedVirtualAddress, client.VirtualAddress)
	}

	if client.VirtualIPv6Address != "" {
		t.Errorf("Expected empty VirtualIPv6Address, got %s", client.VirtualIPv6Address)
	}

	var expectedBytesReceived int64 = 828967
	if client.BytesReceived != expectedBytesReceived {
		t.Errorf("Expected BytesReceived %d, got %d", expectedBytesReceived, client.BytesReceived)
	}

	var expectedBytesSent int64 = 51576
	if client.BytesSent != expectedBytesSent {
		t.Errorf("Expected BytesSent %d, got %d", expectedBytesSent, client.BytesSent)
	}

	expectedBytesReceivedFormatted := "809.5K" // 828967 / 1024 = 809.538...
	if client.BytesReceivedFormatted != expectedBytesReceivedFormatted {
		t.Errorf("Expected BytesReceivedFormatted %s, got %s", expectedBytesReceivedFormatted, client.BytesReceivedFormatted)
	}

	expectedBytesSentFormatted := "50.4K" // 51576 / 1024 = 50.367... -> rounded to 50.4
	if client.BytesSentFormatted != expectedBytesSentFormatted {
		t.Errorf("Expected BytesSentFormatted %s, got %s", expectedBytesSentFormatted, client.BytesSentFormatted)
	}

	expectedConnectedSince := time.Unix(1749637756, 0)
	if !client.ConnectedSince.Equal(expectedConnectedSince) {
		t.Errorf("Expected ConnectedSince %v, got %v", expectedConnectedSince, client.ConnectedSince)
	}

	if client.ConnectedSinceTimeT != 1749637756 {
		t.Errorf("Expected ConnectedSinceTimeT %d, got %d", int64(1749637756), client.ConnectedSinceTimeT)
	}

	expectedUsername := "UNDEF"
	if client.Username != expectedUsername {
		t.Errorf("Expected Username %s, got %s", expectedUsername, client.Username)
	}

	expectedClientID := "1"
	if client.ClientID != expectedClientID {
		t.Errorf("Expected ClientID %s, got %s", expectedClientID, client.ClientID)
	}

	expectedPeerID := "0"
	if client.PeerID != expectedPeerID {
		t.Errorf("Expected PeerID %s, got %s", expectedPeerID, client.PeerID)
	}

	expectedDataChannelCipher := "AES-256-GCM"
	if client.DataChannelCipher != expectedDataChannelCipher {
		t.Errorf("Expected DataChannelCipher %s, got %s", expectedDataChannelCipher, client.DataChannelCipher)
	}

	expectedLastRef := time.Unix(1749647967, 0)
	if !client.LastRef.Equal(expectedLastRef) {
		t.Errorf("Expected LastRef %v, got %v", expectedLastRef, client.LastRef)
	}

	if client.LastRefTimeT != 1749647967 {
		t.Errorf("Expected LastRefTimeT %d, got %d", int64(1749647967), client.LastRefTimeT)
	}

	// IsOnline: logTime (1749647969) - LastRef (1749647967) = 2s. onlineThresholdDuration is 5min. So, should be true.
	if !client.IsOnline {
		t.Errorf("Expected IsOnline to be true, got false. LogTime: %v, LastRef: %v", logTime, client.LastRef)
	}

	// OnlineDurationSeconds: logTime (1749647969) - ConnectedSince (1749637756) = 10213s
	var expectedOnlineDurationSeconds int64 = 10213
	if client.OnlineDurationSeconds != expectedOnlineDurationSeconds {
		t.Errorf("Expected OnlineDurationSeconds %d, got %d", expectedOnlineDurationSeconds, client.OnlineDurationSeconds)
	}
}

func TestParseStatusLog_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "empty.log")
	_, err := os.Create(logFilePath) // Create an empty file
	if err != nil {
		t.Fatalf("Failed to create empty log file: %v", err)
	}

	clients, logTime, err := ParseStatusLog(logFilePath)
	if err != nil {
		t.Fatalf("ParseStatusLog failed for empty file: %v", err)
	}

	if !logTime.IsZero() { // Expect zero time if no TIME line
		t.Errorf("Expected zero log time for empty file, got %v", logTime)
	}

	if len(clients) != 0 {
		t.Errorf("Expected 0 clients for empty file, got %d", len(clients))
	}
}

func TestParseStatusLog_MalformedFile(t *testing.T) {
	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "malformed.log")
	malformedData := "TITLE,Test\nTIME,ThisIsNotATime,NotAnEpoch\nHEADER,CLIENT_LIST,col1,col2\nCLIENT_LIST,data1\n"
	err := os.WriteFile(logFilePath, []byte(malformedData), 0644)
	if err != nil {
		t.Fatalf("Failed to write malformed log file: %v", err)
	}

	clients, logTime, err := ParseStatusLog(logFilePath)
	if err != nil {
		// Depending on strictness, an error might be acceptable or not.
		// The current parser is quite lenient with malformed lines.
		// For now, let's assume it might not error out but return empty/zero data.
		t.Logf("ParseStatusLog returned an error for malformed file (which might be acceptable): %v", err)
	}

	// Check for best-effort parsing: logTime might be zero if TIME line is bad
	if !logTime.IsZero() {
		t.Errorf("Expected zero log time for malformed TIME line, got %v", logTime)
	}

	// Clients list should ideally be empty if CLIENT_LIST lines are malformed or headers don't match
	if len(clients) != 0 {
		t.Errorf("Expected 0 clients for significantly malformed file, got %d. Clients: %+v", len(clients), clients)
	}
}

func TestParseStatusLog_OnlyHeaders(t *testing.T) {
	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "only_headers.log")
	headersData := `TITLE,OpenVPN
TIME,2025-06-11 13:19:29,1749647969
HEADER,CLIENT_LIST,Common Name,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID,Data Channel Cipher
HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
END`
	err := os.WriteFile(logFilePath, []byte(headersData), 0644)
	if err != nil {
		t.Fatalf("Failed to write headers-only log file: %v", err)
	}

	clients, logTime, err := ParseStatusLog(logFilePath)
	if err != nil {
		t.Fatalf("ParseStatusLog failed for headers-only file: %v", err)
	}

	expectedLogTime := time.Unix(1749647969, 0)
	if !logTime.Equal(expectedLogTime) {
		t.Errorf("Expected log time %v, got %v", expectedLogTime, logTime)
	}

	if len(clients) != 0 {
		t.Errorf("Expected 0 clients for headers-only file, got %d", len(clients))
	}
}

func TestParseStatusLog_ClientWithNoRoutingEntry(t *testing.T) {
	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "client_no_route.log")
	logData := `TITLE,OpenVPN Log
TIME,2025-06-11 14:00:00,1749650400
HEADER,CLIENT_LIST,Common Name,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID,Data Channel Cipher
CLIENT_LIST,client1,192.168.1.100:12345,10.8.0.3,,1000,2000,2025-06-11 13:00:00,1749646800,user1,2,0,AES-256-GCM
HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
END
`
	err := os.WriteFile(logFilePath, []byte(logData), 0644)
	if err != nil {
		t.Fatalf("Failed to write log file: %v", err)
	}

	clients, logTime, err := ParseStatusLog(logFilePath)
	if err != nil {
		t.Fatalf("ParseStatusLog failed: %v", err)
	}

	expectedLogTime := time.Unix(1749650400, 0)
	if !logTime.Equal(expectedLogTime) {
		t.Errorf("Expected log time %v, got %v", expectedLogTime, logTime)
	}

	if len(clients) != 1 {
		t.Fatalf("Expected 1 client, got %d", len(clients))
	}

	client := clients[0]
	if client.CommonName != "client1" {
		t.Errorf("Expected CommonName client1, got %s", client.CommonName)
	}
	if !client.LastRef.IsZero() {
		t.Errorf("Expected zero LastRef for client with no routing entry, got %v", client.LastRef)
	}
	if client.IsOnline { // If no LastRef, IsOnline should be false
		t.Errorf("Expected IsOnline to be false for client with no routing entry, got true")
	}
	// OnlineDurationSeconds should still be calculated based on ConnectedSince and logTime
	// logTime (1749650400) - ConnectedSince (1749646800) = 3600s
	var expectedOnlineDurationSeconds int64 = 3600
	if client.OnlineDurationSeconds != expectedOnlineDurationSeconds {
		t.Errorf("Expected OnlineDurationSeconds %d, got %d", expectedOnlineDurationSeconds, client.OnlineDurationSeconds)
	}
}

// This test is for the old format to ensure the parser can gracefully handle it (e.g. by returning no clients or specific error)
// or if we want to ensure it explicitly DOES NOT parse it.
// The current updated parser will likely not parse old format correctly, which is expected.
func TestParseStatusLog_OldFormatGracefulHandling(t *testing.T) {
	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "old_status.log")
	oldFormatLogData := `OpenVPN CLIENT LIST
Updated,Thu Sep 14 10:00:00 2023
Common Name,Real Address,Bytes Received,Bytes Sent,Connected Since
client1,1.2.3.4:1234,100,200,2023-09-14 09:00:00
ROUTING TABLE
Virtual Address,Common Name,Real Address,Last Ref
10.8.0.1,client1,1.2.3.4:1234,2023-09-14 09:59:00
GLOBAL STATS
Max bcast/mcast queue length,0
END
`
	err := os.WriteFile(logFilePath, []byte(oldFormatLogData), 0644)
	if err != nil {
		t.Fatalf("Failed to write old format log file: %v", err)
	}

	clients, logTime, err := ParseStatusLog(logFilePath)
	if err != nil {
		// An error might be acceptable here, or it might parse with zero clients.
		t.Logf("ParseStatusLog for old format returned error (potentially ignorable): %v", err)
	}

	// For the new parser, logTime might be zero because "Updated," is no longer parsed for time.
	if !logTime.IsZero() {
		t.Errorf("Expected zero logTime when parsing old format, got %v", logTime)
	}
	// Expect no clients as the headers and line prefixes won't match.
	if len(clients) != 0 {
		t.Errorf("Expected 0 clients when parsing old format with new parser, got %d. Clients: %+v", len(clients), clients)
	}
}
