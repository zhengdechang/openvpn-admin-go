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

func generateCertificates(serverDir string) error {
	// 确保目录存在
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建证书目录失败: %v", err)
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	// 复制扩展文件到服务器目录
	fmt.Println("正在复制扩展文件到服务器目录...")
	extFiles := []string{"openssl-ca.ext", "openssl-server.ext", "openssl-client.ext"}
	for _, extFile := range extFiles {
		srcPath := filepath.Join(currentDir, "file", extFile)
		dstPath := filepath.Join(serverDir, extFile)
		
		// 检查源文件是否存在
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return fmt.Errorf("扩展文件 %s 不存在于 %s", extFile, srcPath)
		}
		
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("复制扩展文件 %s 失败: %v", extFile, err)
		}
		
		// 验证文件是否成功复制
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			return fmt.Errorf("扩展文件 %s 复制失败，目标文件不存在", extFile)
		}
		
		fmt.Printf("成功复制扩展文件: %s\n", extFile)
	}

	// 生成 DH 参数
	fmt.Println("正在生成 DH 参数（这可能需要几分钟）...")
	dhPath := filepath.Join(serverDir, "dh.pem")
	cmd := exec.Command("openssl", "dhparam", "-out", dhPath, "2048")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成 DH 参数失败: %v\n输出: %s", err, string(output))
	}
	os.Chmod(dhPath, 0600)

	// 生成 CA 私钥和证书
	fmt.Println("正在生成 CA 证书...")
	caKey := filepath.Join(serverDir, "ca.key")
	caCrt := filepath.Join(serverDir, "ca.crt")
	
	cmd = exec.Command("openssl", "genrsa", "-out", caKey, "2048")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成 CA 私钥失败: %v\n输出: %s", err, string(output))
	}
	os.Chmod(caKey, 0600)

	cmd = exec.Command("openssl", "req", "-new", "-x509", "-days", "3650", "-key", caKey, "-out", caCrt, "-subj", "/CN=OpenVPN-CA")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成 CA 证书失败: %v\n输出: %s", err, string(output))
	}
	os.Chmod(caCrt, 0644)

	// 生成服务器私钥和证书
	fmt.Println("正在生成服务器证书...")
	serverKey := filepath.Join(serverDir, "server.key")
	serverCsr := filepath.Join(serverDir, "server.csr")
	serverCrt := filepath.Join(serverDir, "server.crt")

	cmd = exec.Command("openssl", "genrsa", "-out", serverKey, "2048")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成服务器私钥失败: %v\n输出: %s", err, string(output))
	}
	os.Chmod(serverKey, 0600)

	cmd = exec.Command("openssl", "req", "-new", "-key", serverKey, "-out", serverCsr, "-subj", "/CN=OpenVPN-Server")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成服务器证书请求失败: %v\n输出: %s", err, string(output))
	}

	serverExtFile := filepath.Join(serverDir, "openssl-server.ext")
	cmd = exec.Command("openssl", "x509", "-req", "-days", "3650", "-in", serverCsr, "-CA", caCrt, "-CAkey", caKey, "-CAcreateserial", "-out", serverCrt, "-extfile", serverExtFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("签名服务器证书失败: %v\n输出: %s", err, string(output))
	}
	os.Chmod(serverCrt, 0644)

	// 清理临时文件
	os.Remove(serverCsr)
	os.Remove(filepath.Join(serverDir, "ca.srl"))

	fmt.Println("证书生成完成")
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

func generateOpenVPNConfig(serverDir string) error {
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
	serverDir := "/etc/openvpn/server"
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 服务器目录失败: %v", err)
	}

	// 创建客户端配置目录
	clientDir := "/etc/openvpn/client"
	if err := os.MkdirAll(clientDir, 0755); err != nil {
		return fmt.Errorf("创建 OpenVPN 客户端目录失败: %v", err)
	}

	// 生成证书
	fmt.Println("正在生成证书...")
	if err := generateCertificates(serverDir); err != nil {
		return fmt.Errorf("生成证书失败: %v", err)
	}
	
	// 生成Openvpn配置文件
	fmt.Println("正在生成Openvpn配置文件...")
	if err := generateOpenVPNConfig(serverDir); err != nil {
		return fmt.Errorf("生成OpenVPN配置文件失败: %v", err)
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