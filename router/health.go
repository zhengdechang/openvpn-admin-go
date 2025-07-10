package router

import (
	"openvpn-admin-go/controller"

	"github.com/gin-gonic/gin"
)

// SetupHealthRoutes 设置健康检查路由
func SetupHealthRoutes(r *gin.RouterGroup) {
	healthCtrl := &controller.HealthController{}

	// 健康检查路由 - 不需要认证
	r.GET("/health", healthCtrl.HealthCheck)
	r.GET("/health/readiness", healthCtrl.ReadinessCheck)
	r.GET("/health/liveness", healthCtrl.LivenessCheck)
}
