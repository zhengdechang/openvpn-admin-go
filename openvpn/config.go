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

	// 检查配置文件是否存在，如果不存在则创建
	if _, err := os.Stat(constants.ServerConfigPath); os.IsNotExist(err) {
		// 创建配置文件目录
		configDir := filepath.Dir(constants.ServerConfigPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("创建配置目录失败: %v", err)
		}

		// 设置默认配置
		cfg.OpenVPNPort = constants.DefaultOpenVPNPort
		cfg.OpenVPNProto = constants.DefaultOpenVPNProto
		cfg.OpenVPNServerNetwork = constants.DefaultOpenVPNServerNetwork
		cfg.OpenVPNServerNetmask = constants.DefaultOpenVPNServerNetmask
		cfg.OpenVPNServerHostname = getEnv("OPENVPN_SERVER_HOSTNAME",constants.DefaultOPENVPN_SERVER_HOSTNAME)
		cfg.OpenVPNClientToClient = getEnvBool("OPENVPN_CLIENT_TO_CLIENT", false)
		cfg.OpenVPNClientConfigDir = getEnv("OPENVPN_CLIENT_CONFIG_DIR", constants.ClientConfigDir)
		cfg.OpenVPNTLSVersion = getEnv("OPENVPN_TLS_VERSION", constants.DefaultOpenVPNTLSVersion)
		cfg.OpenVPNTLSKey = getEnv("OPENVPN_TLS_KEY",constants.DefaultOpenVPNTLSKey)
		cfg.OpenVPNTLSKeyPath = getEnv("OPENVPN_TLS_KEY_PATH", constants.ServerTLSKeyPath)
		cfg.OpenVPNStatusLogPath = getEnv("OPENVPN_STATUS_LOG_PATH", constants.DefaultOpenVPNStatusLogPath)
		cfg.OpenVPNLogPath = getEnv("OPENVPN_LOG_PATH", constants.DefaultOpenVPNLogPath)

		// 生成默认配置文件
		configContent, err := cfg.GenerateServerConfig()
		if err != nil {
			return nil, fmt.Errorf("生成默认配置文件失败: %v", err)
		}

		// 写入配置文件
		if err := os.WriteFile(constants.ServerConfigPath, []byte(configContent), 0644); err != nil {
			return nil, fmt.Errorf("写入配置文件失败: %v", err)
		}

		return cfg, nil
	}

	// 读取服务端配置文件
	configContent, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return nil, fmt.Errorf("读取服务端配置文件失败: %v", err)
	}

	// 解析配置文件
	lines := strings.Split(string(configContent), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "port":
			cfg.OpenVPNPort, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("解析端口失败: %v", err)
			}
		case "proto":
			cfg.OpenVPNProto = fields[1]
		case "server":
			if len(fields) >= 3 {
				cfg.OpenVPNServerNetwork = fields[1]
				cfg.OpenVPNServerNetmask = fields[2]
			}
		case "push":
			if strings.HasPrefix(fields[1], "route") {
				route := strings.Join(fields[2:], " ")
				cfg.OpenVPNRoutes = append(cfg.OpenVPNRoutes, route)
			}
		}
	}

	// 设置默认值
	if cfg.OpenVPNPort == 0 {
		cfg.OpenVPNPort = constants.DefaultOpenVPNPort
	}
	if cfg.OpenVPNProto == "" {
		cfg.OpenVPNProto = constants.DefaultOpenVPNProto
	}
	if cfg.OpenVPNServerNetwork == "" {
		cfg.OpenVPNServerNetwork = constants.DefaultOpenVPNServerNetwork
	}
	if cfg.OpenVPNServerNetmask == "" {
		cfg.OpenVPNServerNetmask = constants.DefaultOpenVPNServerNetmask
	}

	// 从环境变量加载其他配置
	cfg.OpenVPNServerHostname = getEnv("OPENVPN_SERVER_HOSTNAME", constants.DefaultOPENVPN_SERVER_HOSTNAME)
	cfg.OpenVPNClientToClient = getEnvBool("OPENVPN_CLIENT_TO_CLIENT", false)
	cfg.OpenVPNClientConfigDir = getEnv("OPENVPN_CLIENT_CONFIG_DIR", constants.ClientConfigDir)
	cfg.OpenVPNTLSVersion = getEnv("OPENVPN_TLS_VERSION",constants.DefaultOpenVPNTLSVersion)
	cfg.OpenVPNTLSKey = getEnv("OPENVPN_TLS_KEY",  constants.DefaultOpenVPNTLSKey)
	cfg.OpenVPNTLSKeyPath = getEnv("OPENVPN_TLS_KEY_PATH", constants.ServerTLSKeyPath)
	cfg.OpenVPNStatusLogPath = getEnv("OPENVPN_STATUS_LOG_PATH", constants.DefaultOpenVPNStatusLogPath)
	cfg.OpenVPNLogPath = getEnv("OPENVPN_LOG_PATH", constants.DefaultOpenVPNLogPath)

	// 加载路由配置
	if routes, exists := os.LookupEnv("OPENVPN_ROUTES"); exists {
		cfg.OpenVPNRoutes = strings.Split(routes, ",")
	}

	// 只在环境变量存在时设置 DNS 配置
	if dnsIP, exists := os.LookupEnv("DNS_SERVER_IP"); exists {
		cfg.DNSServerIP = dnsIP
	}
	if dnsDomain, exists := os.LookupEnv("DNS_SERVER_DOMAIN"); exists {
		cfg.DNSServerDomain = dnsDomain
	}

	return cfg, nil
}

// GenerateServerConfig 生成 OpenVPN 服务器配置
func (c *Config) GenerateServerConfig() (string, error) {
	config, err := RenderServerConfig(c)
	if err != nil {
		return "", fmt.Errorf("生成服务器配置失败: %v", err)
	}
	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvBool 获取布尔类型的环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

// getEnvList 获取字符串列表类型的环境变量
func getEnvList(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
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
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	return nil
} 