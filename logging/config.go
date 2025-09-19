package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LogConfig 日志配置结构
type LogConfig struct {
	// 基本配置
	Level         string `json:"level"`          // 日志级别: debug, info, warn, error, fatal
	LogFilePath   string `json:"log_file_path"`  // 日志文件路径
	EnableAPI     bool   `json:"enable_api"`     // 是否启用API日志记录
	EnableConsole bool   `json:"enable_console"` // 是否输出到控制台

	// 文件配置
	MaxFileSize int64 `json:"max_file_size"` // 最大文件大小（字节）
	MaxBackups  int   `json:"max_backups"`   // 最大备份文件数
	MaxAge      int   `json:"max_age"`       // 最大保存天数

	// API日志配置
	API APILogConfig `json:"api"`
}

// APILogConfig API日志配置
type APILogConfig struct {
	EnableRequestBody  bool     `json:"enable_request_body"`  // 是否记录请求体
	EnableResponseBody bool     `json:"enable_response_body"` // 是否记录响应体
	MaxBodySize        int64    `json:"max_body_size"`        // 最大记录的请求/响应体大小
	SkipPaths          []string `json:"skip_paths"`           // 跳过记录的路径
	SensitivePaths     []string `json:"sensitive_paths"`      // 敏感路径（需要特殊处理）
}

// DefaultLogConfig 返回默认的日志配置
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Level:         "info",
		LogFilePath:   "logs/web.log",
		EnableAPI:     true,
		EnableConsole: false,
		MaxFileSize:   100 * 1024 * 1024, // 100MB
		MaxBackups:    5,
		MaxAge:        30, // 30天
		API: APILogConfig{
			EnableRequestBody:  false,
			EnableResponseBody: false,
			MaxBodySize:        1024, // 1KB
			SkipPaths: []string{
				"/api/health",
				"/api/ping",
				"/favicon.ico",
			},
			SensitivePaths: []string{
				"/api/auth/login",
				"/api/auth/logout",
				"/api/users",
				"/api/client",
				"/api/server",
				"/api/departments",
			},
		},
	}
}

// LoadLogConfig 从文件加载日志配置
func LoadLogConfig(configPath string) (LogConfig, error) {
	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := DefaultLogConfig()
		if err := SaveLogConfig(configPath, defaultConfig); err != nil {
			return defaultConfig, fmt.Errorf("创建默认配置文件失败: %v", err)
		}
		return defaultConfig, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return LogConfig{}, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config LogConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return LogConfig{}, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := validateLogConfig(config); err != nil {
		return LogConfig{}, fmt.Errorf("配置验证失败: %v", err)
	}

	return config, nil
}

// SaveLogConfig 保存日志配置到文件
func SaveLogConfig(configPath string, config LogConfig) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 序列化配置
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// validateLogConfig 验证日志配置
func validateLogConfig(config LogConfig) error {
	// 验证日志级别
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLevels[config.Level] {
		return fmt.Errorf("无效的日志级别: %s", config.Level)
	}

	// 验证文件路径
	if config.LogFilePath == "" {
		return fmt.Errorf("日志文件路径不能为空")
	}

	// 验证文件大小
	if config.MaxFileSize <= 0 {
		return fmt.Errorf("最大文件大小必须大于0")
	}

	// 验证备份数量
	if config.MaxBackups < 0 {
		return fmt.Errorf("最大备份数量不能为负数")
	}

	// 验证保存天数
	if config.MaxAge < 0 {
		return fmt.Errorf("最大保存天数不能为负数")
	}

	return nil
}

// ParseLogLevel 解析日志级别字符串
func ParseLogLevel(levelStr string) LogLevel {
	switch levelStr {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// InitFromConfig 从配置初始化日志系统
func InitFromConfig(configPath string) error {
	// 加载配置
	config, err := LoadLogConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载日志配置失败: %v", err)
	}

	// 解析日志级别
	level := ParseLogLevel(config.Level)

	// 初始化日志系统
	logConfig := Config{
		LogLevel:      level,
		LogFilePath:   config.LogFilePath,
		EnableAPI:     config.EnableAPI,
		MaxFileSize:   config.MaxFileSize,
		EnableConsole: config.EnableConsole,
	}

	if err := Init(logConfig); err != nil {
		return fmt.Errorf("初始化日志系统失败: %v", err)
	}

	Info("日志系统初始化成功 - Level: %s, File: %s, API: %v",
		config.Level, config.LogFilePath, config.EnableAPI)

	return nil
}

// GetAPILogConfig 获取API日志配置
func GetAPILogConfig(configPath string) (APILogConfig, error) {
	config, err := LoadLogConfig(configPath)
	if err != nil {
		return APILogConfig{}, err
	}
	return config.API, nil
}

// UpdateLogLevel 动态更新日志级别
func UpdateLogLevel(configPath, newLevel string) error {
	// 加载当前配置
	config, err := LoadLogConfig(configPath)
	if err != nil {
		return err
	}

	// 更新日志级别
	config.Level = newLevel

	// 验证配置
	if err := validateLogConfig(config); err != nil {
		return err
	}

	// 保存配置
	if err := SaveLogConfig(configPath, config); err != nil {
		return err
	}

	// 更新运行时日志级别
	SetLogLevel(ParseLogLevel(newLevel))

	Info("日志级别已更新为: %s", newLevel)
	return nil
}

// GetCurrentConfig 获取当前运行时配置信息
func GetCurrentConfig() map[string]interface{} {
	return map[string]interface{}{
		"level":       GetLogLevel().String(),
		"api_enabled": IsAPILoggingEnabled(),
	}
}
