package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	
	cfg.OpenVPNClientConfigDir = getEnv("OPENVPN_CLIENT_CONFIG_DIR", "/etc/openvpn/client")
	cfg.OpenVPNTLSVersion = getEnv("OPENVPN_TLS_VERSION", "1.2")
	cfg.OpenVPNTLSKey = getEnv("OPENVPN_TLS_KEY", "ta.key")
	cfg.OpenVPNTLSKeyPath = getEnv("OPENVPN_TLS_KEY_PATH", "/etc/openvpn/server/ta.key")

	// 加载 DNS 配置
	cfg.DNSServerIP = getEnv("DNS_SERVER_IP", "10.10.99.44")
	cfg.DNSServerDomain = getEnv("DNS_SERVER_DOMAIN", "corp.jancsitech.net")

	return cfg, nil
}

// GenerateServerConfig 生成 OpenVPN 服务器配置
func (c *Config) GenerateServerConfig() string {
	configDir := "/etc/openvpn/server"

	// 确保TLS相关配置有值
	if c.OpenVPNTLSVersion == "" {
		c.OpenVPNTLSVersion = "1.2"
	}
	if c.OpenVPNTLSKeyPath == "" {
		c.OpenVPNTLSKeyPath = "/etc/openvpn/server/ta.key"
	}

	// 构建路由配置
	var routeConfigs []string
	for _, route := range c.OpenVPNRoutes {
		routeConfigs = append(routeConfigs, fmt.Sprintf(`push "route %s"`, route))
	}

	// 构建客户端到客户端配置
	var clientToClientConfig string
	if c.OpenVPNClientToClient {
		clientToClientConfig = "client-to-client\n"
	}

	// 构建协议相关配置
	var protoConfig string
	if c.OpenVPNProto == "udp" || c.OpenVPNProto == "udp6" {
		protoConfig = "explicit-exit-notify 1\n"
	}

	// 合并所有配置
	config := fmt.Sprintf(`port %d
proto %s
dev tun
ca %s
cert %s
key %s
dh %s
server %s %s
%sifconfig-pool-persist ipp.txt
%s
push "dhcp-option DNS %s"
push "dhcp-option DOMAIN %s"
keepalive 10 120
topology subnet
data-ciphers AES-256-GCM:AES-128-GCM
auth SHA256
tls-server
tls-version-min %s
tls-cipher TLS-ECDHE-ECDSA-WITH-AES-256-GCM-SHA384:TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384:TLS-ECDHE-ECDSA-WITH-AES-128-GCM-SHA256:TLS-ECDHE-RSA-WITH-AES-128-GCM-SHA256
tls-auth %s 0
key-direction 0
user nobody
group nogroup
persist-key
persist-tun
status /var/log/openvpn/status.log
verb 3
%s`, 
		c.OpenVPNPort,
		c.OpenVPNProto,
		filepath.Join(configDir, "ca.crt"),
		filepath.Join(configDir, "server.crt"),
		filepath.Join(configDir, "server.key"),
		filepath.Join(configDir, "dh.pem"),
		c.OpenVPNServerNetwork,
		c.OpenVPNServerNetmask,
		clientToClientConfig,
		strings.Join(routeConfigs, "\n"),
		c.DNSServerIP,
		c.DNSServerDomain,
		c.OpenVPNTLSVersion,
		filepath.Join(configDir, "ta.key"),
		protoConfig,
	)

	return config
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

// 添加配置保存实现
func SaveConfig(cfg *Config) error {
	// 将配置保存到/etc/openvpn/config.json
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}
	
	if err := os.WriteFile("/etc/openvpn/config.json", data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	// 重新生成OpenVPN服务配置
	configContent := cfg.GenerateServerConfig()
	if err := os.WriteFile("/etc/openvpn/server.conf", []byte(configContent), 0644); err != nil {
		return fmt.Errorf("生成服务配置失败: %v", err)
	}
	
	return nil
}

