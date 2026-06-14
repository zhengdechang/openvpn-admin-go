package controller

import (
	"net/http"
	"time"

	"openvpn-admin-go/database"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// NotificationController handles notification endpoints
type NotificationController struct{}

// notificationResponse is the JSON shape returned to the frontend
type notificationResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	UserName  string `json:"userName"`
	RealIP    string `json:"realIP"`
	VirtualIP string `json:"virtualIP"`
	IsRead    bool   `json:"isRead"`
	CreatedAt string `json:"createdAt"`
}

func toResponse(n model.Notification) notificationResponse {
	return notificationResponse{
		ID:        n.ID,
		Type:      string(n.Type),
		UserName:  n.UserName,
		RealIP:    n.RealIP,
		VirtualIP: n.VirtualIP,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// List returns the 50 most recent notifications, newest first
func (nc *NotificationController) List(c *gin.Context) {
	var notifications []model.Notification
	if err := database.DB.
		Order("created_at DESC").
		Limit(50).
		Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to fetch notifications"})
		return
	}

	resp := make([]notificationResponse, 0, len(notifications))
	for _, n := range notifications {
		resp = append(resp, toResponse(n))
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// UnreadCount returns the count of unread notifications
func (nc *NotificationController) UnreadCount(c *gin.Context) {
	var count int64
	if err := database.DB.Model(&model.Notification{}).Where("is_read = ?", false).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to count notifications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"count": count}})
}

// MarkRead marks a single notification as read by ID
func (nc *NotificationController) MarkRead(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to mark notification as read"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "notification not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MarkAllRead marks all notifications as read
func (nc *NotificationController) MarkAllRead(c *gin.Context) {
	if err := database.DB.Model(&model.Notification{}).Where("is_read = ?", false).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to mark all notifications as read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
