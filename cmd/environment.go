package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func CheckCertFiles() error {
	// 定义可能的证书文件路径
	possiblePaths := []string{
		"/etc/openvpn",
		"/etc/openvpn/server",
	}

	// 定义需要检查的证书文件
	certFiles := []string{
		"ca.crt",
		"ca.key",
		"client.key",
		"dh.pem",
		"server.crt",
		"server.csr",
		"server.key",
	}

	// 检查每个文件是否存在于任一路径中
	missingFiles := []string{}
	for _, file := range certFiles {
		found := false
		for _, path := range possiblePaths {
			fullPath := filepath.Join(path, file)
			if _, err := os.Stat(fullPath); err == nil {
				found = true
				break
			}
		}
		if !found {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) > 0 {
		fmt.Println("以下证书文件不存在:")
		for _, file := range missingFiles {
			fmt.Printf("- %s\n", file)
		}
		fmt.Println("\n请检查以下路径:")
		for _, path := range possiblePaths {
			fmt.Printf("- %s\n", path)
		}
		return fmt.Errorf("证书文件缺失")
	}

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

func InstallEnvironment() error {
	// 检查是否以root权限运行
	if os.Geteuid() != 0 {
		// 尝试使用sudo重新运行程序
		cmd := exec.Command("sudo", os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("获取root权限失败: %v\n请使用 sudo 运行程序", err)
		}
		os.Exit(0)
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
	configDir := "/etc/openvpn"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 配置目录失败: %v", err)
	}

	serverDir := "/etc/openvpn/server"
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 服务器目录失败: %v", err)
	}

	// 复制证书文件
	fmt.Println("正在复制证书文件...")
	certFiles := []string{
		"ca.crt",
		"ca.key",
		"client.key",
		"dh.pem",
		"server.crt",
		"server.csr",
		"server.key",
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	for _, file := range certFiles {
		source := filepath.Join(currentDir, "file", file)
		dest := filepath.Join(serverDir, file)

		// 检查源文件是否存在
		if _, err := os.Stat(source); os.IsNotExist(err) {
			return fmt.Errorf("证书文件不存在: %s\n请确保file目录下包含所有必要的证书文件", source)
		}

		// 复制文件
		cmd = exec.Command("cp", source, dest)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("复制证书文件失败: %s -> %s: %v", source, dest, err)
		}

		// 设置文件权限
		cmd = exec.Command("chmod", "600", dest)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("设置证书文件权限失败: %s: %v", dest, err)
		}

		fmt.Printf("已复制并设置权限: %s\n", file)
	}

	// 启动 OpenVPN 服务
	fmt.Println("正在启动 OpenVPN 服务...")
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

	// 验证安装
	fmt.Println("验证安装...")
	if err := verifyInstallation(); err != nil {
		return fmt.Errorf("安装验证失败: %v", err)
	}

	fmt.Println("环境安装完成")
	return nil
}

func verifyInstallation() error {
	// 检查 OpenVPN 是否安装
	if _, err := exec.LookPath("openvpn"); err != nil {
		return fmt.Errorf("OpenVPN 未正确安装: %v", err)
	}

	// 检查 OpenSSL 是否安装
	if _, err := exec.LookPath("openssl"); err != nil {
		return fmt.Errorf("OpenSSL 未正确安装: %v", err)
	}

	// 检查配置目录
	configDir := "/etc/openvpn"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("OpenVPN 配置目录不存在: %s", configDir)
	}

	// 检查 OpenVPN 服务是否运行
	cmd := exec.Command("systemctl", "is-active", "openvpn")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("OpenVPN 服务未运行: %v", err)
	}

	return nil
} 