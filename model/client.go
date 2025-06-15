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
   ID                 string    `gorm:"primaryKey;size:36"`
   Name               string    `gorm:"uniqueIndex;size:100;not null"`
   Email              string    `gorm:"uniqueIndex;size:100;not null"`
   PasswordHash       string    `gorm:"size:255;not null"`
   Role               Role      `gorm:"size:20;not null"`
   DepartmentID       string    `gorm:"size:36"`
   CreatorID          string    `gorm:"size:36"`
   FixedIP            string    `gorm:"size:45"` // IPv4 or IPv6 address
   Subnet             string    `gorm:"size:45"` // Subnet in CIDR format (e.g., 10.10.120.0/23)
   CreatedAt          time.Time
   UpdatedAt          time.Time

   // OpenVPN status fields
   IsOnline           bool      `gorm:"default:false"`
   LastConnectionTime *time.Time
   RealAddress        string    `gorm:"size:45"` // Client's real IP address
   VirtualAddress     string    `gorm:"size:45"` // Client's VPN IP address
   BytesReceived      int64     `gorm:"default:0"`
   BytesSent          int64     `gorm:"default:0"`
   ConnectedSince     *time.Time
   LastRef            *time.Time
   OnlineDuration     int64     `gorm:"default:0"` // Duration in seconds
   IsPaused           bool      `gorm:"default:false"`
}

// BeforeCreate 在创建记录前生成 UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
   u.ID = uuid.NewString()
   return
}