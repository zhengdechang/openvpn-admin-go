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

       // Client Status Routes
       statusGroup := client.Group("/status")
       statusGroup.Use(middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)))
       {
           statusGroup.GET("", clientCtrl.GetAllClientStatuses)      // Existing: /client/status
           statusGroup.GET("/live", clientCtrl.GetLiveConnections) // New: /client/status/live
       }
       // Individual client status - this route might need role protection too, depending on policy
       client.GET("/status/:username", clientCtrl.GetClientStatus) // Path: /client/status/:username

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