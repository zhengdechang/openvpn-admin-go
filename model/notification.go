package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType defines the type of VPN event
type NotificationType string

const (
	NotificationTypeConnected    NotificationType = "user_connected"
	NotificationTypeDisconnected NotificationType = "user_disconnected"
)

// Notification records a VPN connection event for superadmin review
type Notification struct {
	ID         string           `gorm:"primaryKey;size:36"`
	Type       NotificationType `gorm:"size:50;not null"`
	UserName   string           `gorm:"size:100;not null;index"`
	RealIP     string           `gorm:"size:45"`
	VirtualIP  string           `gorm:"size:45"`
	IsRead     bool             `gorm:"default:false;index"`
	CreatedAt  time.Time        `gorm:"index"`
}

// BeforeCreate sets a UUID primary key
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID = uuid.NewString()
	return
}
