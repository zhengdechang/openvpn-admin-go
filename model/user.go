package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role 定义用户角色
type Role string

const (
	RoleSuperAdmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleManager    Role = "manager"
	RoleUser       Role = "user"
)

// User 用户模型
type User struct {
	ID                 string `gorm:"primaryKey;size:36"`
	Name               string `gorm:"size:100;not null"`
	Email              string `gorm:"uniqueIndex;size:100;not null"`
	PasswordHash       string `gorm:"size:255;not null"`
	Role               Role   `gorm:"size:20;not null;default:'user'"`
	DepartmentID       string `gorm:"size:36"`
	LastConnectionTime *time.Time
	IsOnline           bool
	CreatorID          string `gorm:"size:36"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// BeforeCreate 在创建记录前生成 UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}
