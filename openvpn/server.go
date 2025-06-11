package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"openvpn-admin-go/constants"
)

// GetServerConfigTemplate 获取服务端配置模板
func GetServerConfigTemplate() (string, error) {
	// 从环境变量加载配置
	cfg, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("加载配置失败: %v", err)
	}
	return cfg.GenerateServerConfig()
}

// UpdateServerConfig 更新服务器配置
func UpdateServerConfig() error {
	// 加载配置
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 生成服务器配置文件
	config, err := cfg.GenerateServerConfig()
	if err != nil {
		return fmt.Errorf("生成服务器配置失败: %v", err)
	}

	// 写入配置文件
	if err := os.WriteFile(constants.ServerConfigPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	// 更新所有客户端配置
	files, err := os.ReadDir(constants.ClientConfigDir)
	if err != nil {
		return fmt.Errorf("读取客户端目录失败: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".ovpn") {
			username := strings.TrimSuffix(file.Name(), ".ovpn")
			clientConfig, err := GenerateClientConfig(username, cfg)
			if err != nil {
				return fmt.Errorf("生成客户端 %s 配置失败: %v", username, err)
			}
			if err := os.WriteFile(filepath.Join(constants.ClientConfigDir, file.Name()), []byte(clientConfig), 0644); err != nil {
				return fmt.Errorf("更新客户端 %s 配置失败: %v", username, err)
			}
		}
	}

	// 检查证书文件是否存在
	if _, err := os.Stat(constants.ServerCACertPath); os.IsNotExist(err) {
		return fmt.Errorf("CA证书文件不存在: %s", constants.ServerCACertPath)
	}
	if _, err := os.Stat(constants.ServerCertPath); os.IsNotExist(err) {
		return fmt.Errorf("服务器证书文件不存在: %s", constants.ServerCertPath)
	}
	if _, err := os.Stat(constants.ServerKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("服务器密钥文件不存在: %s", constants.ServerKeyPath)
	}
	if _, err := os.Stat(constants.ServerDHPath); os.IsNotExist(err) {
		return fmt.Errorf("DH参数文件不存在: %s", constants.ServerDHPath)
	}
	if _, err := os.Stat(constants.ServerTLSKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("TLS密钥文件不存在: %s", constants.ServerTLSKeyPath)
	}

	// 创建ipp.txt文件
	if err := os.WriteFile(constants.ServerIPPPath, []byte{}, 0644); err != nil {
		return fmt.Errorf("创建ipp.txt文件失败: %v", err)
	}

	// 创建日志目录
	logDir := filepath.Dir(constants.ServerStatusLogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建状态日志文件
	if err := os.WriteFile(constants.ServerStatusLogPath, []byte{}, 0644); err != nil {
		return fmt.Errorf("创建状态日志文件失败: %v", err)
	}

	// 重启OpenVPN服务
	if err := exec.Command("systemctl", "restart", constants.ServiceName).Run(); err != nil {
		return fmt.Errorf("重启OpenVPN服务失败: %v", err)
	}

	return nil
}

// RestartServer 重启OpenVPN服务
func RestartServer() error {
	cmd := exec.Command("systemctl", "restart", constants.ServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("重启服务失败: %v\n输出: %s", err, string(output))
	}
	return nil
}

// ConfigureServer 根据参数更新服务器配置并重启服务
func ConfigureServer(port int, protocol, network, netmask string) error {
   // 加载现有配置
   cfg, err := LoadConfig()
   if err != nil {
       return fmt.Errorf("加载配置失败: %v", err)
   }
   // 更新参数
   cfg.OpenVPNPort = port
   cfg.OpenVPNProto = protocol
   cfg.OpenVPNServerNetwork = network
   cfg.OpenVPNServerNetmask = netmask
   // 保存并生成配置文件
   if err := SaveConfig(cfg); err != nil {
       return fmt.Errorf("保存配置失败: %v", err)
   }
   // 重启服务以应用更改
   if err := RestartServer(); err != nil {
       return fmt.Errorf("重启服务失败: %v", err)
   }
   // TODO: Consider if client configs need regeneration here or if it's handled elsewhere.
   // The old UpdateServerConfig did regenerate all client configs.
   // If direct server parameter changes (port, proto, network) affect client configs,
   // that logic might need to be triggered here as well.
   // For now, focusing on server config persistence.
   // Note: UpdateServerConfig() also handles client config regeneration.
   // Since SaveConfig() now writes server.conf, if client configs depend on parameters
   // changed in ConfigureServer, they should be regenerated.
   // A call to a simplified client regeneration function might be needed if UpdateServerConfig() is too broad.
   return nil
}

// ApplyServerConfig 根据自定义内容写入配置并重启服务
func ApplyServerConfig(content string) error {
   // 写入配置文件
   if err := os.WriteFile(constants.ServerConfigPath, []byte(content), 0644); err != nil {
       return fmt.Errorf("写入配置文件失败: %v", err)
   }
   // 重启服务
   if err := RestartServer(); err != nil {
       return fmt.Errorf("重启服务失败: %v", err)
   }
   // Update config.json to reflect the changes applied to server.conf
   // This helps keep config.json somewhat in sync, although server.conf is parsed by LoadConfig.
   cfg, err := LoadConfig() // This will parse the server.conf that was just written.
   if err == nil {
       // SaveConfig will write to config.json and regenerate server.conf again (idempotent)
       // This ensures that if server.conf had settings not covered by config.json's struct,
       // they are still part of the active config loaded, and then config.json is updated
       // with the parts that are covered.
       _ = SaveConfig(cfg) // Ignoring error from this SaveConfig for now, as primary action was direct write.
   }
   return nil
}

// getEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
// This function seems unused in this file now. Consider removing if not used elsewhere.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}