package openvpn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings" // Required for netmask validation if added, or general string ops
	"strconv"
	"encoding/binary"
	// "log" // For logging errors if necessary
)

// SetClientFixedIP creates or updates a client-specific configuration file (CCD)
// to assign a fixed IP address to a client.
// commonName is typically the user's ID.
// ipAddress is the fixed IP to assign.
// This function will fetch the serverNetmask from the main OpenVPN configuration.
func SetClientFixedIP(commonName string, ipAddress string) error {
	if commonName == "" {
		return fmt.Errorf("commonName cannot be empty")
	}
	if ipAddress == "" {
		return fmt.Errorf("ipAddress cannot be empty")
	}

	// TODO: Add IP address format validation for ipAddress

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		return fmt.Errorf("OpenVPNClientConfigDir is not set in the configuration")
	}
	if cfg.OpenVPNServerNetmask == "" {
		return fmt.Errorf("OpenVPNServerNetmask is not set in the configuration")
	}

	// Ensure the CCD directory exists
	ccdDir := filepath.Join(cfg.OpenVPNClientConfigDir, "ccd")
	if err := os.MkdirAll(ccdDir, 0755); err != nil {
		return fmt.Errorf("failed to create ccd directory '%s': %w", ccdDir, err)
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir,"ccd", commonName)
	newConfig := fmt.Sprintf("ifconfig-push %s %s", ipAddress, cfg.OpenVPNServerNetmask)

	// 读取现有文件内容
	var existingContent string
	if _, err := os.Stat(ccdFilePath); err == nil {
		content, err := os.ReadFile(ccdFilePath)
		if err != nil {
			return fmt.Errorf("failed to read existing config file: %w", err)
		}
		existingContent = string(content)
	}

	// 处理现有配置
	lines := strings.Split(existingContent, "\n")
	newLines := make([]string, 0)
	configAdded := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "ifconfig-push") {
			// 替换已存在的ifconfig-push配置
			newLines = append(newLines, newConfig)
			configAdded = true
		} else if strings.HasPrefix(trimmedLine, "iroute") {
			// 保留iroute配置
			newLines = append(newLines, line)
		} else if trimmedLine != "" {
			// 保留其他配置
			newLines = append(newLines, line)
		}
	}

	// 如果没有找到ifconfig-push配置，添加新的配置
	if !configAdded {
		newLines = append(newLines, newConfig)
	}

	// 写入更新后的内容
	content := strings.Join(newLines, "\n")
	if content != "" {
		content += "\n"
	}

	err = os.WriteFile(ccdFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	// log.Printf("Successfully set fixed IP for %s to %s in %s", commonName, ipAddress, ccdFilePath)
	return nil
}

// RemoveClientFixedIP removes the client-specific configuration file (CCD)
// for the given commonName, effectively removing their fixed IP assignment.
func RemoveClientFixedIP(commonName string) error {
	if commonName == "" {
		return fmt.Errorf("commonName cannot be empty")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		// If the directory isn't set, there's nothing to remove.
		// Depending on desired strictness, this could be an error or a silent success.
		// log.Printf("OpenVPNClientConfigDir is not set, skipping removal for %s", commonName)
		return nil
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir, "ccd", commonName)

	// Check if the file exists before trying to remove it
	if _, err := os.Stat(ccdFilePath); os.IsNotExist(err) {
		// File doesn't exist, so consider it successfully "removed"
		// log.Printf("CCD file for %s does not exist at %s, no action needed.", commonName, ccdFilePath)
		return nil
	} else if err != nil {
		// Other error during stat
		return fmt.Errorf("failed to check client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	err = os.Remove(ccdFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	// log.Printf("Successfully removed fixed IP config for %s from %s", commonName, ccdFilePath)
	return nil
}

// GetClientFixedIP reads the fixed IP address from a client-specific configuration file.
// It returns the IP address string if found, or an empty string if not found or on error.
func GetClientFixedIP(commonName string) (string, error) {
	if commonName == "" {
		return "", fmt.Errorf("commonName cannot be empty")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		// log.Printf("OpenVPNClientConfigDir is not set, cannot get fixed IP for %s", commonName)
		return "", nil // Or return an error: fmt.Errorf("OpenVPNClientConfigDir is not set")
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir, "ccd", commonName)

	if _, err := os.Stat(ccdFilePath); os.IsNotExist(err) {
		// File doesn't exist, no fixed IP assigned
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to check client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	contentBytes, err := os.ReadFile(ccdFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	content := string(contentBytes)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "ifconfig-push") {
			parts := strings.Fields(trimmedLine)
			// Expected: "ifconfig-push" "ip_address" "netmask"
			if len(parts) == 3 {
				// TODO: Validate that parts[1] is a valid IP address
				return parts[1], nil
			}
			// log.Printf("Malformed ifconfig-push line in %s for %s: %s", ccdFilePath, commonName, line)
		}
	}

	// log.Printf("No ifconfig-push directive found in %s for %s", ccdFilePath, commonName)
	return "", nil // No ifconfig-push line found
}

// cidrToNetmask 将CIDR格式转换为子网掩码格式
// 例如：将"10.10.120.0/23"转换为"10.10.120.0 255.255.254.0"
func cidrToNetmask(cidr string) (string, error) {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid CIDR format: %s", cidr)
	}

	ip := parts[0]
	maskBits, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid mask bits: %s", parts[1])
	}

	if maskBits < 0 || maskBits > 32 {
		return "", fmt.Errorf("mask bits must be between 0 and 32")
	}

	// 计算子网掩码
	mask := uint32(0xffffffff) << (32 - maskBits)
	maskBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(maskBytes, mask)

	// 格式化子网掩码
	netmask := fmt.Sprintf("%d.%d.%d.%d", maskBytes[0], maskBytes[1], maskBytes[2], maskBytes[3])
	return fmt.Sprintf("%s %s", ip, netmask), nil
}

// SetClientSubnet 为指定的客户端设置子网配置
// commonName 是客户端的标识名
// subnet 是要设置的子网地址（例如：10.10.120.0/23）
func SetClientSubnet(commonName string, subnet string) error {
	if commonName == "" {
		return fmt.Errorf("commonName cannot be empty")
	}
	if subnet == "" {
		return fmt.Errorf("subnet cannot be empty")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		return fmt.Errorf("OpenVPNClientConfigDir is not set in the configuration")
	}

	// 确保CCD目录存在
	ccdDir := filepath.Join(cfg.OpenVPNClientConfigDir, "ccd")
	if err := os.MkdirAll(ccdDir, 0755); err != nil {
		return fmt.Errorf("failed to create ccd directory '%s': %w", ccdDir, err)
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir, "ccd", commonName)

	// 转换CIDR格式为子网掩码格式
	subnetWithMask, err := cidrToNetmask(subnet)
	if err != nil {
		return fmt.Errorf("failed to convert subnet format: %w", err)
	}

	newConfig := fmt.Sprintf("iroute %s", subnetWithMask)

	// 读取现有文件内容
	var existingContent string
	if _, err := os.Stat(ccdFilePath); err == nil {
		content, err := os.ReadFile(ccdFilePath)
		if err != nil {
			return fmt.Errorf("failed to read existing config file: %w", err)
		}
		existingContent = string(content)
	}

	// 处理现有配置
	lines := strings.Split(existingContent, "\n")
	newLines := make([]string, 0)
	configAdded := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "iroute") {
			// 替换已存在的iroute配置
			newLines = append(newLines, newConfig)
			configAdded = true
		} else if strings.HasPrefix(trimmedLine, "ifconfig-push") {
			// 保留ifconfig-push配置
			newLines = append(newLines, line)
		} else if trimmedLine != "" {
			// 保留其他配置
			newLines = append(newLines, line)
		}
	}

	// 如果没有找到iroute配置，添加新的配置
	if !configAdded {
		newLines = append(newLines, newConfig)
	}

	// 写入更新后的内容
	content := strings.Join(newLines, "\n")
	if content != "" {
		content += "\n"
	}

	err = os.WriteFile(ccdFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write client subnet config file '%s': %w", ccdFilePath, err)
	}

	return nil
}

// RemoveClientSubnet 移除客户端的子网配置
// commonName 是客户端的标识名
func RemoveClientSubnet(commonName string) error {
	if commonName == "" {
		return fmt.Errorf("commonName cannot be empty")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		// 如果目录未设置，则无需移除
		return nil
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir, "ccd", commonName)

	// 检查文件是否存在
	if _, err := os.Stat(ccdFilePath); os.IsNotExist(err) {
		// 文件不存在，视为已成功"移除"
		return nil
	} else if err != nil {
		// 其他错误
		return fmt.Errorf("failed to check client subnet config file '%s': %w", ccdFilePath, err)
	}

	// 读取现有文件内容
	content, err := os.ReadFile(ccdFilePath)
	if err != nil {
		return fmt.Errorf("failed to read client config file '%s': %w", ccdFilePath, err)
	}

	// 处理文件内容，移除iroute配置
	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmedLine, "iroute") && trimmedLine != "" {
			newLines = append(newLines, line)
		}
	}

	// 如果文件为空，则删除文件
	if len(newLines) == 0 {
		err = os.Remove(ccdFilePath)
		if err != nil {
			return fmt.Errorf("failed to remove empty client config file '%s': %w", ccdFilePath, err)
		}
		return nil
	}

	// 写入更新后的内容
	newContent := strings.Join(newLines, "\n")
	if newContent != "" {
		newContent += "\n"
	}

	err = os.WriteFile(ccdFilePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update client config file '%s': %w", ccdFilePath, err)
	}

	return nil
}
