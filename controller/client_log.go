package controller

import (
	"net/http"
	"strconv"
	"time"

	"openvpn-admin-go/database"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// CreateClientLogRequest defines the request body for creating a client log.
type CreateClientLogRequest struct {
	UserID             string     `json:"user_id" binding:"required"`
	IsOnline           bool       `json:"is_online"`
	OnlineDuration     int64      `json:"online_duration"` // in seconds
	TrafficUsage       int64      `json:"traffic_usage"`   // in bytes
	LastConnectionTime *time.Time `json:"last_connection_time"`
}

// CreateClientLog handles the creation of a new client log entry.
func CreateClientLog(c *gin.Context) {
	var req CreateClientLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	clientLog := model.ClientLog{
		UserID:             req.UserID,
		IsOnline:           req.IsOnline,
		OnlineDuration:     req.OnlineDuration,
		TrafficUsage:       req.TrafficUsage,
		LastConnectionTime: req.LastConnectionTime,
		// CreatedAt will be set by GORM by default if not specified, or by BeforeCreate hook if needed for other reasons
	}

	if err := database.DB.Create(&clientLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create client log: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": clientLog})
}

// GetClientLogs retrieves paginated client log entries.
// Accepts query parameters: page (int, default 1), pageSize (int, default 10), user_id (string, optional).
func GetClientLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	userID := c.Query("user_id")

	offset := (page - 1) * pageSize
	var clientLogs []model.ClientLog
	var total int64

	query := database.DB.Model(&model.ClientLog{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to count client logs: " + err.Error()})
		return
	}

	// Retrieve paginated logs
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&clientLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to retrieve client logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"data":     clientLogs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
