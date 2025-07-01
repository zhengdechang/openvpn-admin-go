package router

import (
	"openvpn-admin-go/controller"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// SetupServerRoutes 设置服务器相关路由
func SetupServerRoutes(r *gin.RouterGroup) {
	serverCtrl := &controller.ServerController{}
	server := r.Group("/server")
	server.Use(middleware.JWTAuthMiddleware())
	{
		// 列出服务器列表: 所有认证用户
		server.GET("/list", serverCtrl.ListServers)
		// 查看服务器状态: 所有认证用户
		server.GET("/status", serverCtrl.GetServerStatus)
		// 管理操作: superadmin
		super := server.Group("")
		super.Use(middleware.RoleRequired(string(model.RoleSuperAdmin)))
		{
			super.PUT("/update", serverCtrl.UpdateServer)
			super.DELETE("/delete", serverCtrl.DeleteServer)
			super.POST("/start", serverCtrl.StartServer)
			super.POST("/stop", serverCtrl.StopServer)
			super.POST("/restart", serverCtrl.RestartServer)
			super.GET("/config/template", serverCtrl.GetServerConfigTemplate)
			super.PUT("/config", serverCtrl.UpdateServerConfig)
			super.PUT("/port", serverCtrl.UpdatePort)
			// 配置项管理
			super.GET("/config/items", serverCtrl.GetConfigItems)
			super.PUT("/config/items", serverCtrl.UpdateConfigItems)
			super.PUT("/config/item/:key", serverCtrl.UpdateConfigItem)
		}
	}
}
