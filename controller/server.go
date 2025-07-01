package controller

import (
	"fmt"
	"net/http"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/openvpn"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerController struct{}

// ListServers 列出服务器列表
func (c *ServerController) ListServers(ctx *gin.Context) {
	// 加载当前配置
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 获取运行状态
	status, err := GetServerStatus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 构造返回结构
	server := struct {
		Port      int    `json:"port"`
		Protocol  string `json:"protocol"`
		Network   string `json:"network"`
		Netmask   string `json:"netmask"`
		Status    string `json:"status"`
		Uptime    string `json:"uptime"`
		Connected int    `json:"connected"`
		Total     int    `json:"total"`
	}{
		Port:      cfg.OpenVPNPort,
		Protocol:  cfg.OpenVPNProto,
		Network:   cfg.OpenVPNServerNetwork,
		Netmask:   cfg.OpenVPNServerNetmask,
		Status:    status.Status,
		Uptime:    status.Uptime,
		Connected: status.Connected,
		Total:     status.Total,
	}
	// 返回数组格式
	ctx.JSON(http.StatusOK, []interface{}{server})
}

// ServerStatus 服务器状态
// ServerStatus 服务器状态
type ServerStatus struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Uptime      string `json:"uptime"`      // 运行时长
	Connected   int    `json:"connected"`   // 当前已连接数
	Total       int    `json:"total"`       // 历史总连接数
	LastUpdated string `json:"lastUpdated"` // 最后更新时间
}

// GetServerStatus 获取服务器状态
func GetServerStatus() (*ServerStatus, error) {
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenVPN config: %w", err)
	}

	// 检查服务是否运行
	// 检查服务是否运行，忽略非零退出码，获取服务状态字符串
	cmd := exec.Command("systemctl", "is-active", constants.ServiceName)
	output, _ := cmd.CombinedOutput()

	status := &ServerStatus{
		Name:        "server",
		Status:      strings.TrimSpace(string(output)),
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	// 如果服务正在运行，获取更多信息
	if status.Status == "active" {
		// 获取服务启动时间
		cmd = exec.Command("systemctl", "show", constants.ServiceName, "--property=ActiveEnterTimestamp")
		if output, err := cmd.CombinedOutput(); err == nil {
			if t0, err := time.Parse("Mon 2006-01-02 15:04:05 MST", strings.TrimSpace(strings.TrimPrefix(string(output), "ActiveEnterTimestamp="))); err == nil {
				status.Uptime = time.Since(t0).String()
			}
		}

		// 获取连接数
		if content, err := os.ReadFile(cfg.OpenVPNStatusLogPath); err == nil {
			lines := strings.Split(string(content), "\n")
			status.Total = len(lines)
			status.Connected = 0
			for _, line := range lines {
				if strings.Contains(line, "CONNECTED") {
					status.Connected++
				}
			}
		}
	}

	return status, nil
}

// UpdateServer 更新服务器
func (c *ServerController) UpdateServer(ctx *gin.Context) {
	var server struct {
		Port     int    `json:"port" binding:"required"`
		Protocol string `json:"protocol" binding:"required"`
		Network  string `json:"network" binding:"required"`
		Netmask  string `json:"netmask" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 使用 openvpn/server 包处理参数更新与服务重启
	if err := openvpn.ConfigureServer(server.Port, server.Protocol, server.Network, server.Netmask); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Server updated successfully"})
}

// DeleteServer 删除服务器
func (c *ServerController) DeleteServer(ctx *gin.Context) {
	// 目前不支持删除服务器
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "目前不支持删除服务器"})
}

// GetServerStatus 获取服务器状态
func (c *ServerController) GetServerStatus(ctx *gin.Context) {
	status, err := GetServerStatus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, status)
}

// StartServer 启动服务器
func (c *ServerController) StartServer(ctx *gin.Context) {
	// 启动服务器
	if err := openvpn.RestartServer(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Server started successfully"})
}

// StopServer 停止服务器
func (c *ServerController) StopServer(ctx *gin.Context) {
	// 停止服务器
	cmd := exec.Command("systemctl", "stop", constants.ServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("停止服务失败: %v\n输出: %s", err, string(output))})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Server stopped successfully"})
}

// RestartServer 重启服务器
func (c *ServerController) RestartServer(ctx *gin.Context) {
	if err := openvpn.RestartServer(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Server restarted successfully"})
}

// GetServerConfigTemplate 获取服务器配置模板
func (c *ServerController) GetServerConfigTemplate(ctx *gin.Context) {
	template, err := openvpn.GetServerConfigTemplate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"template": template})
}

// UpdateServerConfig 更新服务器配置
func (c *ServerController) UpdateServerConfig(ctx *gin.Context) {
	var config struct {
		Config string `json:"config" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 使用 openvpn/server 包写入自定义配置并重启服务
	if err := openvpn.ApplyServerConfig(config.Config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Server config updated successfully"})
}

// UpdatePort 更新服务器端口
func (c *ServerController) UpdatePort(ctx *gin.Context) {
	var port struct {
		Port int `json:"port" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&port); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新端口
	if err := openvpn.UpdatePort(port.Port); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Port updated successfully"})
}

// ConfigItem 配置项结构
type ConfigItem struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Type        string      `json:"type"`                 // text, number, boolean, select, array
	Label       string      `json:"label"`                // 显示标签
	Description string      `json:"description"`          // 描述
	Options     []string    `json:"options,omitempty"`    // 选择项（用于select类型）
	Required    bool        `json:"required"`             // 是否必填
	Validation  string      `json:"validation,omitempty"` // 验证规则
}

// GetConfigItems 获取配置项列表
func (c *ServerController) GetConfigItems(ctx *gin.Context) {
	// 加载当前配置
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "加载配置失败: " + err.Error()})
		return
	}

	// 构建配置项列表
	items := []ConfigItem{
		{
			Key:         "openvpn_port",
			Value:       cfg.OpenVPNPort,
			Type:        "number",
			Label:       "OpenVPN端口",
			Description: "OpenVPN服务监听端口",
			Required:    true,
			Validation:  "min:1,max:65535",
		},
		{
			Key:         "openvpn_proto",
			Value:       cfg.OpenVPNProto,
			Type:        "select",
			Label:       "协议类型",
			Description: "OpenVPN使用的协议",
			Options:     []string{"tcp", "tcp6", "udp", "udp6"},
			Required:    true,
		},
		{
			Key:         "openvpn_server_hostname",
			Value:       cfg.OpenVPNServerHostname,
			Type:        "text",
			Label:       "服务器主机名",
			Description: "客户端连接的服务器地址",
			Required:    true,
			Validation:  "ip_or_hostname",
		},
		{
			Key:         "openvpn_server_network",
			Value:       cfg.OpenVPNServerNetwork,
			Type:        "text",
			Label:       "服务器网络",
			Description: "VPN内部网络地址",
			Required:    true,
			Validation:  "ip",
		},
		{
			Key:         "openvpn_server_netmask",
			Value:       cfg.OpenVPNServerNetmask,
			Type:        "text",
			Label:       "子网掩码",
			Description: "VPN内部网络子网掩码",
			Required:    true,
			Validation:  "netmask",
		},
		{
			Key:         "openvpn_sync_certs",
			Value:       cfg.OpenVPNSyncCerts,
			Type:        "boolean",
			Label:       "同步证书",
			Description: "是否自动同步客户端证书",
			Required:    false,
		},
		{
			Key:         "openvpn_use_crl",
			Value:       cfg.OpenVPNUseCRL,
			Type:        "boolean",
			Label:       "使用CRL",
			Description: "是否使用证书撤销列表",
			Required:    false,
		},
		{
			Key:         "openvpn_client_to_client",
			Value:       cfg.OpenVPNClientToClient,
			Type:        "boolean",
			Label:       "客户端互通",
			Description: "允许客户端之间直接通信",
			Required:    false,
		},
		{
			Key:         "openvpn_routes",
			Value:       cfg.OpenVPNRoutes,
			Type:        "array",
			Label:       "路由配置",
			Description: "推送给客户端的路由列表",
			Required:    false,
		},
		{
			Key:         "dns_server_ip",
			Value:       cfg.DNSServerIP,
			Type:        "text",
			Label:       "DNS服务器IP",
			Description: "推送给客户端的DNS服务器地址",
			Required:    false,
			Validation:  "ip",
		},
		{
			Key:         "dns_server_domain",
			Value:       cfg.DNSServerDomain,
			Type:        "text",
			Label:       "DNS域名",
			Description: "推送给客户端的DNS域名",
			Required:    false,
		},
		{
			Key:         "openvpn_management_port",
			Value:       cfg.OpenVPNManagementPort,
			Type:        "number",
			Label:       "管理端口",
			Description: "OpenVPN管理接口端口",
			Required:    false,
			Validation:  "min:1,max:65535",
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"items": items})
}

// UpdateConfigItem 更新单个配置项
func (c *ServerController) UpdateConfigItem(ctx *gin.Context) {
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "配置项key不能为空"})
		return
	}

	var request struct {
		Value interface{} `json:"value" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加载当前配置
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "加载配置失败: " + err.Error()})
		return
	}

	// 更新指定的配置项
	if err := updateSingleConfigItem(cfg, key, request.Value); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存配置
	if err := openvpn.SaveConfig(cfg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存配置失败: " + err.Error()})
		return
	}

	// 重新生成服务器配置
	if err := openvpn.UpdateServerConfig(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务器配置失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "配置项更新成功"})
}

// UpdateConfigItems 批量更新配置项
func (c *ServerController) UpdateConfigItems(ctx *gin.Context) {
	var request struct {
		Items map[string]interface{} `json:"items" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加载当前配置
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "加载配置失败: " + err.Error()})
		return
	}

	// 批量更新配置项
	for key, value := range request.Items {
		if err := updateSingleConfigItem(cfg, key, value); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("更新配置项 %s 失败: %s", key, err.Error())})
			return
		}
	}

	// 保存配置
	if err := openvpn.SaveConfig(cfg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存配置失败: " + err.Error()})
		return
	}

	// 重新生成服务器配置
	if err := openvpn.UpdateServerConfig(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务器配置失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "配置项批量更新成功"})
}

// updateSingleConfigItem 更新单个配置项的辅助函数
func updateSingleConfigItem(cfg *openvpn.Config, key string, value interface{}) error {
	switch key {
	case "openvpn_port":
		if port, ok := value.(float64); ok {
			cfg.OpenVPNPort = int(port)
		} else {
			return fmt.Errorf("端口必须是数字")
		}
	case "openvpn_proto":
		if proto, ok := value.(string); ok {
			validProtos := []string{"tcp", "tcp6", "udp", "udp6"}
			valid := false
			for _, v := range validProtos {
				if proto == v {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("协议类型无效，必须是: %s", strings.Join(validProtos, ", "))
			}
			cfg.OpenVPNProto = proto
		} else {
			return fmt.Errorf("协议类型必须是字符串")
		}
	case "openvpn_server_hostname":
		if hostname, ok := value.(string); ok {
			if hostname == "" {
				return fmt.Errorf("服务器主机名不能为空")
			}
			cfg.OpenVPNServerHostname = hostname
		} else {
			return fmt.Errorf("服务器主机名必须是字符串")
		}
	case "openvpn_server_network":
		if network, ok := value.(string); ok {
			if network == "" {
				return fmt.Errorf("服务器网络不能为空")
			}
			cfg.OpenVPNServerNetwork = network
		} else {
			return fmt.Errorf("服务器网络必须是字符串")
		}
	case "openvpn_server_netmask":
		if netmask, ok := value.(string); ok {
			if netmask == "" {
				return fmt.Errorf("子网掩码不能为空")
			}
			cfg.OpenVPNServerNetmask = netmask
		} else {
			return fmt.Errorf("子网掩码必须是字符串")
		}
	case "openvpn_sync_certs":
		if sync, ok := value.(bool); ok {
			cfg.OpenVPNSyncCerts = sync
		} else {
			return fmt.Errorf("同步证书必须是布尔值")
		}
	case "openvpn_use_crl":
		if useCRL, ok := value.(bool); ok {
			cfg.OpenVPNUseCRL = useCRL
		} else {
			return fmt.Errorf("使用CRL必须是布尔值")
		}
	case "openvpn_client_to_client":
		if clientToClient, ok := value.(bool); ok {
			cfg.OpenVPNClientToClient = clientToClient
		} else {
			return fmt.Errorf("客户端互通必须是布尔值")
		}
	case "openvpn_routes":
		if routes, ok := value.([]interface{}); ok {
			stringRoutes := make([]string, len(routes))
			for i, route := range routes {
				if routeStr, ok := route.(string); ok {
					stringRoutes[i] = routeStr
				} else {
					return fmt.Errorf("路由配置必须是字符串数组")
				}
			}
			cfg.OpenVPNRoutes = stringRoutes
		} else {
			return fmt.Errorf("路由配置必须是数组")
		}
	case "dns_server_ip":
		if dnsIP, ok := value.(string); ok {
			cfg.DNSServerIP = dnsIP
		} else {
			return fmt.Errorf("DNS服务器IP必须是字符串")
		}
	case "dns_server_domain":
		if dnsDomain, ok := value.(string); ok {
			cfg.DNSServerDomain = dnsDomain
		} else {
			return fmt.Errorf("DNS域名必须是字符串")
		}
	case "openvpn_management_port":
		if port, ok := value.(float64); ok {
			cfg.OpenVPNManagementPort = int(port)
		} else {
			return fmt.Errorf("管理端口必须是数字")
		}
	default:
		return fmt.Errorf("未知的配置项: %s", key)
	}
	return nil
}
