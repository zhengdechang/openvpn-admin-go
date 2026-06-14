package database

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const migrationsDir = "migrations"

// prepareGoose 配置 goose 使用内嵌的迁移文件与 MySQL 方言。
func prepareGoose() error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("设置 goose 方言失败: %v", err)
	}
	return nil
}

// RunMigrations 执行所有未应用的迁移（goose up）。
func RunMigrations() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}
	if err := prepareGoose(); err != nil {
		return err
	}
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("执行数据库迁移失败: %v", err)
	}
	return nil
}

// MigrateDown 回滚最近一次迁移（goose down）。
func MigrateDown() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}
	if err := prepareGoose(); err != nil {
		return err
	}
	return goose.Down(sqlDB, migrationsDir)
}

// MigrateStatus 打印迁移状态（goose status）。
func MigrateStatus() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}
	if err := prepareGoose(); err != nil {
		return err
	}
	return goose.Status(sqlDB, migrationsDir)
}
