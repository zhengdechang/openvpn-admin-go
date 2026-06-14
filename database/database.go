package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

var DB *gorm.DB

// Init initializes the database connection.
func Init() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getEnvOrDefault("DB_HOST", "localhost")
		user := getEnvOrDefault("DB_USER", "openvpn")
		password := getEnvOrDefault("DB_PASSWORD", "openvpn")
		dbname := getEnvOrDefault("DB_NAME", "openvpn")
		port := getEnvOrDefault("DB_PORT", "5432")
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}
	DB = db
	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// Migrate runs all pending Goose migrations using the embedded SQL files.
// The variadic parameter is kept for backwards-compat with existing callers.
func Migrate(_ ...interface{}) error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %v", err)
	}

	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %v", err)
	}

	if err := runMigrations(sqlDB); err != nil {
		return fmt.Errorf("goose migrate: %v", err)
	}
	return nil
}

func runMigrations(db *sql.DB) error {
	return goose.Up(db, "migrations")
}

// gooseDB 配置 goose（嵌入式 SQL + postgres 方言）并返回底层 *sql.DB。
func gooseDB() (*sql.DB, error) {
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %v", err)
	}
	goose.SetBaseFS(migrationFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("goose set dialect: %v", err)
	}
	return sqlDB, nil
}

// RunMigrations 应用所有未执行的迁移（goose up），供 CLI `migrate up` 使用。
func RunMigrations() error {
	sqlDB, err := gooseDB()
	if err != nil {
		return err
	}
	return goose.Up(sqlDB, "migrations")
}

// MigrateDown 回滚最近一次迁移（goose down），供 CLI `migrate down` 使用。
func MigrateDown() error {
	sqlDB, err := gooseDB()
	if err != nil {
		return err
	}
	return goose.Down(sqlDB, "migrations")
}

// MigrateStatus 打印迁移状态（goose status），供 CLI `migrate status` 使用。
func MigrateStatus() error {
	sqlDB, err := gooseDB()
	if err != nil {
		return err
	}
	return goose.Status(sqlDB, "migrations")
}
