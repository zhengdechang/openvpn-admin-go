package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"openvpn-admin-go/config"
)

func CheckCertFiles() error {
	// 定义需要检查的证书文件
	certFiles := []string{
		"ca.crt",
		"ca.key",
		"server.crt",
		"server.key",
		"dh.pem",
	}

	// 检查 /etc/openvpn/server 目录下的证书文件
	serverDir := "/etc/openvpn/server"
	missingFiles := []string{}
	for _, file := range certFiles {
		fullPath := filepath.Join(serverDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) > 0 {
		fmt.Println("以下证书文件不存在:")
		for _, file := range missingFiles {
			fmt.Printf("- %s\n", file)
		}
		fmt.Printf("\n请检查路径: %s\n", serverDir)
		return fmt.Errorf("证书文件缺失")
	}

	return nil
}

func generateCertificates() error {
	// 检查证书文件是否已存在
	certFiles := []string{
		"/etc/openvpn/server/ca.crt",
		"/etc/openvpn/server/ca.key",
		"/etc/openvpn/server/server.crt",
		"/etc/openvpn/server/server.key",
		"/etc/openvpn/server/dh.pem",
		"/etc/openvpn/server/ta.key", // 添加TLS密钥检查
	}

	allExist := true
	for _, file := range certFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			allExist = false
			break
		}
	}

	if allExist {
		fmt.Println("证书文件已存在，跳过生成")
		return nil
	}

	// 创建服务器目录
	serverDir := "/etc/openvpn/server"
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建服务器目录失败: %v", err)
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	// 复制扩展文件
	extFiles := []string{
		"openssl-ca.ext",
		"openssl-server.ext",
		"openssl-client.ext",
	}

	for _, file := range extFiles {
		src := filepath.Join(currentDir, "file", file)
		dst := filepath.Join(serverDir, file)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("复制扩展文件 %s 失败: %v", file, err)
		}
		fmt.Printf("已复制扩展文件: %s\n", file)
	}

	// 生成DH参数
	fmt.Println("正在生成DH参数...")
	cmd := exec.Command("openssl", "dhparam", "-out", "/etc/openvpn/server/dh.pem", "2048")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成DH参数失败: %v", err)
	}
	fmt.Println("DH参数生成成功")

	// 生成CA证书
	fmt.Println("正在生成CA证书...")
	cmd = exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-keyout", "/etc/openvpn/server/ca.key", "-out", "/etc/openvpn/server/ca.crt", "-days", "3650", "-nodes", "-subj", "/CN=OpenVPN-CA")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成CA证书失败: %v", err)
	}
	fmt.Println("CA证书生成成功")

	// 生成服务器证书
	fmt.Println("正在生成服务器证书...")
	cmd = exec.Command("openssl", "req", "-newkey", "rsa:2048", "-keyout", "/etc/openvpn/server/server.key", "-out", "/etc/openvpn/server/server.csr", "-nodes", "-subj", "/CN=OpenVPN-Server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成服务器证书请求失败: %v", err)
	}

	cmd = exec.Command("openssl", "x509", "-req", "-in", "/etc/openvpn/server/server.csr", "-CA", "/etc/openvpn/server/ca.crt", "-CAkey", "/etc/openvpn/server/ca.key", "-CAcreateserial", "-out", "/etc/openvpn/server/server.crt", "-days", "3650", "-extfile", "/etc/openvpn/server/openssl-server.ext")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("签名服务器证书失败: %v", err)
	}
	fmt.Println("服务器证书生成成功")

	// 生成TLS密钥
	fmt.Println("正在生成TLS密钥...")
	cmd = exec.Command("openvpn", "--genkey", "secret", "/etc/openvpn/server/ta.key")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成TLS密钥失败: %v", err)
	}
	fmt.Println("TLS密钥生成成功")

	// 清理临时文件
	os.Remove("/etc/openvpn/server/server.csr")
	os.Remove("/etc/openvpn/server/ca.srl")

	return nil
}

// copyFile 复制文件
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

func generateOpenVPNConfig() error {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 生成OpenVPN服务端配置文件内容
	configContent := cfg.GenerateServerConfig()

	// 写入配置文件到 /etc/openvpn 目录
	configPath := filepath.Join("/etc/openvpn", "server.conf")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("生成OpenVPN配置文件失败: %v", err)
	}

	fmt.Println("OpenVPN配置文件生成完成")
	return nil
}

func InstallEnvironment() error {
	// 检查是否以root权限运行
	if os.Geteuid() != 0 {
		return fmt.Errorf("请使用 sudo 运行程序")
	}

	fmt.Println("开始安装所需环境...")

	// 安装 OpenVPN 和 OpenSSL
	fmt.Println("正在安装 OpenVPN 和 OpenSSL...")
	cmd := exec.Command("apt-get", "update")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("更新软件包列表失败: %v", err)
	}

	cmd = exec.Command("apt-get", "install", "-y", "openvpn", "openssl")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("安装 OpenVPN 和 OpenSSL 失败: %v", err)
	}

	// 创建配置目录
	fmt.Println("正在创建配置目录...")
	openvpnDir := "/etc/openvpn"
	if err := os.MkdirAll(openvpnDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 配置目录失败: %v", err)
	}

	// 创建服务器证书目录
	serverDir := "/etc/openvpn/server"
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 服务器证书目录失败: %v", err)
	}

	// 创建客户端配置目录
	clientDir := "/etc/openvpn/client"
	if err := os.MkdirAll(clientDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 客户端目录失败: %v", err)
	}

	// 生成证书
	fmt.Println("正在生成证书...")
	if err := generateCertificates(); err != nil {
		return fmt.Errorf("生成证书失败: %v", err)
	}
	
	// 生成Openvpn配置文件
	fmt.Println("正在生成Openvpn配置文件...")
	if err := generateOpenVPNConfig(); err != nil {
		return fmt.Errorf("生成OpenVPN配置文件失败: %v", err)
	}

	// 停止所有正在运行的 OpenVPN 进程
	fmt.Println("正在停止所有 OpenVPN 进程...")
	cmd = exec.Command("systemctl", "stop", "openvpn")
	cmd.Run() // 忽略错误，因为服务可能没有运行
	cmd = exec.Command("systemctl", "stop", "openvpn-server@server")
	cmd.Run() // 忽略错误，因为服务可能没有运行
	cmd = exec.Command("pkill", "openvpn")
	cmd.Run() // 忽略错误，因为可能没有进程在运行

	// 等待进程完全停止
	time.Sleep(2 * time.Second)

	// 启动 OpenVPN 服务
	fmt.Println("正在启动 OpenVPN 服务...")
	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新加载服务配置失败: %v", err)
	}

	// 创建自定义服务单元
	serviceContent := `[Unit]
Description=OpenVPN service
After=network.target

[Service]
Type=simple
WorkingDirectory=/etc/openvpn
ExecStart=/usr/sbin/openvpn --config server.conf
Restart=on-failure

[Install]
WantedBy=multi-user.target
`
	if err := os.WriteFile("/etc/systemd/system/openvpn.service", []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("创建服务单元文件失败: %v", err)
	}

	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新加载服务配置失败: %v", err)
	}

	cmd = exec.Command("systemctl", "enable", "openvpn")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置 OpenVPN 服务开机自启失败: %v", err)
	}

	cmd = exec.Command("systemctl", "start", "openvpn")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启动 OpenVPN 服务失败: %v", err)
	}

	// 等待服务启动
	fmt.Println("等待服务启动...")
	time.Sleep(5 * time.Second)

	fmt.Println("环境安装完成")
	return nil
}

func CheckEnvironment() error {
	// 检查OpenVPN是否安装
	if _, err := exec.LookPath("openvpn"); err != nil {
		return fmt.Errorf("未找到OpenVPN，请先安装OpenVPN")
	}

	// 检查OpenSSL是否安装
	if _, err := exec.LookPath("openssl"); err != nil {
		return fmt.Errorf("未找到OpenSSL，请先安装OpenSSL")
	}

	// 检查OpenVPN配置目录是否存在
	if _, err := os.Stat("/etc/openvpn"); os.IsNotExist(err) {
		return fmt.Errorf("OpenVPN配置目录不存在，请先安装OpenVPN")
	}

	// 检查证书文件
	if err := CheckCertFiles(); err != nil {
		return err
	}

	return nil
} 