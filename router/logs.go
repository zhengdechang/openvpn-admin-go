package router

import (
   "openvpn-admin-go/controller"
   "openvpn-admin-go/middleware"

   "github.com/gin-gonic/gin"
)

// SetupLogRoutes 设置日志查询路由
func SetupLogRoutes(r *gin.RouterGroup) {
   logCtrl := &controller.LogController{}
   logs := r.Group("/logs")
   logs.Use(middleware.JWTAuthMiddleware())
   {
       logs.GET("/server", logCtrl.GetServerLogs)
       logs.GET("/client/:username", logCtrl.GetClientLogs)
   }
}