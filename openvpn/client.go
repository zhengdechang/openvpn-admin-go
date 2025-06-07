package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"openvpn-admin-go/constants"
)

// CreateClient 创建新的OpenVPN客户端
func CreateClient(username string) error {
	fmt.Printf("开始创建客户端: %s\n", username)

	// 检查证书目录
	fmt.Printf("检查证书目录: %s\n", constants.ClientConfigDir)
	if _, err := os.Stat(constants.ClientConfigDir); os.IsNotExist(err) {
		fmt.Printf("创建证书目录: %s\n", constants.ClientConfigDir)
		cmd := exec.Command("sudo", "mkdir", "-p", constants.ClientConfigDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("创建证书目录失败: %v, 输出: %s", err, string(output))
		}
		fmt.Println("证书目录创建成功")
	}

	// 检查并生成TLS密钥
	if _, err := os.Stat(constants.ServerTLSKeyPath); os.IsNotExist(err) {
		fmt.Printf("正在生成TLS密钥: %s\n", constants.ServerTLSKeyPath)
		cmd := exec.Command("sudo", "openvpn", "--genkey", "secret", constants.ServerTLSKeyPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("生成TLS密钥失败: %v, 输出: %s", err, string(output))
		}
		fmt.Println("TLS密钥生成成功")
	}

	// 检查CA证书和密钥是否存在
	fmt.Printf("检查CA证书: %s\n", constants.ServerCACertPath)
	fmt.Printf("检查CA密钥: %s\n", constants.ServerCAKeyPath)

	if _, err := os.Stat(constants.ServerCACertPath); os.IsNotExist(err) {
		return fmt.Errorf("CA证书不存在: %s", constants.ServerCACertPath)
	}
	if _, err := os.Stat(constants.ServerCAKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("CA密钥不存在: %s", constants.ServerCAKeyPath)
	}

	fmt.Printf("使用CA证书: %s\n", constants.ServerCACertPath)
	fmt.Printf("使用CA密钥: %s\n", constants.ServerCAKeyPath)
	fmt.Println("CA证书和密钥检查通过")

	// 检查客户端扩展文件
	clientExtFile := filepath.Join(filepath.Dir(constants.ServerConfigPath), "openssl-client.ext")
	if _, err := os.Stat(clientExtFile); os.IsNotExist(err) {
		return fmt.Errorf("客户端扩展文件不存在: %s", clientExtFile)
	}
	fmt.Printf("使用客户端扩展文件: %s\n", clientExtFile)

	// 生成客户端私钥
	fmt.Printf("正在为客户端 %s 生成私钥...\n", username)
	keyPath := filepath.Join(constants.ClientConfigDir, username+".key")
	cmd := exec.Command("sudo", "openssl", "genrsa", "-out", keyPath, "2048")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("生成私钥失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("私钥生成成功")

	// 生成证书签名请求
	fmt.Printf("正在为客户端 %s 生成证书签名请求...\n", username)
	csrPath := filepath.Join(constants.ClientConfigDir, username+".csr")
	cmd = exec.Command("sudo", "openssl", "req", "-new", "-key", keyPath, "-out", csrPath, "-subj", "/CN="+username)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("生成证书签名请求失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("证书签名请求生成成功")

	// 使用CA证书签名
	fmt.Printf("正在为客户端 %s 签名证书...\n", username)
	crtPath := filepath.Join(constants.ClientConfigDir, username+".crt")
	cmd = exec.Command("sudo", "openssl", "x509", "-req", "-in", csrPath, "-CA", constants.ServerCACertPath, "-CAkey", constants.ServerCAKeyPath, "-CAcreateserial", "-out", crtPath, "-days", "3650", "-extfile", clientExtFile)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("签名证书失败: %v, 输出: %s", err, string(output))
	}
	fmt.Println("证书签名成功")

	// 复制CA证书到客户端目录
	clientCaPath := filepath.Join(constants.ClientConfigDir, "ca.crt")
	fmt.Printf("正在复制CA证书到: %s\n", clientCaPath)
	cmd = exec.Command("sudo", "cp", constants.ServerCACertPath, clientCaPath)
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
	ovpnPath := filepath.Join(constants.ClientConfigDir, username+".ovpn")

	// 加载配置
	cfg, err := LoadConfig()
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

	fmt.Printf("客户端 %s 的证书和配置文件已生成并复制到 %s 目录\n", username, constants.ClientConfigDir)
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
	files := []string{
		username + ".key",
		username + ".crt",
		username + ".ovpn",
	}

	for _, file := range files {
		path := filepath.Join(constants.ClientConfigDir, file)
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
	ovpnPath := filepath.Join(constants.ClientConfigDir, username+".ovpn")

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
	if logContent, err := os.ReadFile(constants.ServerStatusLogPath); err == nil {
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
	var statuses []ClientStatus

	// 读取客户端目录
	files, err := os.ReadDir(constants.ClientConfigDir)
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
	Username     string    `json:"username"`
	ConnectedAt  time.Time `json:"connectedAt"`
	Disconnected time.Time `json:"disconnectedAt"`
	IsPaused     bool      `json:"isPaused"`
}

// GenerateClientConfig 生成客户端配置文件内容
func GenerateClientConfig(username string, cfg *Config) (string, error) {
	return RenderClientConfig(username, cfg)
}

// 创建客户端配置目录
func createClientConfigDir() error {
	if err := os.MkdirAll(constants.ClientConfigDir, 0755); err != nil {
		return fmt.Errorf("创建客户端配置目录失败: %v", err)
	}
	return nil
}

// 生成客户端配置文件
func generateClientConfigFile(clientName string) error {
	// 加载配置
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	config, err := GenerateClientConfig(clientName, cfg)
	if err != nil {
		return fmt.Errorf("生成客户端配置失败: %v", err)
	}

	// 写入配置文件
	configPath := filepath.Join(constants.ClientConfigDir, clientName+".ovpn")
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	return nil
}

// 复制证书文件
func copyFile(source, destination string) error {
	sourceContent, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("读取源文件失败: %v", err)
	}

	if err := os.WriteFile(destination, sourceContent, 0644); err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}
	return nil
}
