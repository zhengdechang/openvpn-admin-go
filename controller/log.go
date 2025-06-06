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