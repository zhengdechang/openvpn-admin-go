package utils

import (
	"os"
)

// GetEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 