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

	// 生成新的服务端配置
	config, err := cfg.GenerateServerConfig()
	if err != nil {
		return fmt.Errorf("生成服务端配置失败: %v", err)
	}

	// 写入新的配置文件
	if err := os.WriteFile(constants.ServerConfigPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	// 重启 OpenVPN 服务
	if err := RestartServer(); err != nil {
		return fmt.Errorf("重启服务失败: %v", err)
	}

	return nil
}
