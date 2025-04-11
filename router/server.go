package router

import (
	"openvpn-admin-go/controller"

	"github.com/gin-gonic/gin"
)

// SetupServerRoutes 设置服务端路由
func SetupServerRoutes(r *gin.RouterGroup) {
	serverController := &controller.ServerController{}

	// 服务端相关路由
	serverGroup := r.Group("/server")
	{
		serverGroup.GET("/config/template", serverController.GetServerConfigTemplate)
		serverGroup.PUT("/config", serverController.UpdateServerConfig)
		serverGroup.POST("/restart", serverController.RestartServer)
	}
} 