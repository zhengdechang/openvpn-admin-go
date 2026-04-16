package controller

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"openvpn-admin-go/common"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/openvpn"

	"github.com/gin-gonic/gin"
)

// LogController 日志查询
type LogController struct{}

// GetServerLogs 获取服务器日志
func (c *LogController) GetServerLogs(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	if claims.Role != string("superadmin") && claims.Role != string("admin") {
		common.Forbidden(ctx, "forbidden")
		return
	}

	cfg, err := openvpn.LoadConfig()
	if err != nil {
		common.InternalError(ctx, "failed to load OpenVPN config")
		return
	}

	logging.Debug("OpenVPNStatusLogPath: %s", cfg.OpenVPNStatusLogPath)

	if cfg.OpenVPNStatusLogPath == "" {
		common.InternalError(ctx, "OpenVPN status log path not configured")
		return
	}

	data, err := os.ReadFile(cfg.OpenVPNStatusLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			common.OK(ctx, gin.H{"logs": "OpenVPN is not running or status log not yet generated."})
			return
		}
		common.InternalError(ctx, err.Error())
		return
	}
	common.OK(ctx, gin.H{"logs": string(data)})
}

// GetClientLogs 获取客户端日志，支持分页
func (c *LogController) GetClientLogs(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	if claims.Role != string("superadmin") && claims.Role != string("admin") {
		common.Forbidden(ctx, "forbidden")
		return
	}

	cfg, err := openvpn.LoadConfig()
	if err != nil {
		common.InternalError(ctx, "failed to load OpenVPN config")
		return
	}

	logging.Debug("OpenVPNLogPath: %s", cfg.OpenVPNLogPath)

	if cfg.OpenVPNLogPath == "" {
		common.InternalError(ctx, "OpenVPN client log path not configured")
		return
	}

	offsetStr := ctx.DefaultQuery("offset", "0")
	limitStr := ctx.DefaultQuery("limit", "1000")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		common.BadRequest(ctx, "invalid offset parameter")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		common.BadRequest(ctx, "invalid limit parameter")
		return
	}

	file, err := os.Open(cfg.OpenVPNLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			common.OK(ctx, gin.H{
				"logs":       "OpenVPN is not running or client log not yet generated.",
				"totalLines": 0,
				"offset":     0,
				"limit":      limit,
				"hasMore":    false,
			})
			return
		}
		common.InternalError(ctx, err.Error())
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		common.InternalError(ctx, err.Error())
		return
	}

	// Reverse so newest entries appear first
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	totalLines := len(lines)

	if offset < 0 {
		offset = totalLines + offset
		if offset < 0 {
			offset = 0
		}
	}

	start := offset
	end := offset + limit

	if start >= totalLines {
		common.OK(ctx, gin.H{
			"logs":       "",
			"totalLines": totalLines,
			"offset":     offset,
			"limit":      limit,
			"hasMore":    false,
		})
		return
	}

	if end > totalLines {
		end = totalLines
	}

	logsContent := strings.Join(lines[start:end], "\n")
	common.OK(ctx, gin.H{
		"logs":       logsContent,
		"totalLines": totalLines,
		"offset":     offset,
		"limit":      limit,
		"hasMore":    end < totalLines,
	})
}
