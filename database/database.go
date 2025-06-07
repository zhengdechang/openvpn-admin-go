package database

import (
   "fmt"
   "os"
   "path/filepath"

   "gorm.io/driver/sqlite"
   "gorm.io/gorm"
   "openvpn-admin-go/model" // Import the model package
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
   dbPath := os.Getenv("DB_PATH")
   if dbPath == "" {
       dbPath = "data/db.sqlite3"
   }
   dir := filepath.Dir(dbPath)
   if err := os.MkdirAll(dir, 0755); err != nil {
       return fmt.Errorf("failed to create database directory: %v", err)
   }
   db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
   if err != nil {
       return fmt.Errorf("failed to connect database: %v", err)
   }
   DB = db
   return nil
}

// Migrate 自动迁移模型
func Migrate(models ...interface{}) error {
   // Ensure User and ClientLog are always included in migrations
   allModels := []interface{}{&model.User{}, &model.ClientLog{}}
   allModels = append(allModels, models...)
   return DB.AutoMigrate(allModels...)
}