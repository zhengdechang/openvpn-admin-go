package openvpn

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"openvpn-admin-go/constants"
)

// Config 存储所有配置
type Config struct {
	OpenVPNPort            int      `json:"openvpn_port"`
	OpenVPNProto           string   `json:"openvpn_proto"`
	OpenVPNSyncCerts       bool     `json:"openvpn_sync_certs"`
	OpenVPNUseCRL          bool     `json:"openvpn_use_crl"`
	OpenVPNServerHostname  string   `json:"openvpn_server_hostname"`
	OpenVPNServerNetwork   string   `json:"openvpn_server_network"`
	OpenVPNServerNetmask   string   `json:"openvpn_server_netmask"`
	OpenVPNRoutes          []string `json:"openvpn_routes"`
	OpenVPNClientConfigDir string   `json:"openvpn_client_config_dir"`
	OpenVPNTLSVersion      string   `json:"openvpn_tls_version"`
	OpenVPNTLSKey          string   `json:"openvpn_tls_key"`
	OpenVPNTLSKeyPath      string   `json:"openvpn_tls_key_path"`
	OpenVPNClientToClient  bool     `json:"openvpn_client_to_client"`
	DNSServerIP            string   `json:"dns_server_ip"`
	DNSServerDomain        string   `json:"dns_server_domain"`
	OpenVPNStatusLogPath   string   `json:"openvpn_status_log_path"`
	OpenVPNLogPath         string   `json:"openvpn_log_path"`
}

// LoadConfig 从服务端配置文件加载配置
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	// 1. Apply hardcoded defaults first
	cfg.OpenVPNPort = constants.DefaultOpenVPNPort
	cfg.OpenVPNProto = constants.DefaultOpenVPNProto
	cfg.OpenVPNSyncCerts = false // Default, not in constants.go but implied
	cfg.OpenVPNUseCRL = false    // Default, not in constants.go but implied
	cfg.OpenVPNServerHostname = "" // No explicit default in constants.go, empty means not set
	cfg.OpenVPNServerNetwork = constants.DefaultOpenVPNServerNetwork
	cfg.OpenVPNServerNetmask = constants.DefaultOpenVPNServerNetmask
	cfg.OpenVPNRoutes = []string{} // Initialize as empty slice
	cfg.OpenVPNClientConfigDir = constants.ClientConfigDir
	cfg.OpenVPNTLSVersion = constants.DefaultOpenVPNTLSVersion
	cfg.OpenVPNTLSKey = constants.DefaultOpenVPNTLSKey // This seems to be a filename, path is separate
	cfg.OpenVPNTLSKeyPath = constants.ServerTLSKeyPath
	cfg.OpenVPNClientToClient = false // Default
	cfg.DNSServerIP = ""           // No explicit default
	cfg.DNSServerDomain = ""       // No explicit default
	cfg.OpenVPNStatusLogPath = constants.DefaultOpenVPNStatusLogPath
	cfg.OpenVPNLogPath = constants.DefaultOpenVPNLogPath

	// Temporary store for routes from server.conf
	var serverConfRoutes []string

	// 2. Parse server.conf and override defaults
	// We will skip the auto-generation of server.conf for now to focus on loading logic.
	// It can be added back later, using the fully resolved config.
	if _, err := os.Stat(constants.ServerConfigPath); err == nil {
		configContent, err := os.ReadFile(constants.ServerConfigPath)
		if err != nil {
			return nil, fmt.Errorf("读取服务端配置文件失败: %v", err)
		}

		lines := strings.Split(string(configContent), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			key := fields[0]
			value := ""
			if len(fields) > 1 {
				value = fields[1]
			}

			switch key {
			case "port":
				if p, convErr := strconv.Atoi(value); convErr == nil {
					cfg.OpenVPNPort = p
				} else {
					// Log or handle malformed port in server.conf?
					// For now, default remains if parsing fails.
				}
			case "proto":
				cfg.OpenVPNProto = value
			case "server":
				if len(fields) >= 3 {
					cfg.OpenVPNServerNetwork = fields[1]
					cfg.OpenVPNServerNetmask = fields[2]
				}
			case "push":
				// Example: push "route 192.168.1.0 255.255.255.0"
				// Need to handle quotes if present in actual server.conf
				if strings.HasPrefix(value, "\"route") || strings.HasPrefix(value, "route") {
					routeParts := fields[1:] // fields are [push, "route, 1.2.3.4, 255.255.255.0"] or [push, route, ...]
					if len(routeParts) > 0 && (strings.HasPrefix(routeParts[0], "\"route") || routeParts[0] == "route") {
						var route string
						if strings.HasPrefix(routeParts[0], "\"route") { // "route 1.2.3.4 255.255.255.0"
							// Combine all parts starting from fields[1] and remove quotes
							combined := strings.Join(fields[1:], " ")
							route = strings.Trim(strings.TrimPrefix(combined, "route"), "\" ")
						} else { // route 1.2.3.4 255.255.255.0
							route = strings.Join(fields[2:], " ")
						}
                         if route != "" {
						    serverConfRoutes = append(serverConfRoutes, route)
                        }
					}
				} else if strings.HasPrefix(value, "\"dhcp-option DNS") { // push "dhcp-option DNS 8.8.8.8"
					// This is an example, current code doesn't parse DNS from server.conf
					// For now, we are only concerned with "route"
				}
			// Add other server.conf specific directives here if needed
			// e.g. client-to-client, tls-version, etc.
			// For now, assuming they are primarily managed by env vars or have simple defaults.
			}
		}
	} // else, if server.conf doesn't exist, we just proceed with defaults and then env vars.

	// 3. Apply Environment Variables, overriding defaults or server.conf values

	// Integer values
	if valStr, exists := os.LookupEnv("OPENVPN_PORT"); exists {
		if valInt, err := strconv.Atoi(valStr); err == nil {
			cfg.OpenVPNPort = valInt
		}
		// Else: log malformed env var? For now, previous value (default/server.conf) is kept.
	}

	// String values
	if val, exists := os.LookupEnv("OPENVPN_PROTO"); exists {
		cfg.OpenVPNProto = val
	}
	if val, exists := os.LookupEnv("OPENVPN_SERVER_HOSTNAME"); exists {
		cfg.OpenVPNServerHostname = val
	}
	if val, exists := os.LookupEnv("OPENVPN_SERVER_NETWORK"); exists {
		cfg.OpenVPNServerNetwork = val
	}
	if val, exists := os.LookupEnv("OPENVPN_SERVER_NETMASK"); exists {
		cfg.OpenVPNServerNetmask = val
	}
	if val, exists := os.LookupEnv("OPENVPN_CLIENT_CONFIG_DIR"); exists {
		cfg.OpenVPNClientConfigDir = val
	}
	if val, exists := os.LookupEnv("OPENVPN_TLS_VERSION"); exists {
		cfg.OpenVPNTLSVersion = val
	}
	if val, exists := os.LookupEnv("OPENVPN_TLS_KEY"); exists { // This is likely the filename, e.g., "ta.key"
		cfg.OpenVPNTLSKey = val
	}
	if val, exists := os.LookupEnv("OPENVPN_TLS_KEY_PATH"); exists { // Full path to the key
		cfg.OpenVPNTLSKeyPath = val
	}
	if val, exists := os.LookupEnv("OPENVPN_STATUS_LOG_PATH"); exists {
		cfg.OpenVPNStatusLogPath = val
	}
	if val, exists := os.LookupEnv("OPENVPN_LOG_PATH"); exists {
		cfg.OpenVPNLogPath = val
	}
	if val, exists := os.LookupEnv("DNS_SERVER_IP"); exists {
		cfg.DNSServerIP = val
	}
	if val, exists := os.LookupEnv("DNS_SERVER_DOMAIN"); exists {
		cfg.DNSServerDomain = val
	}

	// Boolean values
	if valStr, exists := os.LookupEnv("OPENVPN_SYNC_CERTS"); exists {
		if valBool, err := strconv.ParseBool(valStr); err == nil {
			cfg.OpenVPNSyncCerts = valBool
		}
	}
	if valStr, exists := os.LookupEnv("OPENVPN_USE_CRL"); exists {
		if valBool, err := strconv.ParseBool(valStr); err == nil {
			cfg.OpenVPNUseCRL = valBool
		}
	}
	if valStr, exists := os.LookupEnv("OPENVPN_CLIENT_TO_CLIENT"); exists {
		if valBool, err := strconv.ParseBool(valStr); err == nil {
			cfg.OpenVPNClientToClient = valBool
		}
	}

	// Handle routes: append env var routes to server.conf routes
	cfg.OpenVPNRoutes = serverConfRoutes // Start with routes from server.conf
	if routesStr, exists := os.LookupEnv("OPENVPN_ROUTES"); exists {
		if routesStr != "" {
			envRoutes := strings.Split(routesStr, ",")
			for _, route := range envRoutes {
				trimmedRoute := strings.TrimSpace(route)
				if trimmedRoute != "" {
					// Avoid duplicates if desired, though current task is just to append
					cfg.OpenVPNRoutes = append(cfg.OpenVPNRoutes, trimmedRoute)
				}
			}
		}
	}

	// The section that created a default server.conf if it didn't exist has been removed.
	// This should be handled by the caller or a separate setup function,
	// ensuring it uses the fully resolved configuration.
	// For example, after LoadConfig returns, the caller could check if server.conf
	// exists and if not, call SaveConfig (or a modified GenerateServerConfig)
	// to write it out based on the fully resolved cfg.

	return cfg, nil
}

// GenerateServerConfig 生成 OpenVPN 服务器配置
func (c *Config) GenerateServerConfig() (string, error) {
	config, err := RenderServerConfig(c) // RenderServerConfig is not part of this refactoring scope
	if err != nil {
		return "", fmt.Errorf("生成服务器配置失败: %v", err)
	}
	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(cfg *Config) error {
	// 将配置保存到配置文件
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}
	
	configPath := filepath.Join(filepath.Dir(constants.ServerConfigPath), "config.json")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入 config.json 失败: %v", err)
	}

	// Also generate and write the main server.conf
	serverConfigContent, err := cfg.GenerateServerConfig()
	if err != nil {
		return fmt.Errorf("生成 server.conf 内容失败: %v", err)
	}
	if err := os.WriteFile(constants.ServerConfigPath, []byte(serverConfigContent), 0644); err != nil {
		return fmt.Errorf("写入 server.conf (%s) 失败: %v", constants.ServerConfigPath, err)
	}
	
	return nil
} 