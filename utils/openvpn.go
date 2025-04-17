package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindOpenVPNConfigDir 查找 OpenVPN 配置目录
func FindOpenVPNConfigDir() (string, error) {
	possiblePaths := []string{
		"/etc/openvpn",
		"/etc/openvpn/server",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			// 检查必要的文件是否存在
			files := []string{"ca.crt", "server.crt", "server.key", "dh.pem"}
			allExist := true
			for _, file := range files {
				if _, err := os.Stat(filepath.Join(path, file)); err != nil {
					allExist = false
					break
				}
			}
			if allExist {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的 OpenVPN 配置目录")
} 