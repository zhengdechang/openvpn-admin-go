package controller

import (
	"net/http"
	"openvpn-admin-go/database"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// HealthCheck 健康检查端点
func (c *HealthController) HealthCheck(ctx *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   AppVersion(), // 来自环境变量 APP_VERSION，缺省 "1.0.0"
		Services:  make(map[string]string),
	}

	// 检查数据库连接
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err != nil {
			response.Services["database"] = "error"
			response.Status = "unhealthy"
		} else if err := sqlDB.Ping(); err != nil {
			response.Services["database"] = "error"
			response.Status = "unhealthy"
		} else {
			response.Services["database"] = "healthy"
		}
	} else {
		response.Services["database"] = "unavailable"
		response.Status = "unhealthy"
	}

	// 根据整体状态返回相应的HTTP状态码
	if response.Status == "healthy" {
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusServiceUnavailable, response)
	}
}

// ReadinessCheck 就绪检查端点
func (c *HealthController) ReadinessCheck(ctx *gin.Context) {
	// 检查应用是否准备好接收请求
	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}

// LivenessCheck 存活检查端点
func (c *HealthController) LivenessCheck(ctx *gin.Context) {
	// 简单的存活检查
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, response)
}
