package router

import (
	"openvpn-admin-go/controller"

	"github.com/gin-gonic/gin"
)

// SetupClientRoutes 设置客户端路由
func SetupClientRoutes(r *gin.RouterGroup) {
	clientController := &controller.ClientController{}

	// 客户端相关路由
	clientGroup := r.Group("/client")
	{
		clientGroup.GET("/:username/config", clientController.GenerateClientConfig)
		clientGroup.DELETE("/:username", clientController.DeleteClient)
		clientGroup.GET("/:username/status", clientController.GetClientStatus)
		clientGroup.GET("/status", clientController.GetAllClientStatuses)
		clientGroup.POST("/:username/pause", clientController.PauseClient)
		clientGroup.POST("/:username/resume", clientController.ResumeClient)
	}
} 