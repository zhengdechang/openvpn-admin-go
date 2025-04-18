package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"openvpn-admin-go/config"
)

// CreateClient 创建新的OpenVPN客户端
func CreateClient(username string) error {
	fmt.Printf("开始创建客户端: %s\n", username)
	
	// 检查证书目录
	certDir := "/etc/openvpn/client"
	fmt.Printf("检查证书目录: %s\n", certDir)
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		fmt.Printf("创建证书目录: %s\n", certDir)
		cmd := exec.Command("sudo", "mkdir", "-p", certDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("创建证书目录失败: %v, 输出: %s", err, string(output))
		}
		fmt.Println("证书目录创建成功")
	}

	// 检查并生成TLS密钥
	tlsKeyPath := "/etc/openvpn/server/ta.key"
	if _, err := os.Stat(tlsKeyPath); os.IsNotExist(err) {
		fmt.Printf("正在生成TLS密钥: %s\n", tlsKeyPath)
		cmd := exec.Command("sudo", "openvpn", "--genkey", "secret", tlsKeyPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("生成TLS密钥失败: %v, 输出: %s", err, string(output))
		}
		fmt.Println("TLS密钥生成成功")
	}

	// 检查CA证书和密钥是否存在
	caCertPath := "/etc/openvpn/server/ca.crt"
	caKeyPath := "/etc/openvpn/server/ca.key"
	
	fmt.Printf("检查CA证书: %s\n", caCertPath)
	fmt.Printf("检查CA密钥: %s\n", caKeyPath)
	
	if _, err := os.Stat(caCertPath); os.IsNotExist(err) {
		return fmt.Errorf("CA证书不存在: %s", caCertPath)
	}
	if _, err := os.Stat(caKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("CA密钥不存在: %s", caKeyPath)
	}
	
	fmt.Printf("使用CA证书: %s\n", caCertPath)
	fmt.Printf("使用CA密钥: %s\n", caKeyPath)
	fmt.Println("CA证书和密钥检查通过")

	// 检查客户端扩展文件
	clientExtFile := "/etc/openvpn/server/openssl-client.ext"
	if _, err := os.Stat(clientExtFile); os.IsNotExist(err) {
		return fmt.Errorf("客户端扩展文件不存在: %s", clientExtFile)
	}
	fmt.Printf("使用客户端扩展文件: %s\n", clientExtFile)

	// 生成客户端私钥
	fmt.Printf("正在为客户端 %s 生成私钥...\n", username)
	keyPath := certDir + "/" + username + ".key"
	cmd := exec.Command("sudo", "openssl", "genrsa", "-out", keyPath, "2048")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("生成私钥失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("私钥生成成功")

	// 生成证书签名请求
	fmt.Printf("正在为客户端 %s 生成证书签名请求...\n", username)
	csrPath := certDir + "/" + username + ".csr"
	cmd = exec.Command("sudo", "openssl", "req", "-new", "-key", keyPath, "-out", csrPath, "-subj", "/CN="+username)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("生成证书签名请求失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("证书签名请求生成成功")

	// 使用CA证书签名
	fmt.Printf("正在为客户端 %s 签名证书...\n", username)
	crtPath := certDir + "/" + username + ".crt"
	cmd = exec.Command("sudo", "openssl", "x509", "-req", "-in", csrPath, "-CA", caCertPath, "-CAkey", caKeyPath, "-CAcreateserial", "-out", crtPath, "-days", "3650", "-extfile", clientExtFile)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("签名证书失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("证书签名成功")

	// 复制CA证书到客户端目录
	clientCaPath := certDir + "/ca.crt"
	fmt.Printf("正在复制CA证书到: %s\n", clientCaPath)
	cmd = exec.Command("sudo", "cp", caCertPath, clientCaPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("复制CA证书失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("CA证书复制成功")

	// 清理临时文件
	fmt.Printf("清理临时文件: %s\n", csrPath)
	cmd = exec.Command("sudo", "rm", csrPath)
	cmd.Run() // 忽略错误，因为文件可能不存在

	// 生成.ovpn配置文件
	fmt.Printf("正在为客户端 %s 生成配置文件...\n", username)
	ovpnPath := certDir + "/" + username + ".ovpn"
	
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	
	// 生成客户端配置
	clientConfig, err := GenerateClientConfig(username, cfg)
	if err != nil {
		return fmt.Errorf("生成客户端配置失败: %v", err)
	}

	// 使用sudo写入配置文件
	tempFile := "/tmp/" + username + ".ovpn"
	fmt.Printf("创建临时配置文件: %s\n", tempFile)
	if err := os.WriteFile(tempFile, []byte(clientConfig), 0644); err != nil {
		return fmt.Errorf("创建临时配置文件失败: %v", err)
	}
	
	fmt.Printf("移动配置文件到: %s\n", ovpnPath)
	cmd = exec.Command("sudo", "mv", tempFile, ovpnPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("移动配置文件失败: %v, 输出: %s", err, string(output))
	}

	// 设置文件权限
	fmt.Printf("设置文件权限: %s\n", ovpnPath)
	cmd = exec.Command("sudo", "chmod", "644", ovpnPath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("设置文件权限失败: %v, 输出: %s", err, string(output))
	}

	fmt.Printf("客户端 %s 的证书和配置文件已生成并复制到 %s 目录\n", username, certDir)
	return nil
}

// readFile 读取文件内容
func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

// DeleteClient 删除OpenVPN客户端
func DeleteClient(username string) error {
	certDir := "/etc/openvpn/client"
	files := []string{
		username + ".key",
		username + ".crt",
		username + ".ovpn",
	}

	for _, file := range files {
		path := certDir + "/" + file
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("删除文件 %s 失败: %v", file, err)
		}
	}

	return nil
}

// PauseClient 暂停OpenVPN客户端
func PauseClient(username string) error {
	// TODO: 实现暂停客户端的功能
	return nil
}

// ResumeClient 恢复OpenVPN客户端
func ResumeClient(username string) error {
	// TODO: 实现恢复客户端的功能
	return nil
}

// GetClientStatus 获取OpenVPN客户端状态
func GetClientStatus(username string) (*ClientStatus, error) {
	// 检查客户端配置文件是否存在
	certDir := "/etc/openvpn/client"
	ovpnPath := certDir + "/" + username + ".ovpn"
	
	// 获取文件信息
	fileInfo, err := os.Stat(ovpnPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("客户端 %s 不存在", username)
	}
	if err != nil {
		return nil, fmt.Errorf("获取客户端信息失败: %v", err)
	}

	// 创建状态对象，使用文件创建时间作为客户端创建时间
	status := &ClientStatus{
		Username:    username,
		ConnectedAt: fileInfo.ModTime(), // 使用文件修改时间作为创建时间
		IsPaused:    false,
	}

	// 检查OpenVPN状态日志获取连接状态
	statusLog := "/var/log/openvpn/status.log"
	if logContent, err := os.ReadFile(statusLog); err == nil {
		lines := strings.Split(string(logContent), "\n")
		for _, line := range lines {
			// 查找客户端连接信息
			if strings.Contains(line, username) {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					// 解析连接时间
					if t, err := time.Parse("2006-01-02 15:04:05", fields[1]); err == nil {
						status.ConnectedAt = t
					}
				}
			}
		}
	}

	return status, nil
}

// GetAllClientStatuses 获取所有OpenVPN客户端状态
func GetAllClientStatuses() ([]ClientStatus, error) {
	certDir := "/etc/openvpn/client"
	var statuses []ClientStatus

	// 读取客户端目录
	files, err := os.ReadDir(certDir)
	if err != nil {
		return nil, fmt.Errorf("读取客户端目录失败: %v", err)
	}

	// 遍历所有.ovpn文件
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".ovpn") {
			username := strings.TrimSuffix(file.Name(), ".ovpn")
			if status, err := GetClientStatus(username); err == nil {
				statuses = append(statuses, *status)
			}
		}
	}

	return statuses, nil
}

// ClientStatus 客户端状态
type ClientStatus struct {
	Username      string
	ConnectedAt   time.Time
	Disconnected  time.Time
	IsPaused      bool
}

// GenerateClientConfig 生成客户端配置文件内容
func GenerateClientConfig(username string, cfg *config.Config) (string, error) {
	// 读取服务器配置以获取端口和协议
	serverConfig, err := os.ReadFile("/etc/openvpn/server/server.conf")
	if err != nil {
		return "", fmt.Errorf("读取服务器配置失败: %v", err)
	}

	// 解析端口和协议
	var port, proto string
	lines := strings.Split(string(serverConfig), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "port ") {
			port = strings.TrimSpace(strings.TrimPrefix(line, "port "))
		}
		if strings.HasPrefix(line, "proto ") {
			proto = strings.TrimSpace(strings.TrimPrefix(line, "proto "))
		}
	}

	// 读取客户端证书和密钥
	certPath := filepath.Join("/etc/openvpn/client", username+".crt")
	keyPath := filepath.Join("/etc/openvpn/client", username+".key")
	
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return "", fmt.Errorf("读取客户端证书失败: %v", err)
	}
	
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("读取客户端密钥失败: %v", err)
	}

	// 读取CA证书
	caCert, err := os.ReadFile("/etc/openvpn/server/ca.crt")
	if err != nil {
		return "", fmt.Errorf("读取CA证书失败: %v", err)
	}

	// 读取TLS密钥
	tlsKey, err := os.ReadFile(cfg.OpenVPNTLSKeyPath)
	if err != nil {
		return "", fmt.Errorf("读取TLS密钥失败: %v", err)
	}

	// 生成客户端配置
	config := fmt.Sprintf(`client
dev tun
proto %s
remote %s %s
resolv-retry 5
nobind
persist-key
persist-tun
remote-cert-tls server
cipher AES-256-GCM
auth SHA256
key-direction 1
tls-client
tls-version-min %s
tls-cipher TLS-ECDHE-ECDSA-WITH-AES-256-GCM-SHA384:TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384:TLS-ECDHE-ECDSA-WITH-AES-128-GCM-SHA256:TLS-ECDHE-RSA-WITH-AES-128-GCM-SHA256
auth-nocache
<ca>
%s
</ca>
<cert>
%s
</cert>
<key>
%s
</key>
<tls-auth>
%s
</tls-auth>
`, 
		proto,
		cfg.OpenVPNServerHostname,
		port,
		cfg.OpenVPNTLSVersion,
		caCert,
		cert,
		key,
		tlsKey,
	)

	return config, nil
}
