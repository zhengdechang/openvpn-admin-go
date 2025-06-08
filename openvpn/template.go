package openvpn

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"openvpn-admin-go/constants"
)

// RenderServerConfig 渲染服务端配置模板
func RenderServerConfig(cfg *Config) (string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 使用相对路径查找模板文件
	templatePath := filepath.Join(wd, "template", "server.conf.j2")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("解析服务端配置模板失败: %v", err)
	}

	data := map[string]interface{}{
		"openvpn_port":            cfg.OpenVPNPort,
		"openvpn_proto":           cfg.OpenVPNProto,
		"openvpn_server_network":  cfg.OpenVPNServerNetwork,
		"openvpn_server_netmask":  cfg.OpenVPNServerNetmask,
		"openvpn_client_to_client": cfg.OpenVPNClientToClient,
		"openvpn_routes":          cfg.OpenVPNRoutes,
		"dns_server_ip":           cfg.DNSServerIP,
		"dns_server_domain":       cfg.DNSServerDomain,
		"openvpn_tls_version":     cfg.OpenVPNTLSVersion,
		"ca_cert_path":            constants.ServerCACertPath,
		"server_cert_path":        constants.ServerCertPath,
		"server_key_path":         constants.ServerKeyPath,
		"dh_path":                 constants.ServerDHPath,
		"ipp_path":                constants.ServerIPPPath,
		"tls_key_path":            constants.ServerTLSKeyPath,
		"status_log_path":         constants.ServerStatusLogPath,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染服务端配置模板失败: %v", err)
	}

	return buf.String(), nil
}

// RenderClientConfig 渲染客户端配置模板
func RenderClientConfig(username string, cfg *Config) (string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 使用相对路径查找模板文件
	templatePath := filepath.Join(wd, "template", "client.ovpn.j2")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("解析客户端配置模板失败: %v", err)
	}

	// 读取证书文件
	caCert, err := os.ReadFile(constants.ServerCACertPath)
	if err != nil {
		return "", fmt.Errorf("读取CA证书失败: %v", err)
	}

	clientCert, err := os.ReadFile(filepath.Join(constants.ClientConfigDir, username+".crt"))
	if err != nil {
		return "", fmt.Errorf("读取客户端证书失败: %v", err)
	}

	clientKey, err := os.ReadFile(filepath.Join(constants.ClientConfigDir, username+".key"))
	if err != nil {
		return "", fmt.Errorf("读取客户端密钥失败: %v", err)
	}

	tlsAuthKey, err := os.ReadFile(constants.ServerTLSKeyPath)
	if err != nil {
		return "", fmt.Errorf("读取TLS密钥失败: %v", err)
	}

	// 确保proto值正确传递
	proto := cfg.OpenVPNProto
	if proto == "tcp6" {
		proto = "tcp"
	} else if proto == "udp6" {
		proto = "udp"
	} else if proto == "udp" {
		proto = "udp"
	} else {
		proto = "tcp"
	}

	data := map[string]interface{}{
		"openvpn_proto":           proto,
		"openvpn_port":            cfg.OpenVPNPort,
		"openvpn_server_hostname": cfg.OpenVPNServerHostname,
		"openvpn_tls_version":     cfg.OpenVPNTLSVersion,
		"openvpn_routes":          cfg.OpenVPNRoutes,
		"ca_cert":                 string(caCert),
		"client_cert":             string(clientCert),
		"client_key":              string(clientKey),
		"tls_auth_key":            string(tlsAuthKey),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染客户端配置模板失败: %v", err)
	}

	return buf.String(), nil
} 