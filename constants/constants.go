package constants

import "path/filepath"

// OpenVPN 相关常量
const (
	OpenVPNConfigPath = "/etc/openvpn/server/server.conf"
	// 服务名称
	ServiceName = "openvpn-server@server.service"

	// 服务器配置路径
	ServerConfigPath = "/etc/openvpn/server/server.conf"

	// 服务器证书路径
	ServerCACertPath = "/etc/openvpn/server/ca.crt"
	ServerCAKeyPath  = "/etc/openvpn/server/ca.key"
	ServerCertPath   = "/etc/openvpn/server/server.crt"
	ServerKeyPath    = "/etc/openvpn/server/server.key"
	ServerDHPath     = "/etc/openvpn/server/dh.pem"
	ServerTLSKeyPath = "/etc/openvpn/server/tls-auth.key"

	// Default log paths
	DefaultOpenVPNStatusLogPath = "/etc/openvpn/status.log"
	DefaultOpenVPNLogPath       = "/etc/openvpn/openvpn.log"

	// 服务器状态日志路径
	// ServerStatusLogPath is now set from config in openvpn/config.go
	// This line is kept for compatibility but should be removed in the future.
	ServerStatusLogPath = DefaultOpenVPNStatusLogPath

	// 服务器 IP 分配文件路径
	ServerIPPPath = "/etc/openvpn/server/ipp.txt"

	// 客户端配置目录
	ClientConfigDir = "/etc/openvpn/client"

)

// openssl 扩展文件
var OpenSSLExtFiles = []string{
	"openssl-ca.ext",
	"openssl-server.ext",
	"openssl-client.ext",
}


// GetClientConfigPath 获取客户端配置文件路径
func GetClientConfigPath(username string) string {
	return filepath.Join(ClientConfigDir, username+".ovpn")
}

// GetClientCertPath 获取客户端证书路径
func GetClientCertPath(username string) string {
	return filepath.Join(ClientConfigDir, username+".crt")
}

// GetClientKeyPath 获取客户端密钥路径
func GetClientKeyPath(username string) string {
	return filepath.Join(ClientConfigDir, username+".key")
}
