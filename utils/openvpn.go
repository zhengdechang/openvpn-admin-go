package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindOpenVPNConfigDir 查找 OpenVPN 配置目录
func FindOpenVPNConfigDir() (string, error) {
	configDir := "/etc/openvpn"
	if _, err := os.Stat(configDir); err != nil {
		return "", fmt.Errorf("OpenVPN 配置目录不存在: %v", err)
	}

	// 检查必要的文件是否存在
	files := []string{"ca.crt", "server.crt", "server.key", "dh.pem"}
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(configDir, "server", file)); err != nil {
			return "", fmt.Errorf("证书文件 %s 不存在", file)
		}
	}

	return configDir, nil
} 