package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

// GetEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetOpenVPNStatusLogPath retrieves the OpenVPN status log path from env or returns a default.
func GetOpenVPNStatusLogPath() string {
	// Common locations for OpenVPN status log. Adjust default as needed.
	return GetEnvOrDefault("OPENVPN_STATUS_LOG_PATH", "/etc/openvpn/openvpn-status.log")
}

// GetOpenVPNSyncInterval retrieves the sync interval from env, parses it as seconds,
// and returns it as time.Duration. Defaults if not set or parsing fails.
func GetOpenVPNSyncInterval() time.Duration {
	defaultIntervalSeconds := 60
	intervalStr := GetEnvOrDefault("OPENVPN_SYNC_INTERVAL_SECONDS", strconv.Itoa(defaultIntervalSeconds))

	intervalSeconds, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Printf("Warning: Could not parse OPENVPN_SYNC_INTERVAL_SECONDS value '%s' as integer: %v. Using default %d seconds.", intervalStr, err, defaultIntervalSeconds)
		return time.Duration(defaultIntervalSeconds) * time.Second
	}

	if intervalSeconds <= 0 {
		log.Printf("Warning: OPENVPN_SYNC_INTERVAL_SECONDS must be positive. Using default %d seconds.", defaultIntervalSeconds)
		return time.Duration(defaultIntervalSeconds) * time.Second
	}

	return time.Duration(intervalSeconds) * time.Second
}
