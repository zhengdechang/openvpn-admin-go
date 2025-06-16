package controller

import (
   "net/http"
   "os"
   "fmt"
   "bufio"   // Added for line-by-line reading
   "math"    // Added for pagination
   "strconv" // Added for pagination

   "github.com/gin-gonic/gin"
   "openvpn-admin-go/middleware"
   "openvpn-admin-go/openvpn"
)

// LogController 日志查询
type LogController struct{}

// paginateLogFile reads a file line by line and returns a specific page of logs and total line count.
func paginateLogFile(filePath string, page int, pageSize int) (lines []string, totalItems int64, totalPages int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// First pass: count total lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalItems++
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("error reading log file for counting lines: %w", err)
	}

	totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	if totalPages == 0 && totalItems > 0 {
		totalPages = 1
	}
    // Ensure page is within valid range
    if page < 1 {
        page = 1
    }
    if page > totalPages && totalPages > 0 { // if totalPages is 0, page 1 is still valid (empty result)
        page = totalPages
    }


	// Reset file pointer to the beginning for the second pass
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to seek to beginning of log file: %w", err)
	}

	// Second pass: read the specific page
	scanner = bufio.NewScanner(file) // Re-initialize scanner
	currentLine := int64(0)
	startLine := int64((page - 1) * pageSize)

	for scanner.Scan() {
		if currentLine >= startLine && currentLine < startLine+int64(pageSize) {
			lines = append(lines, scanner.Text())
		}
		currentLine++
		if currentLine >= startLine+int64(pageSize) && len(lines) >= pageSize { // Optimization: stop if page is full
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("error reading log file for pagination: %w", err)
	}

	return lines, totalItems, totalPages, nil
}

// GetServerLogs 获取服务器日志
func (c *LogController) GetServerLogs(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	if claims.Role != string("superadmin") && claims.Role != string("admin") {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	pageQuery := ctx.DefaultQuery("page", "1")
	pageSizeQuery := ctx.DefaultQuery("pageSize", "100")

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if parsing fails or value is invalid
	}

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize < 1 {
		pageSize = 100 // Default to 100 if parsing fails or value is invalid
	}
	if pageSize > 1000 { // Max pageSize 1000
		pageSize = 1000
	}

	cfg, err := openvpn.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load OpenVPN config: " + err.Error()})
		return
	}
	if cfg.OpenVPNStatusLogPath == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OpenVPN status log path not configured"})
		return
	}

	logs, totalItems, totalPages, actualPageUsed, err := paginateLogFile(cfg.OpenVPNStatusLogPath, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"logs":        logs,
		"totalItems":  totalItems,
		"totalPages":  totalPages,
		"currentPage": actualPageUsed,
		"pageSize":    pageSize,
	})
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

   fmt.Printf("OpenVPNLogPath: %s\n", cfg.OpenVPNLogPath)

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