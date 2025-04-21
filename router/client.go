package router

import (
	"openvpn-admin-go/controller"

	"github.com/gin-gonic/gin"
)

// SetupClientRoutes 设置客户端相关路由
func SetupClientRoutes(r *gin.RouterGroup) {
	clientCtrl := &controller.ClientController{}
	client := r.Group("/client")
	{
		client.GET("/list", clientCtrl.GetClientList)
		client.POST("/add", clientCtrl.AddClient)
		client.PUT("/update", clientCtrl.UpdateClient)
		client.DELETE("/delete/:username", clientCtrl.DeleteClient)
		client.GET("/config/:username", clientCtrl.GetClientConfig)
		client.POST("/revoke", clientCtrl.RevokeClient)
		client.POST("/renew", clientCtrl.RenewClient)
	}
} 