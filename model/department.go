package model

import (
   "time"

   "github.com/google/uuid"
   "gorm.io/gorm"
)

// Department 部门模型
type Department struct {
   ID        string    `gorm:"primaryKey;size:36"`
   Name      string    `gorm:"size:100;uniqueIndex;not null"`
   CreatedAt time.Time
   UpdatedAt time.Time
   Users     []User    `gorm:"foreignKey:DepartmentID"`
}

// BeforeCreate 在创建记录前生成 UUID
func (d *Department) BeforeCreate(tx *gorm.DB) (err error) {
   d.ID = uuid.NewString()
   return
}