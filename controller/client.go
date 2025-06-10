package controller

import (
	"net/http"
	"openvpn-admin-go/database" // Added
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"
	"openvpn-admin-go/openvpn"
	"time" // Added

	"github.com/gin-gonic/gin"
)

type ClientController struct{}

// UserClientStatusDetailed combines user data with OpenVPN client status
type UserClientStatusDetailed struct {
	// Fields from model.User (excluding PasswordHash)
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	Email              string     `json:"email"`
	Role               model.Role `json:"role"`
	DepartmentID       string     `json:"departmentId,omitempty"`
	CreatorID          string     `json:"creatorId,omitempty"`
	FixedIP            string     `json:"fixedIp,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DBIsOnline         bool       `json:"dbIsOnline"` // User.IsOnline from DB
	DBLastConnectionTime *time.Time `json:"dbLastConnectionTime,omitempty"`

	// Fields from openvpn.OpenVPNClientStatus
	ClientCommonName      string    `json:"clientCommonName"` // This is status.CommonName (user.Name)
	RealAddress           string    `json:"realAddress,omitempty"`
	VirtualAddress        string    `json:"virtualAddress,omitempty"`
	BytesReceived         int64     `json:"bytesReceived"`
	BytesSent             int64     `json:"bytesSent"`
	ConnectedSince        time.Time `json:"connectedSince,omitempty"`
	LastRef               time.Time `json:"lastRef,omitempty"`
	LiveIsOnline          bool      `json:"liveIsOnline"` // status.IsOnline from status log
	OnlineDurationSeconds int64     `json:"onlineDurationSeconds"`
}

// GetDetailedClientStatuses godoc
// @Summary Get detailed status for all OpenVPN clients including user data
// @Description Retrieves a list of all clients with live data and associated user details.
// @Tags Client
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} UserClientStatusDetailed "List of detailed client statuses"
// @Failure 500 {object} gin.H "Error message"
// @Router /client/status/detailed [get]
func (cc *ClientController) GetDetailedClientStatuses(c *gin.Context) {
	parsedStatuses, err := openvpn.ParseAllClientStatuses() // From openvpn/status_parser.go
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse client statuses: " + err.Error()})
		return
	}

	if parsedStatuses == nil {
		parsedStatuses = []openvpn.OpenVPNClientStatus{}
	}

	var detailedStatuses []UserClientStatusDetailed

	for _, status := range parsedStatuses {
		var user model.User
		// Query user by name, as status.CommonName is now user.Name
		dbErr := database.DB.Where("name = ?", status.CommonName).First(&user).Error

		item := UserClientStatusDetailed{
			ClientCommonName:      status.CommonName,
			RealAddress:           status.RealAddress,
			VirtualAddress:        status.VirtualAddress,
			BytesReceived:         status.BytesReceived,
			BytesSent:             status.BytesSent,
			ConnectedSince:        status.ConnectedSince,
			LastRef:               status.LastRef,
			LiveIsOnline:          status.IsOnline,
			OnlineDurationSeconds: status.OnlineDurationSeconds,
		}

		if dbErr == nil { // User found
			item.ID = user.ID
			item.Name = user.Name // Should be same as status.CommonName if matched
			item.Email = user.Email
			item.Role = user.Role
			item.DepartmentID = user.DepartmentID
			item.CreatorID = user.CreatorID
			item.FixedIP = user.FixedIP
			item.CreatedAt = user.CreatedAt
			item.UpdatedAt = user.UpdatedAt
			item.DBIsOnline = user.IsOnline
			item.DBLastConnectionTime = user.LastConnectionTime
		}
		detailedStatuses = append(detailedStatuses, item)
	}
	c.JSON(http.StatusOK, detailedStatuses)
}

// GetClientList 获取客户端列表
func (c *ClientController) GetClientList(ctx *gin.Context) {
	statuses, err := openvpn.GetAllClientStatuses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, statuses)
}

// GetLiveConnections godoc
// @Summary Get all currently connected OpenVPN clients' live status
// @Description Retrieves a list of all clients currently connected to the OpenVPN server with live data like IP, duration, data transfer.
// @Tags Client
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} openvpn.OpenVPNClientStatus "List of live client statuses"
// @Failure 500 {object} gin.H "Error message"
// @Router /client/status/live [get]
func (c *ClientController) GetLiveConnections(ctx *gin.Context) {
	statuses, err := openvpn.GetAllClientStatuses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get client statuses: " + err.Error()})
		return
	}
	if statuses == nil { // Handle case where parsing might return nil, nil
		statuses = []openvpn.ClientStatus{}
	}
	ctx.JSON(http.StatusOK, statuses)
}

// AddClient 添加客户端
func (c *ClientController) AddClient(ctx *gin.Context) {
	var client struct {
		Username string `json:"username" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := openvpn.CreateClient(client.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Client added successfully"})
}

// UpdateClient 更新客户端
func (c *ClientController) UpdateClient(ctx *gin.Context) {
	var client struct {
		Username string `json:"username" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新客户端实际上就是重新生成证书和配置
	if err := openvpn.CreateClient(client.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Client updated successfully"})
}

// DeleteClient 删除客户端
func (c *ClientController) DeleteClient(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if err := openvpn.DeleteClient(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}

// GetClientConfig 获取客户端配置
func (c *ClientController) GetClientConfig(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	// 加载配置
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 权限检查: 普通用户仅能下载自己的配置
	claims := ctx.MustGet("claims").(*middleware.Claims)
	if claims.Role == string(model.RoleUser) && claims.UserID != username {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	// 生成客户端配置
	config, err := openvpn.GenerateClientConfig(username, cfg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"config": config})
}

// RevokeClient 吊销客户端证书
func (c *ClientController) RevokeClient(ctx *gin.Context) {
	var client struct {
		Username string `json:"username" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := openvpn.DeleteClient(client.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Client certificate revoked successfully"})
}

// RenewClient 续期客户端证书
func (c *ClientController) RenewClient(ctx *gin.Context) {
	var client struct {
		Username string `json:"username" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 续期证书实际上就是重新生成证书和配置
	if err := openvpn.CreateClient(client.Username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Client certificate renewed successfully"})
}

// GetClientStatus 获取客户端状态
func (c *ClientController) GetClientStatus(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	claims := ctx.MustGet("claims").(*middleware.Claims)
	if claims.Role == string(model.RoleUser) && claims.UserID != username {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	status, err := openvpn.GetClientStatus(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if status == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	ctx.JSON(http.StatusOK, status)
}

// GetAllClientStatuses 获取所有客户端状态
func (c *ClientController) GetAllClientStatuses(ctx *gin.Context) {
	// 只有管理员和部门负责人可以查看所有状态
	// 路由已限制，此处无需重复检查
	statuses, err := openvpn.GetAllClientStatuses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get client statuses: " + err.Error()})
		return
	}
	if statuses == nil {
		// Ensure an empty array is returned instead of null if no statuses are found
		ctx.JSON(http.StatusOK, []openvpn.OpenVPNClientStatus{})
		return
	}
	ctx.JSON(http.StatusOK, statuses)
}

// PauseClient 暂停客户端连接
func (c *ClientController) PauseClient(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if err := openvpn.PauseClient(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Client paused successfully"})
}

// ResumeClient 恢复客户端连接
func (c *ClientController) ResumeClient(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if err := openvpn.ResumeClient(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Client resumed successfully"})
}
