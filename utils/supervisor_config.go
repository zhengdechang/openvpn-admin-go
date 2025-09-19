package utils

import (
	"bytes"
	"fmt"
	"openvpn-admin-go/constants"
	"os"
	"path/filepath"
	"text/template"
)

// SupervisorConfig 包含 supervisor 配置的数据结构
type SupervisorConfig struct {
	BinaryPath       string
	WorkingDirectory string
	Port             int
	DBPath           string
	OpenVPNConfigDir string
	OpenVPNAutoStart bool
	WebAutoStart     bool
}

// ServiceConfig 单个服务的配置结构
type ServiceConfig struct {
	BinaryPath       string
	WorkingDirectory string
	Port             int
	DBPath           string
	OpenVPNConfigDir string
	AutoStart        bool
}

// InstallSupervisorMainConfig 安装 supervisor 主配置文件
func InstallSupervisorMainConfig() error {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 解析模板
	templatePath := filepath.Join(wd, "template", "supervisord.conf.j2")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("解析 supervisor 主配置模板失败: %v", err)
	}

	// 生成配置文件内容（主配置不需要动态数据）
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return fmt.Errorf("渲染 supervisor 主配置模板失败: %v", err)
	}

	// 确保配置目录存在
	if err := os.MkdirAll(constants.SupervisorConfDir, 0755); err != nil {
		return fmt.Errorf("创建 supervisor 配置目录失败: %v", err)
	}

	// 写入 supervisor 主配置文件
	if err := os.WriteFile(constants.SupervisorConfigPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入 supervisor 主配置文件失败: %v", err)
	}

	fmt.Printf("Supervisor 主配置已安装: %s\n", constants.SupervisorConfigPath)
	return nil
}

// InstallOpenVPNServiceConfig 安装 OpenVPN 服务配置
func InstallOpenVPNServiceConfig(autoStart bool) error {
	return installServiceConfig("openvpn-server.conf.j2", constants.SupervisorOpenVPNConfigPath, ServiceConfig{
		AutoStart: autoStart,
	})
}

// InstallWebServiceConfig 安装 Web 服务配置
func InstallWebServiceConfig(config ServiceConfig) error {
	// 设置默认值
	if config.BinaryPath == "" {
		binaryPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("获取可执行文件路径失败: %v", err)
		}
		config.BinaryPath = binaryPath
	}
	if config.WorkingDirectory == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("获取工作目录失败: %v", err)
		}
		config.WorkingDirectory = wd
	}
	if config.Port == 0 {
		config.Port = 8085
	}
	if config.DBPath == "" {
		config.DBPath = "/app/data/db.sqlite3"
	}
	if config.OpenVPNConfigDir == "" {
		config.OpenVPNConfigDir = "/etc/openvpn"
	}

	return installServiceConfig("openvpn-go-api.conf.j2", constants.SupervisorWebConfigPath, config)
}

// InstallFrontendServiceConfig 安装前端 (Nginx) 服务配置
func InstallFrontendServiceConfig(autoStart bool) error {
	return installServiceConfig("openvpn-frontend.conf.j2", constants.SupervisorFrontendConfigPath, ServiceConfig{
		BinaryPath:       "/usr/sbin/nginx",
		WorkingDirectory: "/app",
		AutoStart:        autoStart,
	})
}

// installServiceConfig 通用的服务配置安装函数
func installServiceConfig(templateName, configPath string, config ServiceConfig) error {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 模板数据
	data := map[string]interface{}{
		"BinaryPath":       config.BinaryPath,
		"WorkingDirectory": config.WorkingDirectory,
		"Port":             config.Port,
		"DBPath":           config.DBPath,
		"OpenVPNConfigDir": config.OpenVPNConfigDir,
		"AutoStart":        config.AutoStart,
	}

	// 解析模板
	templatePath := filepath.Join(wd, "template", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("解析服务配置模板失败: %v", err)
	}

	// 生成配置文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染服务配置模板失败: %v", err)
	}

	// 确保配置目录存在
	if err := os.MkdirAll(constants.SupervisorConfDir, 0755); err != nil {
		return fmt.Errorf("创建 supervisor 配置目录失败: %v", err)
	}

	// 写入服务配置文件
	if err := os.WriteFile(configPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入服务配置文件失败: %v", err)
	}

	fmt.Printf("服务配置已安装: %s\n", configPath)
	return nil
}

// UpdateWebServiceConfig 更新 Web 服务配置并重新加载
func UpdateWebServiceConfig(config ServiceConfig) error {
	// 安装新配置
	if err := InstallWebServiceConfig(config); err != nil {
		return err
	}

	// 重新加载配置
	if err := SupervisorctlReload(); err != nil {
		return fmt.Errorf("重新加载 supervisor 配置失败: %v", err)
	}

	return nil
}

// UpdateOpenVPNServiceConfig 更新 OpenVPN 服务配置并重新加载
func UpdateOpenVPNServiceConfig(autoStart bool) error {
	// 安装新配置
	if err := InstallOpenVPNServiceConfig(autoStart); err != nil {
		return err
	}

	// 重新加载配置
	if err := SupervisorctlReload(); err != nil {
		return fmt.Errorf("重新加载 supervisor 配置失败: %v", err)
	}

	return nil
}

// RemoveServiceConfig 移除服务配置文件
func RemoveServiceConfig(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // 文件不存在，无需删除
	}

	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("删除服务配置文件失败: %v", err)
	}

	// 重新加载配置
	if err := SupervisorctlReload(); err != nil {
		return fmt.Errorf("重新加载 supervisor 配置失败: %v", err)
	}

	fmt.Printf("服务配置已移除: %s\n", configPath)
	return nil
}

// RemoveWebServiceConfig 移除 Web 服务配置
func RemoveWebServiceConfig() error {
	return RemoveServiceConfig(constants.SupervisorWebConfigPath)
}

// RemoveOpenVPNServiceConfig 移除 OpenVPN 服务配置
func RemoveOpenVPNServiceConfig() error {
	return RemoveServiceConfig(constants.SupervisorOpenVPNConfigPath)
}

// RemoveFrontendServiceConfig 移除前端服务配置
func RemoveFrontendServiceConfig() error {
	return RemoveServiceConfig(constants.SupervisorFrontendConfigPath)
}

// GetSupervisorConfigPath 获取 supervisor 主配置文件路径
func GetSupervisorConfigPath() string {
	return constants.SupervisorConfigPath
}

// IsSupervisorConfigExists 检查 supervisor 主配置文件是否存在
func IsSupervisorConfigExists() bool {
	_, err := os.Stat(GetSupervisorConfigPath())
	return err == nil
}

// IsWebServiceConfigExists 检查 Web 服务配置文件是否存在
func IsWebServiceConfigExists() bool {
	_, err := os.Stat(constants.SupervisorWebConfigPath)
	return err == nil
}

// IsOpenVPNServiceConfigExists 检查 OpenVPN 服务配置文件是否存在
func IsOpenVPNServiceConfigExists() bool {
	_, err := os.Stat(constants.SupervisorOpenVPNConfigPath)
	return err == nil
}

// IsFrontendServiceConfigExists 检查前端服务配置文件是否存在
func IsFrontendServiceConfigExists() bool {
	_, err := os.Stat(constants.SupervisorFrontendConfigPath)
	return err == nil
}

// BackupSupervisorConfig 备份当前的 supervisor 配置
func BackupSupervisorConfig() error {
	configPath := GetSupervisorConfigPath()
	if !IsSupervisorConfigExists() {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	backupPath := configPath + ".backup"
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return fmt.Errorf("创建备份文件失败: %v", err)
	}

	fmt.Printf("配置文件已备份到: %s\n", backupPath)
	return nil
}

// RestoreSupervisorConfig 恢复 supervisor 配置
func RestoreSupervisorConfig() error {
	configPath := GetSupervisorConfigPath()
	backupPath := configPath + ".backup"

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupPath)
	}

	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %v", err)
	}

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		return fmt.Errorf("恢复配置文件失败: %v", err)
	}

	// 重新加载配置
	if err := SupervisorctlReload(); err != nil {
		return fmt.Errorf("重新加载配置失败: %v", err)
	}

	fmt.Printf("配置文件已从备份恢复: %s\n", backupPath)
	return nil
}
