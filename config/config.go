package config

import (
	"fmt"
	"os"
	"strconv"
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
	return fmt.Sprintf(`port %d
proto %s
dev tun
ca ca.crt
cert server.crt
key server.key
dh dh.pem
server %s %s
ifconfig-pool-persist ipp.txt
push "redirect-gateway def1 bypass-dhcp"
push "dhcp-option DNS %s"
push "dhcp-option DOMAIN %s"
keepalive 10 120
cipher AES-256-CBC
comp-lzo
user nobody
group nogroup
persist-key
persist-tun
status openvpn-status.log
verb 3
`, c.OpenVPNPort, c.OpenVPNProto, c.OpenVPNServerNetwork, c.OpenVPNServerNetmask, c.DNSServerIP, c.DNSServerDomain)
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