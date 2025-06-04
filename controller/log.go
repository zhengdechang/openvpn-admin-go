package controller

import (
   "net/http"
   "os"
   "strings"

   "github.com/gin-gonic/gin"
   "openvpn-admin-go/constants"
   "openvpn-admin-go/middleware"
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
   data, err := os.ReadFile(constants.ServerStatusLogPath)
   if err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"logs": string(data)})
}

// GetClientLogs 获取客户端日志 (过滤服务器日志中对应用户名)
func (c *LogController) GetClientLogs(ctx *gin.Context) {
   username := ctx.Param("username")
   claims := ctx.MustGet("claims").(*middleware.Claims)
   // 普通用户仅能查看自己
   if claims.Role == string("user") && claims.UserID != username {
       ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
       return
   }
   // 其他角色: superadmin/admin/manager 可查看
   data, err := os.ReadFile(constants.ServerStatusLogPath)
   if err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   lines := strings.Split(string(data), "\n")
   var filtered []string
   for _, line := range lines {
       if strings.Contains(line, username) {
           filtered = append(filtered, line)
       }
   }
   ctx.JSON(http.StatusOK, gin.H{"logs": filtered})
}