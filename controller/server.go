package controller

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/openvpn/server"

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
	// 检查服务是否运行
	cmd := exec.Command("systemctl", "is-active", constants.ServiceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("检查服务状态失败: %v", err)
	}

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
		if content, err := os.ReadFile(constants.ServerStatusLogPath); err == nil {
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
	// 更新服务器配置
	if err := openvpn.UpdateServerConfig(); err != nil {
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
	if err := openvpn.UpdateServerConfig(); err != nil {
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
	if err := server.UpdatePort(port.Port); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Port updated successfully"})
}
