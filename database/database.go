package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// buildDSN 根据环境变量构造 MySQL DSN。
// 优先使用 DB_DSN 整体覆盖，否则由离散变量拼装。
func buildDSN() string {
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		return dsn
	}

	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "openvpn")
	password := os.Getenv("DB_PASSWORD")
	dbName := getEnv("DB_NAME", "openvpn")

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Init 初始化数据库连接（MySQL）。
// 容器编排下 MySQL 可能稍晚就绪，这里带有限次数的重试。
func Init() error {
	dsn := buildDSN()

	const maxAttempts = 10
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = db
			return nil
		}
		lastErr = err
		fmt.Printf("等待数据库就绪 (%d/%d): %v\n", attempt, maxAttempts, err)
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("failed to connect database after %d attempts: %v", maxAttempts, lastErr)
}
