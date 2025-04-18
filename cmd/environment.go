package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/openvpn"
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

	// 检查服务器目录下的证书文件
	serverDir := filepath.Dir(constants.ServerConfigPath)
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

func checkCertificates() error {
	// 检查 /etc/openvpn/server 目录下的证书文件
	serverDir := filepath.Dir(constants.ServerConfigPath)
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		return fmt.Errorf("OpenVPN服务器目录不存在: %s", serverDir)
	}

	// 检查必要的证书文件
	requiredFiles := []string{
		constants.ServerCACertPath,
		constants.ServerCAKeyPath,
		constants.ServerCertPath,
		constants.ServerKeyPath,
		constants.ServerDHPath,
		constants.ServerTLSKeyPath,
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("证书文件不存在: %s", file)
		}
	}

	return nil
}

func generateCertificates() error {
	// 创建服务器目录
	serverDir := filepath.Dir(constants.ServerConfigPath)
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建服务器目录失败: %v", err)
	}

	// 生成DH参数
	cmd := exec.Command("openssl", "dhparam", "-out", constants.ServerDHPath, "2048")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成DH参数失败: %v\n输出: %s", err, string(output))
	}

	// 生成CA证书
	cmd = exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-keyout", constants.ServerCAKeyPath, "-out", constants.ServerCACertPath, "-days", "3650", "-nodes", "-subj", "/CN=OpenVPN-CA")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成CA证书失败: %v\n输出: %s", err, string(output))
	}

	// 生成服务器证书
	cmd = exec.Command("openssl", "req", "-newkey", "rsa:2048", "-keyout", constants.ServerKeyPath, "-out", filepath.Join(serverDir, "server.csr"), "-nodes", "-subj", "/CN=OpenVPN-Server")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成服务器证书请求失败: %v\n输出: %s", err, string(output))
	}

	// 签名服务器证书
	cmd = exec.Command("openssl", "x509", "-req", "-in", filepath.Join(serverDir, "server.csr"), "-CA", constants.ServerCACertPath, "-CAkey", constants.ServerCAKeyPath, "-CAcreateserial", "-out", constants.ServerCertPath, "-days", "3650", "-extfile", filepath.Join(serverDir, "openssl-server.ext"))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("签名服务器证书失败: %v\n输出: %s", err, string(output))
	}

	// 生成TLS密钥
	cmd = exec.Command("openvpn", "--genkey", "secret", constants.ServerTLSKeyPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成TLS密钥失败: %v\n输出: %s", err, string(output))
	}

	// 清理临时文件
	os.Remove(filepath.Join(serverDir, "server.csr"))
	os.Remove(filepath.Join(serverDir, "ca.srl"))

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
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 生成OpenVPN服务端配置文件内容
	configContent, err := cfg.GenerateServerConfig()
	if err != nil {
		return fmt.Errorf("生成服务器配置失败: %v", err)
	}

	// 写入配置文件到服务器目录
	if err := os.WriteFile(constants.ServerConfigPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("生成OpenVPN配置文件失败: %v", err)
	}

	fmt.Println("OpenVPN配置文件生成完成")
	return nil
}

// WriteConfig 写入配置文件到服务器目录
func WriteConfig(config string) error {
	// 写入配置文件到服务器目录
	configPath := constants.ServerConfigPath
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	return nil
}

func checkOpenVPNDirectory() error {
	// 检查OpenVPN目录
	openvpnDir := filepath.Dir(constants.ServerConfigPath)
	if _, err := os.Stat(openvpnDir); os.IsNotExist(err) {
		return fmt.Errorf("OpenVPN目录不存在: %s", openvpnDir)
	}

	// 检查服务器目录
	serverDir := filepath.Dir(constants.ServerConfigPath)
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		return fmt.Errorf("OpenVPN服务器目录不存在: %s", serverDir)
	}

	// 检查客户端目录
	if _, err := os.Stat(constants.ClientConfigDir); os.IsNotExist(err) {
		return fmt.Errorf("OpenVPN客户端目录不存在: %s", constants.ClientConfigDir)
	}

	return nil
}

func createOpenVPNDirectory() error {
	// 创建OpenVPN目录
	openvpnDir := filepath.Dir(constants.ServerConfigPath)
	if err := os.MkdirAll(openvpnDir, 0755); err != nil {
		return fmt.Errorf("创建OpenVPN目录失败: %v", err)
	}

	// 创建服务器目录
	serverDir := filepath.Dir(constants.ServerConfigPath)
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建OpenVPN服务器目录失败: %v", err)
	}

	// 创建客户端目录
	if err := os.MkdirAll(constants.ClientConfigDir, 0755); err != nil {
		return fmt.Errorf("创建OpenVPN客户端目录失败: %v", err)
	}

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
	openvpnDir := filepath.Dir(constants.ServerConfigPath)
	if err := os.MkdirAll(openvpnDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 配置目录失败: %v", err)
	}

	// 创建服务器证书目录
	serverDir := filepath.Dir(constants.ServerConfigPath)
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 服务器证书目录失败: %v", err)
	}

	// 创建客户端配置目录
	if err := os.MkdirAll(constants.ClientConfigDir, 0755); err != nil {
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
	cmd = exec.Command("pkill", "openvpn")
	cmd.Run() // 忽略错误，因为可能没有进程在运行

	// 等待进程完全停止
	time.Sleep(2 * time.Second)

	// 启动 OpenVPN 服务
	fmt.Println("正在启动 OpenVPN 服务...")
	cmd = exec.Command("systemctl", "enable", constants.ServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置 OpenVPN 服务开机自启失败: %v", err)
	}

	cmd = exec.Command("systemctl", "start", constants.ServiceName)
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
	if err := checkOpenVPNDirectory(); err != nil {
		return err
	}

	// 检查证书文件
	if err := checkCertificates(); err != nil {
		return err
	}

	// 检查服务端配置文件是否存在
	if _, err := os.Stat(constants.ServerConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("服务端配置文件不存在: %s", constants.ServerConfigPath)
	}

	// 检查CA证书和密钥文件是否存在
	if _, err := os.Stat(constants.ServerCACertPath); os.IsNotExist(err) {
		return fmt.Errorf("CA证书文件不存在: %s", constants.ServerCACertPath)
	}
	if _, err := os.Stat(constants.ServerCAKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("CA密钥文件不存在: %s", constants.ServerCAKeyPath)
	}

	return nil
} 