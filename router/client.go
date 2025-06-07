package router

import (
	"openvpn-admin-go/controller"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// SetupClientRoutes 设置客户端相关路由
func SetupClientRoutes(r *gin.RouterGroup) {
	clientCtrl := &controller.ClientController{}
	client := r.Group("/client")
	client.Use(middleware.JWTAuthMiddleware())
	{
		// 客户端配置下载: all roles, but users only own
		client.GET("/config/:username", clientCtrl.GetClientConfig)
		// 客户端状态查看
		client.GET("/status", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.GetAllClientStatuses)
		client.GET("/status/:username", clientCtrl.GetClientStatus)
		// 管理操作: superadmin, admin, manager
		admin := client.Group("")
		admin.Use(middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)))
		{
			admin.GET("/list", clientCtrl.GetClientList)
			admin.POST("/add", clientCtrl.AddClient)
			admin.PUT("/update", clientCtrl.UpdateClient)
			admin.DELETE("/delete/:username", clientCtrl.DeleteClient)
			admin.POST("/revoke", clientCtrl.RevokeClient)
			admin.POST("/renew", clientCtrl.RenewClient)
		}
	}
}
