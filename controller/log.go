package controller

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	logging.Debug("OpenVPNStatusLogPath: %s", cfg.OpenVPNStatusLogPath)

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

// GetClientLogs 获取客户端日志，支持分页
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

	logging.Debug("OpenVPNLogPath: %s", cfg.OpenVPNLogPath)

	if cfg.OpenVPNLogPath == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OpenVPN client log path not configured"})
		return
	}

	// 获取分页参数
	offsetStr := ctx.DefaultQuery("offset", "0")
	limitStr := ctx.DefaultQuery("limit", "1000")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	// 读取日志文件并按行分页
	file, err := os.Open(cfg.OpenVPNLogPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// 读取所有行
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalLines := len(lines)

	// 如果offset为负数，表示从末尾开始计算
	if offset < 0 {
		offset = totalLines + offset
		if offset < 0 {
			offset = 0
		}
	}

	// 计算实际的起始和结束位置
	start := offset
	end := offset + limit

	if start >= totalLines {
		// 如果起始位置超出范围，返回空结果
		ctx.JSON(http.StatusOK, gin.H{
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

	// 提取指定范围的行
	selectedLines := lines[start:end]
	logsContent := strings.Join(selectedLines, "\n")

	ctx.JSON(http.StatusOK, gin.H{
		"logs":       logsContent,
		"totalLines": totalLines,
		"offset":     offset,
		"limit":      limit,
		"hasMore":    end < totalLines,
	})
}
