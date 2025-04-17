package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"openvpn-admin-go/utils"
)

// Config 存储所有配置
type Config struct {
	OpenVPNPort          int
	OpenVPNProto         string
	OpenVPNSyncCerts     bool
	OpenVPNUseCRL        bool
	OpenVPNServerHostname string
	OpenVPNServerNetmask  string
	OpenVPNServerNetwork  string
	DNSServerIP          string
	DNSServerDomain      string
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

	cfg.OpenVPNProto = getEnv("OPENVPN_PROTO", "tcp")
	cfg.OpenVPNSyncCerts = getEnvBool("OPENVPN_SYNC_CERTS", true)
	cfg.OpenVPNUseCRL = getEnvBool("OPENVPN_USE_CRL", true)
	cfg.OpenVPNServerHostname = getEnv("OPENVPN_SERVER_HOSTNAME", "network.jancsitech.net")
	cfg.OpenVPNServerNetmask = getEnv("OPENVPN_SERVER_NETMASK", "255.255.255.0")
	cfg.OpenVPNServerNetwork = getEnv("OPENVPN_SERVER_NETWORK", "10.9.97.0")

	// 加载 DNS 配置
	cfg.DNSServerIP = getEnv("DNS_SERVER_IP", "10.10.99.44")
	cfg.DNSServerDomain = getEnv("DNS_SERVER_DOMAIN", "corp.jancsitech.net")

	return cfg, nil
}

// GenerateServerConfig 生成 OpenVPN 服务器配置
func (c *Config) GenerateServerConfig() string {
	configDir, err := utils.FindOpenVPNConfigDir()
	if err != nil {
		return ""
	}

	return fmt.Sprintf(`port %d
proto tcp6
dev tun
ca %s
cert %s
key %s
dh %s
client-to-client
server %s %s
ifconfig-pool-persist %s
push "dhcp-option DNS %s"
push "dhcp-option DOMAIN %s"
keepalive 10 120
cipher AES-256-CBC
comp-lzo
user nobody
group nogroup
persist-key
persist-tun
status %s
verb 3`,
		c.OpenVPNPort,
		filepath.Join(configDir, "ca.crt"),
		filepath.Join(configDir, "server.crt"),
		filepath.Join(configDir, "server.key"),
		filepath.Join(configDir, "dh.pem"),
		c.OpenVPNServerNetwork,
		c.OpenVPNServerNetmask,
		filepath.Join(configDir, "ipp.txt"),
		c.DNSServerIP,
		c.DNSServerDomain,
		filepath.Join(configDir, "openvpn-status.log"))
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