package openvpn

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"openvpn-admin-go/constants"
)

// Config 存储所有配置
type Config struct {
	OpenVPNPort            int      `json:"openvpn_port"`
	OpenVPNProto           string   `json:"openvpn_proto"`
	OpenVPNSyncCerts       bool     `json:"openvpn_sync_certs"`
	OpenVPNUseCRL          bool     `json:"openvpn_use_crl"`
	OpenVPNServerHostname  string   `json:"openvpn_server_hostname"`
	OpenVPNServerNetwork   string   `json:"openvpn_server_network"`
	OpenVPNServerNetmask   string   `json:"openvpn_server_netmask"`
	OpenVPNRoutes          []string `json:"openvpn_routes"`
	OpenVPNClientConfigDir string   `json:"openvpn_client_config_dir"`
	OpenVPNTLSVersion      string   `json:"openvpn_tls_version"`
	OpenVPNTLSKey          string   `json:"openvpn_tls_key"`
	OpenVPNTLSKeyPath      string   `json:"openvpn_tls_key_path"`
	OpenVPNClientToClient  bool     `json:"openvpn_client_to_client"`
	DNSServerIP            string   `json:"dns_server_ip"`
	DNSServerDomain        string   `json:"dns_server_domain"`
	OpenVPNStatusLogPath   string   `json:"openvpn_status_log_path"`
	OpenVPNLogPath         string   `json:"openvpn_log_path"`
	OpenVPNManagementPort  int      `json:"openvpn_management_port,omitempty"`
	OpenVPNBlacklistFile   string   `json:"openvpn_blacklist_file,omitempty"`
}

// LoadConfig 从配置文件加载配置，优先使用 JSON 配置，回退到解析 server.conf
func LoadConfig() (*Config, error) {
	// 首先尝试从新的配置管理系统加载
	appCfg, err := loadAppConfig()
	if err == nil {
		// 转换为 OpenVPN Config 结构
		return convertAppConfigToConfig(appCfg), nil
	}

	// 如果新配置系统失败，回退到原有的解析方式
	return loadFromServerConfig()
}

// loadAppConfig 加载应用配置
func loadAppConfig() (*AppConfig, error) {
	// 这里我们需要手动实现配置加载，因为导入有问题
	if _, err := os.Stat(constants.ConfigJSONPath); err != nil {
		// JSON 配置文件不存在，创建默认配置
		return createDefaultAppConfig()
	}

	// 读取 JSON 配置文件
	data, err := os.ReadFile(constants.ConfigJSONPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &cfg, nil
}

// AppConfig 应用程序配置结构（临时定义，避免导入问题）
type AppConfig struct {
	OpenVPNPort            int      `json:"openvpn_port"`
	OpenVPNProto           string   `json:"openvpn_proto"`
	OpenVPNSyncCerts       bool     `json:"openvpn_sync_certs"`
	OpenVPNUseCRL          bool     `json:"openvpn_use_crl"`
	OpenVPNServerHostname  string   `json:"openvpn_server_hostname"`
	OpenVPNServerNetwork   string   `json:"openvpn_server_network"`
	OpenVPNServerNetmask   string   `json:"openvpn_server_netmask"`
	OpenVPNRoutes          []string `json:"openvpn_routes"`
	OpenVPNClientConfigDir string   `json:"openvpn_client_config_dir"`
	OpenVPNTLSVersion      string   `json:"openvpn_tls_version"`
	OpenVPNTLSKey          string   `json:"openvpn_tls_key"`
	OpenVPNTLSKeyPath      string   `json:"openvpn_tls_key_path"`
	OpenVPNClientToClient  bool     `json:"openvpn_client_to_client"`
	DNSServerIP            string   `json:"dns_server_ip"`
	DNSServerDomain        string   `json:"dns_server_domain"`
	OpenVPNStatusLogPath   string   `json:"openvpn_status_log_path"`
	OpenVPNLogPath         string   `json:"openvpn_log_path"`
	OpenVPNManagementPort  int      `json:"openvpn_management_port"`
	OpenVPNBlacklistFile   string   `json:"openvpn_blacklist_file"`
}

// createDefaultAppConfig 创建默认应用配置
func createDefaultAppConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		OpenVPNPort:            constants.DefaultOpenVPNPort,
		OpenVPNProto:           constants.DefaultOpenVPNProto,
		OpenVPNSyncCerts:       constants.DefaultOpenVPNSyncCerts,
		OpenVPNUseCRL:          constants.DefaultOpenVPNUseCRL,
		OpenVPNServerHostname:  constants.DefaultOPENVPN_SERVER_HOSTNAME,
		OpenVPNServerNetwork:   constants.DefaultOpenVPNServerNetwork,
		OpenVPNServerNetmask:   constants.DefaultOpenVPNServerNetmask,
		OpenVPNRoutes:          append([]string{}, constants.DefaultOpenVPNRoutes...),
		OpenVPNClientConfigDir: constants.DefaultOpenVPNClientConfigDir,
		OpenVPNTLSVersion:      constants.DefaultOpenVPNTLSVersion,
		OpenVPNTLSKey:          constants.DefaultOpenVPNTLSKey,
		OpenVPNTLSKeyPath:      constants.DefaultOpenVPNTLSKeyPath,
		OpenVPNClientToClient:  constants.DefaultOpenVPNClientToClient,
		DNSServerIP:            "",
		DNSServerDomain:        "",
		OpenVPNStatusLogPath:   constants.DefaultOpenVPNStatusLogPath,
		OpenVPNLogPath:         constants.DefaultOpenVPNLogPath,
		OpenVPNManagementPort:  constants.DefaultOpenVPNManagementPort,
		OpenVPNBlacklistFile:   constants.DefaultOpenVPNBlacklistFile,
	}

	// 保存默认配置
	if err := saveAppConfig(cfg); err != nil {
		return nil, fmt.Errorf("保存默认配置失败: %v", err)
	}

	return cfg, nil
}

// saveAppConfig 保存应用配置
func saveAppConfig(cfg *AppConfig) error {
	// 确保目录存在
	dir := filepath.Dir(constants.ConfigJSONPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 序列化配置
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(constants.ConfigJSONPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// convertAppConfigToConfig 将 AppConfig 转换为 Config
func convertAppConfigToConfig(appCfg *AppConfig) *Config {
	return &Config{
		OpenVPNPort:            appCfg.OpenVPNPort,
		OpenVPNProto:           appCfg.OpenVPNProto,
		OpenVPNSyncCerts:       appCfg.OpenVPNSyncCerts,
		OpenVPNUseCRL:          appCfg.OpenVPNUseCRL,
		OpenVPNServerHostname:  appCfg.OpenVPNServerHostname,
		OpenVPNServerNetwork:   appCfg.OpenVPNServerNetwork,
		OpenVPNServerNetmask:   appCfg.OpenVPNServerNetmask,
		OpenVPNRoutes:          append([]string{}, appCfg.OpenVPNRoutes...),
		OpenVPNClientConfigDir: appCfg.OpenVPNClientConfigDir,
		OpenVPNTLSVersion:      appCfg.OpenVPNTLSVersion,
		OpenVPNTLSKey:          appCfg.OpenVPNTLSKey,
		OpenVPNTLSKeyPath:      appCfg.OpenVPNTLSKeyPath,
		OpenVPNClientToClient:  appCfg.OpenVPNClientToClient,
		DNSServerIP:            appCfg.DNSServerIP,
		DNSServerDomain:        appCfg.DNSServerDomain,
		OpenVPNStatusLogPath:   appCfg.OpenVPNStatusLogPath,
		OpenVPNLogPath:         appCfg.OpenVPNLogPath,
		OpenVPNManagementPort:  appCfg.OpenVPNManagementPort,
		OpenVPNBlacklistFile:   appCfg.OpenVPNBlacklistFile,
	}
}

// loadFromServerConfig 从 server.conf 文件加载配置（回退方案）
func loadFromServerConfig() (*Config, error) {
	cfg := &Config{}

	// 检查配置文件是否存在，如果不存在则创建
	if _, err := os.Stat(constants.ServerConfigPath); os.IsNotExist(err) {
		// 使用默认配置创建
		appCfg, err := createDefaultAppConfig()
		if err != nil {
			return nil, fmt.Errorf("创建默认配置失败: %v", err)
		}

		cfg = convertAppConfigToConfig(appCfg)

		// 生成默认配置文件
		configContent, err := cfg.GenerateServerConfig()
		if err != nil {
			return nil, fmt.Errorf("生成默认配置文件失败: %v", err)
		}

		// 写入配置文件
		if err := os.WriteFile(constants.ServerConfigPath, []byte(configContent), 0644); err != nil {
			return nil, fmt.Errorf("写入配置文件失败: %v", err)
		}

		return cfg, nil
	}

	// 读取服务端配置文件
	configContent, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return nil, fmt.Errorf("读取服务端配置文件失败: %v", err)
	}

	// 解析配置文件
	lines := strings.Split(string(configContent), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "port":
			cfg.OpenVPNPort, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("解析端口失败: %v", err)
			}
		case "proto":
			cfg.OpenVPNProto = fields[1]
		case "server":
			if len(fields) >= 3 {
				cfg.OpenVPNServerNetwork = fields[1]
				cfg.OpenVPNServerNetmask = fields[2]
			}
		case "push":
			if strings.HasPrefix(fields[1], "route") {
				route := strings.Join(fields[2:], " ")
				cfg.OpenVPNRoutes = append(cfg.OpenVPNRoutes, route)
			}
		}
	}

	// 设置默认值（不再依赖环境变量）
	if cfg.OpenVPNPort == 0 {
		cfg.OpenVPNPort = constants.DefaultOpenVPNPort
	}
	if cfg.OpenVPNProto == "" {
		cfg.OpenVPNProto = constants.DefaultOpenVPNProto
	}
	if cfg.OpenVPNServerNetwork == "" {
		cfg.OpenVPNServerNetwork = constants.DefaultOpenVPNServerNetwork
	}
	if cfg.OpenVPNServerNetmask == "" {
		cfg.OpenVPNServerNetmask = constants.DefaultOpenVPNServerNetmask
	}

	// 使用常量默认值设置其他配置
	cfg.OpenVPNServerHostname = constants.DefaultOPENVPN_SERVER_HOSTNAME
	cfg.OpenVPNSyncCerts = constants.DefaultOpenVPNSyncCerts
	cfg.OpenVPNUseCRL = constants.DefaultOpenVPNUseCRL
	cfg.OpenVPNClientToClient = constants.DefaultOpenVPNClientToClient
	cfg.OpenVPNClientConfigDir = constants.DefaultOpenVPNClientConfigDir
	cfg.OpenVPNTLSVersion = constants.DefaultOpenVPNTLSVersion
	cfg.OpenVPNTLSKey = constants.DefaultOpenVPNTLSKey
	cfg.OpenVPNTLSKeyPath = constants.DefaultOpenVPNTLSKeyPath
	cfg.OpenVPNStatusLogPath = constants.DefaultOpenVPNStatusLogPath
	cfg.OpenVPNLogPath = constants.DefaultOpenVPNLogPath
	cfg.OpenVPNManagementPort = constants.DefaultOpenVPNManagementPort
	cfg.OpenVPNBlacklistFile = constants.DefaultOpenVPNBlacklistFile

	// 设置默认路由
	if len(cfg.OpenVPNRoutes) == 0 {
		cfg.OpenVPNRoutes = append([]string{}, constants.DefaultOpenVPNRoutes...)
	}

	// DNS 配置默认为空
	cfg.DNSServerIP = ""
	cfg.DNSServerDomain = ""

	return cfg, nil
}

// GenerateServerConfig 生成 OpenVPN 服务器配置
func (c *Config) GenerateServerConfig() (string, error) {
	config, err := RenderServerConfig(c)
	if err != nil {
		return "", fmt.Errorf("生成服务器配置失败: %v", err)
	}
	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数类型的环境变量
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// getEnvBool 获取布尔类型的环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

// getEnvList 获取字符串列表类型的环境变量
func getEnvList(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// SaveConfig 保存配置到 JSON 文件
func SaveConfig(cfg *Config) error {
	// 转换为 AppConfig
	appCfg := convertConfigToAppConfig(cfg)

	// 保存到 JSON 文件
	return saveAppConfig(appCfg)
}

// convertConfigToAppConfig 将 Config 转换为 AppConfig
func convertConfigToAppConfig(cfg *Config) *AppConfig {
	return &AppConfig{
		OpenVPNPort:            cfg.OpenVPNPort,
		OpenVPNProto:           cfg.OpenVPNProto,
		OpenVPNSyncCerts:       cfg.OpenVPNSyncCerts,
		OpenVPNUseCRL:          cfg.OpenVPNUseCRL,
		OpenVPNServerHostname:  cfg.OpenVPNServerHostname,
		OpenVPNServerNetwork:   cfg.OpenVPNServerNetwork,
		OpenVPNServerNetmask:   cfg.OpenVPNServerNetmask,
		OpenVPNRoutes:          append([]string{}, cfg.OpenVPNRoutes...),
		OpenVPNClientConfigDir: cfg.OpenVPNClientConfigDir,
		OpenVPNTLSVersion:      cfg.OpenVPNTLSVersion,
		OpenVPNTLSKey:          cfg.OpenVPNTLSKey,
		OpenVPNTLSKeyPath:      cfg.OpenVPNTLSKeyPath,
		OpenVPNClientToClient:  cfg.OpenVPNClientToClient,
		DNSServerIP:            cfg.DNSServerIP,
		DNSServerDomain:        cfg.DNSServerDomain,
		OpenVPNStatusLogPath:   cfg.OpenVPNStatusLogPath,
		OpenVPNLogPath:         cfg.OpenVPNLogPath,
		OpenVPNManagementPort:  cfg.OpenVPNManagementPort,
		OpenVPNBlacklistFile:   cfg.OpenVPNBlacklistFile,
	}
}
