package openvpn

import (
	"fmt"
	"os"
	"openvpn-admin-go/utils"
)

// GetServerConfigTemplate 获取服务器配置模板
func GetServerConfigTemplate() string {
	// 从环境变量读取配置
	port := utils.GetEnvOrDefault("OPENVPN_PORT", "1194")
	proto := utils.GetEnvOrDefault("OPENVPN_PROTO", "udp")
	serverNetwork := utils.GetEnvOrDefault("OPENVPN_SERVER_NETWORK", "10.8.0.0")
	serverNetmask := utils.GetEnvOrDefault("OPENVPN_SERVER_NETMASK", "255.255.255.0")
	serverHostname := utils.GetEnvOrDefault("OPENVPN_SERVER_HOSTNAME", "")
	serverIP := utils.GetEnvOrDefault("OPENVPN_SERVER_IP", "")
	dnsServer := utils.GetEnvOrDefault("DNS_SERVER_IP", "8.8.8.8")
	dnsDomain := utils.GetEnvOrDefault("DNS_SERVER_DOMAIN", "")
	useCRL := utils.GetEnvOrDefault("OPENVPN_USE_CRL", "false")
	syncCerts := utils.GetEnvOrDefault("OPENVPN_SYNC_CERTS", "false")

	config := fmt.Sprintf(`port %s
proto %s
dev tun
ca ca.crt
cert server.crt
key server.key
dh dh.pem
server %s %s
ifconfig-pool-persist ipp.txt
push "redirect-gateway def1 bypass-dhcp"
push "dhcp-option DNS %s"`, 
		port, proto, serverNetwork, serverNetmask, dnsServer)

	// 如果配置了DNS域名，添加域名推送
	if dnsDomain != "" {
		config += fmt.Sprintf("\npush \"dhcp-option DOMAIN %s\"", dnsDomain)
	}

	// 如果配置了服务器主机名，添加推送
	if serverHostname != "" {
		config += fmt.Sprintf("\npush \"dhcp-option DOMAIN-SEARCH %s\"", serverHostname)
	}

	// 如果配置了服务器IP，添加推送
	if serverIP != "" {
		config += fmt.Sprintf("\npush \"dhcp-option DNS %s\"", serverIP)
	}

	// 如果启用了CRL
	if useCRL == "true" {
		config += "\ncrl-verify crl.pem"
	}

	// 如果启用了证书同步
	if syncCerts == "true" {
		config += "\nclient-cert-not-required"
	}

	// 添加其他通用配置
	config += `
keepalive 10 120
cipher AES-256-GCM
comp-lzo
user nobody
group nogroup
persist-key
auth SHA256
persist-tun
status openvpn-status.log
verb 3`

	return config
}

// getEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}