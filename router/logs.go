package router

import (
   "openvpn-admin-go/controller"
   "openvpn-admin-go/middleware"

   "github.com/gin-gonic/gin"
)

// SetupLogRoutes 设置日志查询路由
func SetupLogRoutes(r *gin.RouterGroup) {
   logCtrl := &controller.LogController{} // Keep this if GetServerLogs is still used from LogController
   clientLogCtrl := &controller.ClientLogController{} // Assuming a similar pattern or direct use

   logs := r.Group("/logs")
   logs.Use(middleware.JWTAuthMiddleware())
   {
       logs.GET("/server", logCtrl.GetServerLogs)
       // Routes for the new ClientLog controller
       logs.POST("/client", controller.CreateClientLog) // Direct function reference
       logs.GET("/client", controller.GetClientLogs)    // Direct function reference
   }
}