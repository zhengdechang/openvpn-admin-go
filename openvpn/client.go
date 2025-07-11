package openvpn

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/utils"
)

// CreateClient 创建新的OpenVPN客户端
func CreateClient(username string) error {
	fmt.Printf("开始创建客户端: %s\n", username)

	// 检查证书目录
	fmt.Printf("检查证书目录: %s\n", constants.ClientConfigDir)
	if _, err := os.Stat(constants.ClientConfigDir); os.IsNotExist(err) {
		fmt.Printf("创建证书目录: %s\n", constants.ClientConfigDir)
		if err := utils.ExecCommand(fmt.Sprintf("mkdir -p %s", constants.ClientConfigDir)); err != nil {
			return fmt.Errorf("创建证书目录失败: %v", err)
		}
		fmt.Println("证书目录创建成功")
	}

	// 检查并生成TLS密钥
	if _, err := os.Stat(constants.ServerTLSKeyPath); os.IsNotExist(err) {
		fmt.Printf("正在生成TLS密钥: %s\n", constants.ServerTLSKeyPath)
		if err := utils.ExecCommand(fmt.Sprintf("openvpn --genkey secret %s", constants.ServerTLSKeyPath)); err != nil {
			return fmt.Errorf("生成TLS密钥失败: %v", err)
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
	if err := utils.ExecCommand(fmt.Sprintf("openssl genrsa -out %s 2048", keyPath)); err != nil {
		return fmt.Errorf("生成私钥失败: %v", err)
	}
	fmt.Println("私钥生成成功")

	// 生成证书签名请求
	fmt.Printf("正在为客户端 %s 生成证书签名请求...\n", username)
	csrPath := filepath.Join(constants.ClientConfigDir, username+".csr")
	if err := utils.ExecCommand(fmt.Sprintf("openssl req -new -key %s -out %s -subj '/CN=%s'", keyPath, csrPath, username)); err != nil {
		return fmt.Errorf("生成证书签名请求失败: %v", err)
	}
	fmt.Println("证书签名请求生成成功")

	// 使用CA证书签名
	fmt.Printf("正在为客户端 %s 签名证书...\n", username)
	crtPath := filepath.Join(constants.ClientConfigDir, username+".crt")
	if err := utils.ExecCommand(fmt.Sprintf("openssl x509 -req -in %s -CA %s -CAkey %s -CAcreateserial -out %s -days 3650 -extfile %s", csrPath, constants.ServerCACertPath, constants.ServerCAKeyPath, crtPath, clientExtFile)); err != nil {
		return fmt.Errorf("签名证书失败: %v", err)
	}
	fmt.Println("证书签名成功")

	// 复制CA证书到客户端目录
	clientCaPath := filepath.Join(constants.ClientConfigDir, "ca.crt")
	fmt.Printf("正在复制CA证书到: %s\n", clientCaPath)
	if err := utils.ExecCommand(fmt.Sprintf("cp %s %s", constants.ServerCACertPath, clientCaPath)); err != nil {
		return fmt.Errorf("复制CA证书失败: %v", err)
	}
	fmt.Println("CA证书复制成功")

	// 清理临时文件
	fmt.Printf("清理临时文件: %s\n", csrPath)
	utils.ExecCommand(fmt.Sprintf("rm %s", csrPath)) // 忽略错误，因为文件可能不存在

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
	if err := utils.ExecCommand(fmt.Sprintf("mv %s %s", tempFile, ovpnPath)); err != nil {
		return fmt.Errorf("移动配置文件失败: %v", err)
	}

	// 设置文件权限
	fmt.Printf("设置文件权限: %s\n", ovpnPath)
	if err := utils.ExecCommand(fmt.Sprintf("chmod 644 %s", ovpnPath)); err != nil {
		return fmt.Errorf("设置文件权限失败: %v", err)
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
	// Connect to OpenVPN management interface
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", constants.DefaultOpenVPNManagementPort))
	if err != nil {
		return fmt.Errorf("failed to connect to OpenVPN management interface: %w", err)
	}
	defer conn.Close()

	// Send kill command
	_, err = fmt.Fprintf(conn, "kill %s\n", username)
	if err != nil {
		// It's possible the client is not connected, so log but don't necessarily fail hard
		fmt.Printf("Failed to send kill command for %s: %v. This might be okay if client was not connected.\n", username, err)
	} else {
		// Optional: Read response
		status, readErr := bufio.NewReader(conn).ReadString('\n')
		if readErr != nil {
			fmt.Printf("Failed to read response from management interface after killing %s: %v\n", username, readErr)
		} else {
			fmt.Printf("OpenVPN management interface response for kill %s: %s\n", username, status)
		}
	}

	// Append username to blacklist file
	// Using constants.DefaultOpenVPNBlacklistFile as per project structure
	// Note: The constants package has `DefaultOpenVPNBlacklistFile` (string)
	// and `BlacklistFile` (array of strings). We need the specific file path.
	// The original code in PauseClient correctly used `constants.BlacklistFile`
	// which was a string variable in that context.
	// Let's assume constants.DefaultOpenVPNBlacklistFile is the correct one to use.
	// If constants.BlacklistFile was intended to be a single path, its type in constants.go is confusing.
	// For now, sticking to DefaultOpenVPNBlacklistFile for clarity.
	blacklistFilePath := constants.DefaultOpenVPNBlacklistFile // Use the specific constant for the file path
	f, err := os.OpenFile(blacklistFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open blacklist file %s: %w", blacklistFilePath, err)
	}
	defer f.Close()

	// Check if user is already in blacklist to avoid duplicates
	existingContent, err := os.ReadFile(blacklistFilePath)
	if err != nil {
		// If cannot read, proceed to write, but log it
		fmt.Printf("Could not read blacklist file %s before appending: %v\n", blacklistFilePath, err)
	} else {
		// Ensure matching the whole line for the username
		lines := strings.Split(string(existingContent), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == username {
				fmt.Printf("User %s already in blacklist %s\n", username, blacklistFilePath)
				return nil // Already blacklisted
			}
		}
	}

	if _, err := fmt.Fprintln(f, username); err != nil {
		return fmt.Errorf("failed to write to blacklist file %s: %w", blacklistFilePath, err)
	}
	fmt.Printf("User %s added to blacklist %s\n", username, blacklistFilePath)
	return nil
}

// ResumeClient 恢复OpenVPN客户端
func ResumeClient(username string) error {
	blacklistFilePath := constants.DefaultOpenVPNBlacklistFile // Use the specific constant for the file path

	content, err := os.ReadFile(blacklistFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If blacklist file doesn't exist, user is not paused.
			fmt.Printf("Blacklist file %s does not exist, user %s is not paused.\n", blacklistFilePath, username)
			return nil
		}
		return fmt.Errorf("failed to read blacklist file %s: %w", blacklistFilePath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	found := false
	for _, line := range lines {
		// Trim space to ensure exact match, and check if line is the username
		if strings.TrimSpace(line) == username {
			found = true
			// Skip this line to remove the user from blacklist
			continue
		}
		// Keep non-empty lines, effectively filtering out the target username and empty lines
		if strings.TrimSpace(line) != "" {
			newLines = append(newLines, line)
		}
	}

	if !found {
		fmt.Printf("User %s not found in blacklist file %s.\n", username, blacklistFilePath)
		return nil // User not found, considered success for idempotency
	}

	// Join the remaining lines. If newLines is empty, newContent will be empty.
	// If newLines has items, join them with \n. Add a trailing \n if there's content.
	newContent := strings.Join(newLines, "\n")
	if len(newContent) > 0 && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	} else if len(newLines) == 0 { // If all lines were removed or file was empty to begin with
		newContent = "" // Ensure file is empty
	}


	err = os.WriteFile(blacklistFilePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated blacklist file %s: %w", blacklistFilePath, err)
	}

	fmt.Printf("User %s removed from blacklist %s\n", username, blacklistFilePath)
	return nil
}

// GetClientStatus 获取OpenVPN客户端状态
func GetClientStatus(username string) (*ClientStatus, error) {
	status, err := ParseClientStatus(username)
	if err != nil {
		return nil, err
	}
	if status == nil {
		return nil, nil
	}
	return &ClientStatus{
		CommonName:     status.CommonName,
		RealAddress:    status.RealAddress,
		VirtualAddress: status.VirtualAddress,
		BytesReceived:  status.BytesReceived,
		BytesSent:      status.BytesSent,
		ConnectedSince: status.ConnectedSince,
		LastRef:        status.LastRef,
	}, nil
}

// GetAllClientStatuses 获取所有OpenVPN客户端状态
func GetAllClientStatuses() ([]ClientStatus, error) {
	statuses, err := ParseAllClientStatuses()
	if err != nil {
		return nil, err
	}

	result := make([]ClientStatus, len(statuses))
	for i, status := range statuses {
		result[i] = ClientStatus{
			CommonName:     status.CommonName,
			RealAddress:    status.RealAddress,
			VirtualAddress: status.VirtualAddress,
			BytesReceived:  status.BytesReceived,
			BytesSent:      status.BytesSent,
			ConnectedSince: status.ConnectedSince,
			LastRef:        status.LastRef,
		}
	}
	return result, nil
}

// ClientStatus 客户端状态
type ClientStatus struct {
	CommonName     string    `json:"commonName"`
	RealAddress    string    `json:"realAddress"`
	VirtualAddress string    `json:"virtualAddress"`
	BytesReceived  int64     `json:"bytesReceived"`
	BytesSent      int64     `json:"bytesSent"`
	ConnectedSince time.Time `json:"connectedSince"`
	LastRef        time.Time `json:"lastRef"`
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