package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// WebLogger Web日志记录器
type WebLogger struct {
	logger    *log.Logger
	level     LogLevel
	logFile   *os.File
	enableAPI bool // 是否启用API日志记录
}

var (
	// DefaultLogger 默认的Web日志记录器实例
	DefaultLogger *WebLogger
)

// Config 日志配置
type Config struct {
	LogLevel    LogLevel
	LogFilePath string
	EnableAPI   bool
	MaxFileSize int64 // 最大文件大小（字节）
}

// Init 初始化日志系统
func Init(config Config) error {
	// 确保日志目录存在
	logDir := filepath.Dir(config.LogFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开日志文件
	logFile, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 创建多重写入器（同时写入文件和控制台）
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 创建日志记录器
	logger := log.New(multiWriter, "", 0)

	DefaultLogger = &WebLogger{
		logger:    logger,
		level:     config.LogLevel,
		logFile:   logFile,
		enableAPI: config.EnableAPI,
	}

	return nil
}

// Close 关闭日志系统
func Close() error {
	if DefaultLogger != nil && DefaultLogger.logFile != nil {
		return DefaultLogger.logFile.Close()
	}
	return nil
}

// formatMessage 格式化日志消息
func (w *WebLogger) formatMessage(level LogLevel, message string) string {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(3)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level.String(), caller, message)
}

// log 内部日志记录方法
func (w *WebLogger) log(level LogLevel, format string, args ...interface{}) {
	if w == nil || level < w.level {
		return
	}

	message := fmt.Sprintf(format, args...)
	formattedMessage := w.formatMessage(level, message)
	w.logger.Println(formattedMessage)
}

// Debug 记录调试日志
func Debug(format string, args ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.log(DEBUG, format, args...)
	}
}

// Info 记录信息日志
func Info(format string, args ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.log(INFO, format, args...)
	}
}

// Warn 记录警告日志
func Warn(format string, args ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.log(WARN, format, args...)
	}
}

// Error 记录错误日志
func Error(format string, args ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.log(ERROR, format, args...)
	}
}

// Fatal 记录致命错误日志并退出程序
func Fatal(format string, args ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.log(FATAL, format, args...)
		os.Exit(1)
	}
}

// LogAPIRequest 记录API请求日志
func LogAPIRequest(c *gin.Context, statusCode int, duration time.Duration, errorMsg string) {
	if DefaultLogger == nil || !DefaultLogger.enableAPI {
		return
	}

	method := c.Request.Method
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 获取用户信息（如果存在）
	userInfo := "anonymous"
	if claims, exists := c.Get("claims"); exists {
		if claimsData, ok := claims.(map[string]interface{}); ok {
			if username, ok := claimsData["username"].(string); ok {
				userInfo = username
			}
		}
	}

	logMessage := fmt.Sprintf("API Request - Method: %s, Path: %s, Status: %d, Duration: %v, IP: %s, User: %s, UserAgent: %s",
		method, path, statusCode, duration, clientIP, userInfo, userAgent)

	if errorMsg != "" {
		logMessage += fmt.Sprintf(", Error: %s", errorMsg)
	}

	if statusCode >= 400 {
		Error(logMessage)
	} else {
		Info(logMessage)
	}
}

// LogUserAction 记录用户操作日志
func LogUserAction(username, action, resource, details string) {
	if DefaultLogger == nil {
		return
	}

	message := fmt.Sprintf("User Action - User: %s, Action: %s, Resource: %s, Details: %s",
		username, action, resource, details)
	Info(message)
}

// LogSystemEvent 记录系统事件日志
func LogSystemEvent(event, details string) {
	if DefaultLogger == nil {
		return
	}

	message := fmt.Sprintf("System Event - Event: %s, Details: %s", event, details)
	Info(message)
}

// LogSecurityEvent 记录安全事件日志
func LogSecurityEvent(event, username, ip, details string) {
	if DefaultLogger == nil {
		return
	}

	message := fmt.Sprintf("Security Event - Event: %s, User: %s, IP: %s, Details: %s",
		event, username, ip, details)
	Warn(message)
}

// GetLogLevel 获取当前日志级别
func GetLogLevel() LogLevel {
	if DefaultLogger != nil {
		return DefaultLogger.level
	}
	return INFO
}

// SetLogLevel 设置日志级别
func SetLogLevel(level LogLevel) {
	if DefaultLogger != nil {
		DefaultLogger.level = level
	}
}

// IsAPILoggingEnabled 检查是否启用API日志记录
func IsAPILoggingEnabled() bool {
	return DefaultLogger != nil && DefaultLogger.enableAPI
}

// SetAPILogging 设置API日志记录开关
func SetAPILogging(enabled bool) {
	if DefaultLogger != nil {
		DefaultLogger.enableAPI = enabled
	}
}
