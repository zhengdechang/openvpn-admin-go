package openvpn

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/utils"
)

// EnsureMgmtPassword 确保 management 接口的密码文件存在。
// server.conf 的 `management 127.0.0.1 7505 <file>` 在 OpenVPN 启动(降权前)时读取该文件首行做口令，
// 缺失会导致服务起不来；这里幂等生成一次随机口令(0600)，重启复用同一口令，PauseClient/auth 脚本读同一文件。
func EnsureMgmtPassword() error {
	if info, err := os.Stat(constants.ServerMgmtPasswordPath); err == nil && info.Size() > 0 {
		return nil
	}
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Errorf("生成 management 口令失败: %v", err)
	}
	// 仅十六进制字符，避免 OpenVPN 管理协议里发送口令时的转义/截断问题
	pw := hex.EncodeToString(buf) + "\n"
	if err := os.WriteFile(constants.ServerMgmtPasswordPath, []byte(pw), 0600); err != nil {
		return fmt.Errorf("写入 management 口令文件失败: %v", err)
	}
	return nil
}

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

	// server.conf 会引用 management 口令文件，必须保证它先存在，否则 OpenVPN 起不来
	if err := EnsureMgmtPassword(); err != nil {
		return err
	}

	// server.conf 无条件引用 tls-verify.sh（按 CN 拉黑），开启 CRL 时还引用 crl.pem。
	// 这些文件必须在写配置/重启之前就位，否则 OpenVPN 拒绝启动 = 全员锁死。
	// EnsureServerHelperFiles 把脚本/配置刷进持久卷；EnsureCRLSetup 在有 CA 时生成初始空 CRL。
	if err := EnsureServerHelperFiles(); err != nil {
		return err
	}
	if err := EnsureCRLSetup(); err != nil {
		return err
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
	logDir := filepath.Dir(constants.DefaultOpenVPNStatusLogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建状态日志文件
	if err := os.WriteFile(constants.DefaultOpenVPNStatusLogPath, []byte{}, 0644); err != nil {
		return fmt.Errorf("创建状态日志文件失败: %v", err)
	}

	// 重启OpenVPN服务
	utils.SupervisorctlRestart(constants.SupervisorOpenVPNServiceName)

	return nil
}

// RestartServer 重启OpenVPN服务
func RestartServer() error {
	utils.SupervisorctlRestart(constants.SupervisorOpenVPNServiceName)
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
   // 重新写入 server.conf 并更新所有客户端、重启服务
   if err := UpdateServerConfig(); err != nil {
       return fmt.Errorf("更新服务器配置失败: %v", err)
   }
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
   return nil
}

// getEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}