package controller

import (
	"net/http"
	"openvpn-admin-go/openvpn"

	"github.com/gin-gonic/gin"
)

type ServerController struct{}

// GetServerConfigTemplate 获取服务端配置模板
func (c *ServerController) GetServerConfigTemplate(ctx *gin.Context) {
	template := openvpn.GetServerConfigTemplate()
	ctx.JSON(http.StatusOK, gin.H{
		"template": template,
	})
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

	if err := openvpn.UpdateServerConfig(config.Config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "配置更新成功"})
}

// RestartServer 重启OpenVPN服务
func (c *ServerController) RestartServer(ctx *gin.Context) {
	if err := openvpn.RestartServer(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "服务重启成功"})
} 