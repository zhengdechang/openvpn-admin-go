package openvpn

import (
	"fmt"
	"os"

	"openvpn-admin-go/constants"
)

// UpdatePort 更新端口号
func UpdatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("端口号必须在 1-65535 之间")
	}
	// 加载配置
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	// 更新配置
	cfg.OpenVPNPort = port
	// 保存配置
	if err := SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	// 重启 OpenVPN 服务
	// SaveConfig now handles writing to constants.ServerConfigPath by calling GenerateServerConfig
	if err := RestartServer(); err != nil {
		return fmt.Errorf("重启服务失败: %v", err)
	}

	return nil
}
