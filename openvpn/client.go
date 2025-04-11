package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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

	// 检查CA证书和密钥是否存在
	caCertPath := "/etc/openvpn/ca.crt"
	caKeyPath := "/etc/openvpn/ca.key"
	serverCaCertPath := "/etc/openvpn/server/ca.crt"
	serverCaKeyPath := "/etc/openvpn/server/ca.key"
	
	fmt.Printf("检查CA证书: %s\n", caCertPath)
	fmt.Printf("检查CA密钥: %s\n", caKeyPath)
	
	// 检查主目录下的文件
	caCertExists := false
	caKeyExists := false
	
	if _, err := os.Stat(caCertPath); err == nil {
		caCertExists = true
	}
	if _, err := os.Stat(caKeyPath); err == nil {
		caKeyExists = true
	}
	
	// 如果主目录下没有，检查server目录
	if !caCertExists || !caKeyExists {
		fmt.Printf("检查server目录下的CA证书: %s\n", serverCaCertPath)
		fmt.Printf("检查server目录下的CA密钥: %s\n", serverCaKeyPath)
		
		if _, err := os.Stat(serverCaCertPath); err == nil {
			caCertPath = serverCaCertPath
			caCertExists = true
		}
		if _, err := os.Stat(serverCaKeyPath); err == nil {
			caKeyPath = serverCaKeyPath
			caKeyExists = true
		}
	}
	
	if !caCertExists {
		return fmt.Errorf("CA证书不存在: %s 或 %s", "/etc/openvpn/ca.crt", "/etc/openvpn/server/ca.crt")
	}
	if !caKeyExists {
		return fmt.Errorf("CA密钥不存在: %s 或 %s", "/etc/openvpn/ca.key", "/etc/openvpn/server/ca.key")
	}
	
	fmt.Printf("使用CA证书: %s\n", caCertPath)
	fmt.Printf("使用CA密钥: %s\n", caKeyPath)
	fmt.Println("CA证书和密钥检查通过")

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
	cmd = exec.Command("sudo", "openssl", "x509", "-req", "-in", csrPath, "-CA", caCertPath, "-CAkey", caKeyPath, "-CAcreateserial", "-out", crtPath, "-days", "3650")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("签名证书失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("证书签名成功")

	// 复制CA证书
	caPath := certDir + "/ca.crt"
	fmt.Printf("正在复制CA证书到: %s\n", caPath)
	cmd = exec.Command("sudo", "cp", caCertPath, caPath)
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
	
	// 读取文件内容
	fmt.Println("正在读取证书文件内容...")
	caContent := readFile(caPath)
	crtContent := readFile(crtPath)
	keyContent := readFile(keyPath)
	
	if caContent == "" || crtContent == "" || keyContent == "" {
		return fmt.Errorf("读取证书文件失败")
	}
	fmt.Println("证书文件内容读取成功")
	
	// 读取服务器配置文件获取路由信息
	serverConfig := readFile("/etc/openvpn/server.conf")
	var routes []string
	if serverConfig != "" {
		lines := strings.Split(serverConfig, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "push \"route ") {
				// 提取路由信息
				route := strings.TrimPrefix(line, "push \"route ")
				route = strings.TrimSuffix(route, "\"")
				routes = append(routes, route)
			}
		}
	}
	
	// 如果没有找到路由配置，使用默认路由
	if len(routes) == 0 {
		routes = []string{"10.8.0.0 255.255.255.0"}
	}
	
	// 生成配置文件
	config := fmt.Sprintf(`client
dev tun
proto udp
remote %s 1194
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
cipher AES-256-CBC
verb 3
`, getServerIP())

	// 添加路由配置
	for _, route := range routes {
		config += fmt.Sprintf("route %s\n", route)
	}

	// 添加证书和密钥
	config += fmt.Sprintf(`
<ca>
%s
</ca>
<cert>
%s
</cert>
<key>
%s
</key>
`, caContent, crtContent, keyContent)

	// 使用sudo写入配置文件
	tempFile := "/tmp/" + username + ".ovpn"
	fmt.Printf("创建临时配置文件: %s\n", tempFile)
	if err := os.WriteFile(tempFile, []byte(config), 0644); err != nil {
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
	// TODO: 实现获取客户端状态的功能
	return &ClientStatus{
		Username:    username,
		ConnectedAt: time.Now(),
		IsPaused:    false,
	}, nil
}

// GetAllClientStatuses 获取所有OpenVPN客户端状态
func GetAllClientStatuses() ([]ClientStatus, error) {
	certDir := "/etc/openvpn/client"
	
	// 检查目录是否存在
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("客户端目录不存在: %s", certDir)
	}
	
	// 读取目录中的所有文件
	files, err := os.ReadDir(certDir)
	if err != nil {
		return nil, fmt.Errorf("读取客户端目录失败: %v", err)
	}
	
	var clients []ClientStatus
	
	// 遍历所有文件，查找.ovpn文件
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".ovpn") {
			// 从文件名中提取用户名（去掉.ovpn后缀）
			username := strings.TrimSuffix(file.Name(), ".ovpn")
			
			// 检查对应的证书文件是否存在
			crtPath := certDir + "/" + username + ".crt"
			keyPath := certDir + "/" + username + ".key"
			
			// 检查证书和密钥文件是否存在
			_, crtErr := os.Stat(crtPath)
			_, keyErr := os.Stat(keyPath)
			
			if crtErr == nil && keyErr == nil {
				// 获取证书的创建时间
				info, err := os.Stat(crtPath)
				if err != nil {
					continue
				}
				
				clients = append(clients, ClientStatus{
					Username:    username,
					ConnectedAt: info.ModTime(),
					IsPaused:    false,
				})
			}
		}
	}
	
	return clients, nil
}

// ClientStatus 表示OpenVPN客户端的状态
type ClientStatus struct {
	Username      string
	ConnectedAt   time.Time
	Disconnected  time.Time
	IsPaused      bool
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readFile(path string) string {
	fmt.Printf("正在读取文件: %s\n", path)
	cmd := exec.Command("sudo", "cat", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("读取文件 %s 失败: %v, 输出: %s\n", path, err, string(output))
		return ""
	}
	fmt.Printf("成功读取文件: %s, 内容长度: %d\n", path, len(output))
	return string(output)
}

func getServerIP() string {
	// 读取服务器配置文件获取IP
	fmt.Println("正在获取服务器IP...")
	config := readFile("/etc/openvpn/server.conf")
	if config == "" {
		fmt.Println("无法读取服务器配置文件，使用默认IP")
		return "your-server-ip" // 默认值
	}

	lines := strings.Split(config, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "server ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				ip := parts[1]
				fmt.Printf("找到服务器IP: %s\n", ip)
				return ip
			}
		}
	}

	fmt.Println("未找到服务器IP配置，使用默认IP")
	return "your-server-ip" // 默认值
} 