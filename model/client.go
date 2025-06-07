package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClientLog records client connection status and usage
type ClientLog struct {
	ID                 string    `gorm:"primaryKey;size:36"`
	UserID             string    `gorm:"size:36;index"` // Indexed for faster lookups
	IsOnline           bool
	OnlineDuration     int64 // in seconds
	TrafficUsage       int64 // in bytes
	LastConnectionTime *time.Time
	CreatedAt          time.Time
}

// BeforeCreate will set a UUID rather than numeric ID.
func (cl *ClientLog) BeforeCreate(tx *gorm.DB) (err error) {
	cl.ID = uuid.NewString()
	return
}
