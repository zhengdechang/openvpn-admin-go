/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-07-01 14:40:50
 */
package constants

import "path/filepath"

// OpenVPN 相关常量
const (
	OpenVPNConfigPath = "/etc/openvpn/server/server.conf"
	// Systemd 服务名称
	ServiceName = "openvpn-server@server.service"
	// Web 服务名称
	WebServiceName = "openvpn-go-api.service"

	// Supervisor 服务名称
	SupervisorOpenVPNServiceName = "openvpn-server"
	SupervisorWebServiceName     = "openvpn-go-api"
	SupervisorServiceGroup       = "openvpn-services"

	// 服务器配置路径
	ServerConfigPath = "/etc/openvpn/server/server.conf"

	// 服务器证书路径
	ServerCACertPath = "/etc/openvpn/server/ca.crt"
	ServerCAKeyPath  = "/etc/openvpn/server/ca.key"
	ServerCertPath   = "/etc/openvpn/server/server.crt"
	ServerKeyPath    = "/etc/openvpn/server/server.key"
	ServerDHPath     = "/etc/openvpn/server/dh.pem"
	ServerTLSKeyPath = "/etc/openvpn/server/tls-auth.key"

	// management 接口的密码文件：让 `management 127.0.0.1 7505 <file>` 带口令，
	// 消除 "Using --management on a TCP port WITHOUT passwords" 警告。
	// 文件首行即口令，OpenVPN 启动时(降权前,root)读取；PauseClient/auth 脚本连接后先发同一口令。
	ServerMgmtPasswordPath = "/etc/openvpn/server/mgmt-pw.txt"

	// Default log paths
	DefaultOpenVPNStatusLogPath = "/etc/openvpn/status.log"
	DefaultOpenVPNLogPath       = "/etc/openvpn/openvpn.log"

	// 服务器 IP 分配文件路径
	ServerIPPPath = "/etc/openvpn/server/ipp.txt"

	// CRL（证书吊销列表）路径：删除用户时吊销其证书并重生成此文件，
	// server.conf 通过 crl-verify 引用它。crl.cnf 是配套的 openssl CA 配置，
	// ca-db/ 是 openssl ca 所需的最小数据库目录（index.txt / crlnumber / newcerts）。
	ServerCRLPath    = "/etc/openvpn/server/crl.pem"
	ServerCRLConfig  = "/etc/openvpn/server/crl.cnf"
	ServerCRLDBDir   = "/etc/openvpn/server/ca-db"

	// 客户端配置目录
	ClientConfigDir = "/etc/openvpn/client"

	// 配置文件路径
	ConfigJSONPath = "/etc/openvpn/server/config.json"

	// Supervisor 配置路径
	SupervisorConfigPath        = "/etc/supervisor/supervisord.conf"
	SupervisorConfDir           = "/etc/supervisor/conf.d"
	SupervisorLogDir            = "/var/log/supervisor"
	SupervisorSocketPath        = "/var/run/supervisor.sock"
	SupervisorPidPath           = "/var/run/supervisord.pid"
	SupervisorOpenVPNConfigPath = "/etc/supervisor/conf.d/openvpn-server.conf"
	SupervisorWebConfigPath     = "/etc/supervisor/conf.d/openvpn-go-api.conf"

	// OpenVPN 默认配置值
	DefaultOpenVPNPort             = 4500
	DefaultOpenVPNProto            = "tcp6"
	DefaultOpenVPNServerNetwork    = "10.8.0.0"
	DefaultOpenVPNServerNetmask    = "255.255.255.0"
	DefaultOpenVPNTLSVersion       = "1.2"
	DefaultOpenVPNTLSKey           = "ta.key"
	DefaultOPENVPN_SERVER_HOSTNAME = "192.168.2.1"
	DefaultOpenVPNManagementPort   = 7505
	DefaultOpenVPNBlacklistFile    = "/etc/openvpn/server/blacklist.txt"
	DefaultOpenVPNSyncCerts        = true
	DefaultOpenVPNUseCRL           = true
	DefaultOpenVPNClientToClient   = false
	DefaultOpenVPNClientConfigDir  = "/etc/openvpn/client"
	DefaultOpenVPNTLSKeyPath       = "/etc/openvpn/server/tls-auth.key"
)

// 默认路由配置
var DefaultOpenVPNRoutes = []string{
	"10.10.100.0 255.255.255.0",
	"10.10.98.0 255.255.255.0",
}

// openssl 扩展文件
var OpenSSLExtFiles = []string{
	"openssl-ca.ext",
	"openssl-server.ext",
	"openssl-client.ext",
}

// 这些文件在初始化时从 <cwd>/file/ 复制到 /etc/openvpn/server/ 并 chmod 755
// （见 cmd/environment.go generateCertificates）。
// tls-verify.sh：按 CN 拉黑的脚本（替代旧的 auth-blacklist.sh）。
// crl.cnf：证书吊销用的 openssl CA 配置。
var BlacklistFile = []string{
	"tls-verify.sh",
	"crl.cnf",
	"blacklist.txt",
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
