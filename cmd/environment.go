package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/utils"
)

// ErrRootRequired 表示自动安装环境需要 root 权限
var ErrRootRequired = errors.New("自动安装需要 root 权限")

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

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	for _, file := range constants.OpenSSLExtFiles {
		src := filepath.Join(currentDir, "file", file)
		dst := filepath.Join(serverDir, file)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("复制扩展文件 %s 失败: %v", file, err)
		}
		fmt.Printf("已复制扩展文件: %s\n", file)
	}

	for _, file := range constants.BlacklistFile {
		src := filepath.Join(currentDir, "file", file)
		dst := filepath.Join(serverDir, file)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("复制文件 %s 失败: %v", file, err)
		}
		// 设置文件权限为 644 (rw-r--r--)
		if err := os.Chmod(dst, 777); err != nil {
			return fmt.Errorf("设置文件权限失败 %s: %v", dst, err)
		}
		fmt.Printf("已复制文件: %s\n", file)
	}

	// 确保 ccd 目录存在并设置正确的权限
	ccdDir := filepath.Join(constants.ClientConfigDir, "ccd")
	if err := os.MkdirAll(ccdDir, 0755); err != nil {
		return fmt.Errorf("创建 ccd 目录失败: %v", err)
	}
	// 设置 ccd 目录权限为 755 (rwxr-xr-x)
	if err := os.Chmod(ccdDir, 0755); err != nil {
		return fmt.Errorf("设置 ccd 目录权限失败: %v", err)
	}

	// 生成DH参数
	if err := utils.ExecCommand(fmt.Sprintf("openssl dhparam -out %s 2048", constants.ServerDHPath)); err != nil {
		return fmt.Errorf("生成DH参数失败: %v", err)
	}

	// 生成CA证书
	if err := utils.ExecCommand(fmt.Sprintf("openssl req -x509 -newkey rsa:2048 -keyout %s -out %s -days 3650 -nodes -subj '/CN=OpenVPN-CA'", constants.ServerCAKeyPath, constants.ServerCACertPath)); err != nil {
		return fmt.Errorf("生成CA证书失败: %v", err)
	}

	// 生成服务器证书
	if err := utils.ExecCommand(fmt.Sprintf("openssl req -newkey rsa:2048 -keyout %s -out %s -nodes -subj '/CN=OpenVPN-Server'", constants.ServerKeyPath, filepath.Join(serverDir, "server.csr"))); err != nil {
		return fmt.Errorf("生成服务器证书请求失败: %v", err)
	}

	// 签名服务器证书
	if err := utils.ExecCommand(fmt.Sprintf("openssl x509 -req -in %s -CA %s -CAkey %s -CAcreateserial -out %s -days 3650 -extfile %s", filepath.Join(serverDir, "server.csr"), constants.ServerCACertPath, constants.ServerCAKeyPath, constants.ServerCertPath, filepath.Join(serverDir, "openssl-server.ext"))); err != nil {
		return fmt.Errorf("签名服务器证书失败: %v", err)
	}

	// 生成TLS密钥
	if err := utils.ExecCommand(fmt.Sprintf("openvpn --genkey secret %s", constants.ServerTLSKeyPath)); err != nil {
		return fmt.Errorf("生成TLS密钥失败: %v", err)
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
		return ErrRootRequired
	}

	fmt.Println("开始安装所需环境...")

	// 安装 OpenVPN、OpenSSL 和 Supervisor
	fmt.Println("正在安装 OpenVPN、OpenSSL 和 Supervisor...")
	if err := utils.ExecCommand("apt-get update"); err != nil {
		return fmt.Errorf("更新软件包列表失败: %v", err)
	}

	if err := utils.ExecCommand("apt-get install -y openvpn openssl supervisor"); err != nil {
		return fmt.Errorf("安装 OpenVPN、OpenSSL 或 Supervisor 失败: %v", err)
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
	utils.ExecCommand("pkill openvpn") // 忽略错误，因为可能没有进程在运行

	// 等待进程完全停止
	time.Sleep(2 * time.Second)

	// 启动 OpenVPN 服务
	fmt.Println("正在启动 OpenVPN 服务...")
	utils.SystemctlEnable(constants.ServiceName)
	utils.SystemctlStart(constants.ServiceName)

	// 等待服务启动
	fmt.Println("等待服务启动...")
	time.Sleep(5 * time.Second)

	fmt.Println("环境安装完成")
	return nil
}

func CheckEnvironment() error {
	// 检查OpenVPN是否安装
	if !utils.CheckCommandExists("openvpn") {
		return fmt.Errorf("未找到OpenVPN，请先安装OpenVPN")
	}

	// 检查OpenSSL是否安装
	if !utils.CheckCommandExists("openssl") {
		return fmt.Errorf("未找到OpenSSL，请先安装OpenSSL")
	}

	// 检查supervisor是否安装
	if !utils.CheckSupervisorInstalled() {
		return fmt.Errorf("未找到supervisor，请先安装supervisor")
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

// RunEnvironmentSetup 交互式执行环境检查与安装
func RunEnvironmentSetup() {
	fmt.Println("\n=== 运行环境检查 ===")
	if err := CheckEnvironment(); err != nil {
		fmt.Printf("检测到环境缺失: %v\n", err)
		fmt.Println("正在尝试自动安装必要组件...")
		if errInstall := InstallEnvironment(); errInstall != nil {
			fmt.Printf("自动安装失败: %v\n", errInstall)
			fmt.Println("请手动检查系统依赖或日志后重试。")
			fmt.Println("\n按回车键返回主菜单...")
			fmt.Scanln()
			return
		}

		// 安装完成后重新验证
		if errRecheck := CheckEnvironment(); errRecheck != nil {
			fmt.Printf("环境仍然存在问题: %v\n", errRecheck)
			fmt.Println("请手动排查后重新运行。")
		} else {
			fmt.Println("环境安装并验证完成。")
		}
	} else {
		fmt.Println("环境已就绪，无需额外操作。")
	}

	fmt.Println("\n按回车键返回主菜单...")
	fmt.Scanln()
}
