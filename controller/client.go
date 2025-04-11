package controller

import (
	"net/http"
	"openvpn-admin-go/openvpn"

	"github.com/gin-gonic/gin"
)

type ClientController struct{}

// GenerateClientConfig 生成客户端配置文件
func (c *ClientController) GenerateClientConfig(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	configPath, err := openvpn.GenerateClientConfig(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.File(configPath)
}

// DeleteClient 删除客户端并吊销证书
func (c *ClientController) DeleteClient(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	if err := openvpn.RevokeClientCert(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户删除成功，证书已吊销"})
}

// GetClientStatus 获取客户端状态
func (c *ClientController) GetClientStatus(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	status := openvpn.GetClientStatus(username)
	if status == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, status)
}

// GetAllClientStatuses 获取所有客户端状态
func (c *ClientController) GetAllClientStatuses(ctx *gin.Context) {
	statuses := openvpn.GetAllClientStatuses()
	ctx.JSON(http.StatusOK, statuses)
}

// PauseClient 暂停客户端连接
func (c *ClientController) PauseClient(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	if err := openvpn.PauseClient(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "客户端已暂停"})
}

// ResumeClient 恢复客户端连接
func (c *ClientController) ResumeClient(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	if err := openvpn.ResumeClient(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "客户端已恢复"})
} 