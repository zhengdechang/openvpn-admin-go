package router

import (
	"openvpn-admin-go/controller"

	"github.com/gin-gonic/gin"
)

// SetupServerRoutes 设置服务器相关路由
func SetupServerRoutes(r *gin.RouterGroup) {
	serverCtrl := &controller.ServerController{}
	server := r.Group("/server")
	{
		server.GET("/list", serverCtrl.GetServerList)
		server.POST("/add", serverCtrl.AddServer)
		server.PUT("/update", serverCtrl.UpdateServer)
		server.DELETE("/delete", serverCtrl.DeleteServer)
		server.GET("/status", serverCtrl.GetServerStatus)
		server.POST("/start", serverCtrl.StartServer)
		server.POST("/stop", serverCtrl.StopServer)
		server.POST("/restart", serverCtrl.RestartServer)
		server.GET("/config/template", serverCtrl.GetServerConfigTemplate)
		server.PUT("/config", serverCtrl.UpdateServerConfig)
	}
} 