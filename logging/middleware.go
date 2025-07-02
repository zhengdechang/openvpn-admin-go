package logging

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装gin.ResponseWriter以捕获响应数据
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// GinLoggingMiddleware 返回一个Gin中间件用于记录API请求
func GinLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 读取请求体（如果需要记录的话）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = responseWriter

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取错误信息
		errorMsg := ""
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// 记录API请求日志
		LogAPIRequest(c, c.Writer.Status(), duration, errorMsg)

		// 如果是敏感操作，记录详细信息
		if isSensitiveOperation(c) {
			logSensitiveOperation(c, requestBody, responseWriter.body.Bytes())
		}
	}
}

// isSensitiveOperation 判断是否为敏感操作
func isSensitiveOperation(c *gin.Context) bool {
	method := c.Request.Method
	path := c.Request.URL.Path

	// 定义敏感操作的路径模式
	sensitivePatterns := []string{
		"/api/auth/login",
		"/api/auth/logout",
		"/api/users",
		"/api/client",
		"/api/server",
		"/api/departments",
	}

	// 检查是否为POST、PUT、DELETE操作
	if method == "POST" || method == "PUT" || method == "DELETE" {
		for _, pattern := range sensitivePatterns {
			if contains(path, pattern) {
				return true
			}
		}
	}

	return false
}

// logSensitiveOperation 记录敏感操作的详细信息
func logSensitiveOperation(c *gin.Context, requestBody, responseBody []byte) {
	username := "anonymous"
	if claims, exists := c.Get("claims"); exists {
		if claimsData, ok := claims.(map[string]interface{}); ok {
			if user, ok := claimsData["username"].(string); ok {
				username = user
			}
		}
	}

	action := c.Request.Method
	resource := c.Request.URL.Path
	details := ""

	// 根据不同的操作类型记录不同的详细信息
	switch {
	case contains(resource, "/auth/login"):
		LogSecurityEvent("LOGIN_ATTEMPT", username, c.ClientIP(), "User login attempt")
	case contains(resource, "/auth/logout"):
		LogSecurityEvent("LOGOUT", username, c.ClientIP(), "User logout")
	case contains(resource, "/users"):
		if action == "POST" {
			details = "Create new user"
		} else if action == "PUT" {
			details = "Update user information"
		} else if action == "DELETE" {
			details = "Delete user"
		}
		LogUserAction(username, action, "USER", details)
	case contains(resource, "/client"):
		if action == "POST" {
			details = "Create new client certificate"
		} else if action == "PUT" {
			details = "Update client configuration"
		} else if action == "DELETE" {
			details = "Revoke client certificate"
		}
		LogUserAction(username, action, "CLIENT", details)
	case contains(resource, "/server"):
		if action == "POST" {
			details = "Start/Stop OpenVPN server"
		} else if action == "PUT" {
			details = "Update server configuration"
		}
		LogUserAction(username, action, "SERVER", details)
	case contains(resource, "/departments"):
		if action == "POST" {
			details = "Create new department"
		} else if action == "PUT" {
			details = "Update department information"
		} else if action == "DELETE" {
			details = "Delete department"
		}
		LogUserAction(username, action, "DEPARTMENT", details)
	}
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && s[:len(substr)] == substr))
}

// RequestLoggingConfig 请求日志配置
type RequestLoggingConfig struct {
	EnableRequestBody  bool
	EnableResponseBody bool
	MaxBodySize        int64
	SkipPaths          []string
}

// GinDetailedLoggingMiddleware 返回一个详细的Gin日志中间件
func GinDetailedLoggingMiddleware(config RequestLoggingConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过此路径
		for _, skipPath := range config.SkipPaths {
			if c.Request.URL.Path == skipPath {
				c.Next()
				return
			}
		}

		startTime := time.Now()

		// 记录请求详细信息
		Debug("Request Details - Method: %s, Path: %s, Query: %s, Headers: %v, IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.URL.RawQuery,
			c.Request.Header,
			c.ClientIP())

		// 读取请求体（如果启用）
		var requestBody []byte
		if config.EnableRequestBody && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			if int64(len(requestBody)) <= config.MaxBodySize {
				Debug("Request Body: %s", string(requestBody))
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器（如果需要记录响应体）
		var respWriter *responseWriter
		if config.EnableResponseBody {
			respWriter = &responseWriter{
				ResponseWriter: c.Writer,
				body:           bytes.NewBufferString(""),
			}
			c.Writer = respWriter
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录响应详细信息
		if config.EnableResponseBody && respWriter != nil {
			responseBody := respWriter.body.Bytes()
			if int64(len(responseBody)) <= config.MaxBodySize {
				Debug("Response Body: %s", string(responseBody))
			}
		}

		// 记录处理结果
		Debug("Request Completed - Status: %d, Duration: %v", c.Writer.Status(), duration)

		// 获取错误信息
		errorMsg := ""
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// 记录API请求日志
		LogAPIRequest(c, c.Writer.Status(), duration, errorMsg)
	}
}
