package controller

import (
   "net/http"
   "os"

   "github.com/gin-gonic/gin"
   "openvpn-admin-go/middleware"
   "openvpn-admin-go/openvpn"
)

// LogController 日志查询
type LogController struct{}

// GetServerLogs 获取服务器日志
func (c *LogController) GetServerLogs(ctx *gin.Context) {
   claims := ctx.MustGet("claims").(*middleware.Claims)
   // 仅 superadmin 和 admin
   if claims.Role != string("superadmin") && claims.Role != string("admin") {
       ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
       return
   }

   // 从配置中获取 status 日志路径
   cfg, err := openvpn.LoadConfig()
   if err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load OpenVPN config"})
       return
   }

   if cfg.OpenVPNStatusLogPath == "" {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OpenVPN status log path not configured"})
       return
   }

   data, err := os.ReadFile(cfg.OpenVPNStatusLogPath)
   if err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"logs": string(data)})
}

// GetClientLogs 获取客户端日志
func (c *LogController) GetClientLogs(ctx *gin.Context) {
   claims := ctx.MustGet("claims").(*middleware.Claims)
   // 仅 superadmin 和 admin
   if claims.Role != string("superadmin") && claims.Role != string("admin") {
       ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
       return
   }

   // 从配置中获取客户端日志路径
   cfg, err := openvpn.LoadConfig()
   if err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load OpenVPN config"})
       return
   }

   if cfg.OpenVPNLogPath == "" {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OpenVPN client log path not configured"})
       return
   }

   data, err := os.ReadFile(cfg.OpenVPNLogPath)
   if err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"logs": string(data)})
}