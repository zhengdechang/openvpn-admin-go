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
	OpenVPNServerIP        string   `json:"openvpn_server_ip"`
	OpenVPNRoutes          []string `json:"openvpn_routes"`
	OpenVPNClientConfigDir string   `json:"openvpn_client_config_dir"`
	OpenVPNTLSVersion      string   `json:"openvpn_tls_version"`
	OpenVPNTLSKey          string   `json:"openvpn_tls_key"`
	OpenVPNTLSKeyPath      string   `json:"openvpn_tls_key_path"`
	OpenVPNClientToClient  bool     `json:"openvpn_client_to_client"`
	DNSServerIP            string   `json:"dns_server_ip"`
	DNSServerDomain        string   `json:"dns_server_domain"`
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	var err error

	// 加载 OpenVPN 配置
	cfg.OpenVPNPort, err = strconv.Atoi(getEnv("OPENVPN_PORT", "4500"))
	if err != nil {
		return nil, fmt.Errorf("invalid OPENVPN_PORT: %v", err)
	}

	cfg.OpenVPNProto = getEnv("OPENVPN_PROTO", "tcp6")
	cfg.OpenVPNSyncCerts = getEnvBool("OPENVPN_SYNC_CERTS", true)
	cfg.OpenVPNUseCRL = getEnvBool("OPENVPN_USE_CRL", true)
	cfg.OpenVPNServerHostname = getEnv("OPENVPN_SERVER_HOSTNAME", "network.jancsitech.net")
	cfg.OpenVPNServerNetwork = getEnv("OPENVPN_SERVER_NETWORK", "10.8.0.0")
	cfg.OpenVPNServerNetmask = getEnv("OPENVPN_SERVER_NETMASK", "255.255.255.0")
	cfg.OpenVPNServerIP = getEnv("OPENVPN_SERVER_IP", "172.16.10.10")
	cfg.OpenVPNClientToClient = getEnvBool("OPENVPN_CLIENT_TO_CLIENT", false)
	
	// 只在环境变量存在时设置路由
	if routes, exists := os.LookupEnv("OPENVPN_ROUTES"); exists {
		cfg.OpenVPNRoutes = strings.Split(routes, ",")
	} else {
		cfg.OpenVPNRoutes = []string{}
	}
	
	cfg.OpenVPNClientConfigDir = getEnv("OPENVPN_CLIENT_CONFIG_DIR", constants.ClientConfigDir)
	cfg.OpenVPNTLSVersion = getEnv("OPENVPN_TLS_VERSION", "1.2")
	cfg.OpenVPNTLSKey = getEnv("OPENVPN_TLS_KEY", "ta.key")
	cfg.OpenVPNTLSKeyPath = getEnv("OPENVPN_TLS_KEY_PATH", constants.ServerTLSKeyPath)

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