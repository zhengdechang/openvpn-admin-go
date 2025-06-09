package controller

import (
	"net/http"
	"os" // Added
	"path/filepath" // Added
	"strings"

	"openvpn-admin-go/common"
	"openvpn-admin-go/constants" // Added
	"openvpn-admin-go/database"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"
	"openvpn-admin-go/openvpn"

	"github.com/gin-gonic/gin"
)

// AdminUserController 管理用户
type AdminUserController struct{}

// CreateUser 创建用户 (superadmin/admin 全量，manager 仅限本部门、仅 user 角色)
func (c *AdminUserController) CreateUser(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	var req struct {
		Name         string  `json:"name" binding:"required"`
		Email        string  `json:"email" binding:"required,email"`
		Password     string  `json:"password" binding:"required,min=6"`
		Role         string  `json:"role" binding:"required,oneof=superadmin admin manager user"`
		DepartmentID string  `json:"departmentId"`
		FixedIP      *string `json:"fixedIp" binding:"omitempty,ip|cidrv4|cidrv6"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// manager 权限限制
	if claims.Role == string(model.RoleManager) {
		if req.DepartmentID != claims.DeptID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only create users in own department"})
			return
		}
		if req.Role != string(model.RoleUser) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only assign user role"})
			return
		}
		if req.FixedIP != nil && strings.TrimSpace(*req.FixedIP) != "" { // Check against trimmed value
			ctx.JSON(http.StatusForbidden, gin.H{"error": "manager cannot set fixed IP"})
			return
		}
	}

	hash, err := common.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "hash password failed"})
		return
	}
	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         model.Role(req.Role),
		DepartmentID: req.DepartmentID,
		CreatorID:    claims.UserID,
	}

	// Handle FixedIP assignment on creation
	if req.FixedIP != nil {
		trimmedFixedIP := strings.TrimSpace(*req.FixedIP)
		if trimmedFixedIP != "" {
			if !(claims.Role == string(model.RoleSuperAdmin) || claims.Role == string(model.RoleAdmin)) {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "only superadmin or admin can set fixed IP during creation"})
				return
			}
			user.FixedIP = trimmedFixedIP
		}
	}

	if err := database.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user in database: " + err.Error()})
		return
	}

	// Check for existing client configuration
	clientOvpnFile := filepath.Join(constants.ClientConfigDir, user.Name+".ovpn")
	if _, err := os.Stat(clientOvpnFile); err == nil {
		// File exists, so client with this name likely exists
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "A VPN client with the username '" + user.Name + "' already exists."})
		return
	} else if !os.IsNotExist(err) {
		// Another error occurred during stat, not just file not existing
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing VPN client config: " + err.Error()})
		return
	}

	// If fixed IP was set for the user, apply it via CCD
	if user.FixedIP != "" {
		if err := openvpn.SetClientFixedIP(user.Name, user.FixedIP); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user created, but failed to set fixed IP in OpenVPN config: " + err.Error()})
			return
		}
	}

	// Create OpenVPN client certs
	if err := openvpn.CreateClient(user.Name); err != nil {
		// If cert creation fails, and fixed IP was set, should we try to remove CCD?
		if user.FixedIP != "" {
			openvpn.RemoveClientFixedIP(user.Name) // Attempt to clean up CCD
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create OpenVPN client certificate: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":           user.ID,
		"name":         user.Name,
		"email":        user.Email,
		"role":         user.Role,
		"departmentId": user.DepartmentID,
		"fixedIp":      user.FixedIP,
	}})
}

// ListUsers 列出用户列表 (manager 仅本部门)
func (c *AdminUserController) ListUsers(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	var users []model.User
	db := database.DB
	// 部门负责人仅查看本部门用户；普通用户仅查看自身
	if claims.Role == string(model.RoleManager) {
		db = db.Where("department_id = ?", claims.DeptID)
	} else if claims.Role == string(model.RoleUser) {
		db = db.Where("id = ?", claims.UserID)
	}

	// Added Order by created_at desc
	if err := db.Order("created_at desc").Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users: " + err.Error()})
		return
	}

	// 获取在线状态
	statuses, err := openvpn.GetAllClientStatuses() // 使用 GetAllClientStatuses 而不是 ParseAllClientStatuses
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取客户端状态失败"})
		return
	}

	var resp []gin.H
	for _, u := range users {
		userData := gin.H{
			"id":                 u.ID,
			"name":               u.Name,
			"email":              u.Email,
			"role":               u.Role,
			"departmentId":       u.DepartmentID,
			"creatorId":          u.CreatorID,
			"lastConnectionTime": u.LastConnectionTime, // This is from DB, might be historical
			"fixedIp":            u.FixedIP,            // From DB
			"createdAt":          u.CreatedAt,
			"updatedAt":          u.UpdatedAt,
			"isOnline":           false,
			"connectionIp":       nil,
			"allocatedVpnIp":     nil,
		}

		// 检查用户是否在线
		for _, status := range statuses {
			if status.CommonName == u.ID {
				userData["isOnline"] = true
				userData["connectionIp"] = status.RealAddress
				userData["allocatedVpnIp"] = status.VirtualAddress
				break
			}
		}

		resp = append(resp, userData)
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetUser 获取单个用户 (manager 仅本部门)
func (c *AdminUserController) GetUser(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	id := ctx.Param("id")
	var u model.User
	if err := database.DB.First(&u, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if claims.Role == string(model.RoleManager) && u.DepartmentID != claims.DeptID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only view own department users"})
		return
	}
	// Stricter RBAC for GetUser based on the comprehensive plan
	if claims.Role == string(model.RoleUser) && u.ID != claims.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "user can only view self"})
		return
	}
	// Manager can view self even if not in their department (e.g. if they are head of a parent dept)
	if claims.Role == string(model.RoleManager) && u.DepartmentID != claims.DeptID && u.ID != claims.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only view own department users or self"})
		return
	}

	isOnline := false
	var connectionIp interface{} = nil
	var allocatedVpnIp interface{} = nil

	liveStatus, err := openvpn.GetClientStatus(u.ID)
	if err == nil && liveStatus != nil {
		isOnline = true
		connectionIp = liveStatus.RealAddress
		allocatedVpnIp = liveStatus.VirtualAddress
	} else if err != nil {
		// Log error getting individual client status, but don't fail the request
		// log.Printf("Warning: Failed to get live status for user %s: %v", u.ID, err)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":                 u.ID,
		"name":               u.Name,
		"email":              u.Email,
		"role":               u.Role,
		"departmentId":       u.DepartmentID,
		"fixedIp":            u.FixedIP,
		"isOnline":           isOnline,
		"connectionIp":       connectionIp,
		"allocatedVpnIp":     allocatedVpnIp,
		"lastConnectionTime": u.LastConnectionTime,
		"creatorId":          u.CreatorID,
		"createdAt":          u.CreatedAt,
		"updatedAt":          u.UpdatedAt,
	}})
}

// UpdateUser 更新用户 (manager 对自身部门用户权限受限)
func (c *AdminUserController) UpdateUser(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	id := ctx.Param("id")
	var u model.User
	if err := database.DB.First(&u, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// RBAC: Managers can only update users in their own department.
	// Users cannot update other users, only their own profile (but this endpoint is admin-focused).
	if claims.Role == string(model.RoleManager) && u.DepartmentID != claims.DeptID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only update users in own department"})
		return
	}
	if claims.Role == string(model.RoleUser) && u.ID != claims.UserID { // Should not happen if routes are correct
		ctx.JSON(http.StatusForbidden, gin.H{"error": "user can only update their own profile via dedicated endpoint"})
		return
	}

	var req struct {
		Name         *string `json:"name"`
		Email        *string `json:"email" binding:"omitempty,email"`
		Role         *string `json:"role" binding:"omitempty,oneof=superadmin admin manager user"`
		DepartmentID *string `json:"departmentId"`
		FixedIP      *string `json:"fixedIp" binding:"omitempty,ip|cidrv4|cidrv6"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if claims.Role == string(model.RoleManager) {
		if req.DepartmentID != nil && *req.DepartmentID != claims.DeptID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only update users in own department"})
			return
		}
		if req.Role != nil && *req.Role != string(model.RoleUser) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only assign user role"})
			return
		}
		if req.FixedIP != nil && strings.TrimSpace(*req.FixedIP) != "" { // Check against trimmed value
			ctx.JSON(http.StatusForbidden, gin.H{"error": "manager cannot set fixed IP"})
			return
		}
	}

	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Role != nil {
		u.Role = model.Role(*req.Role)
	}
	if req.DepartmentID != nil {
		u.DepartmentID = *req.DepartmentID
	}
	if req.FixedIP != nil {
		trimmedFixedIP := strings.TrimSpace(*req.FixedIP)
		if trimmedFixedIP != "" {
			if !(claims.Role == string(model.RoleSuperAdmin) || claims.Role == string(model.RoleAdmin)) {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "only superadmin or admin can set fixed IP during update"})
				return
			}
			u.FixedIP = trimmedFixedIP
		}
	}

	if err := database.DB.Save(&u).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user in database: " + err.Error()})
		return
	}

	// If fixed IP was set for the user, apply it via CCD
	if u.FixedIP != "" {
		if err := openvpn.SetClientFixedIP(u.Name, u.FixedIP); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user updated, but failed to set fixed IP in OpenVPN config: " + err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":           u.ID,
		"name":         u.Name,
		"email":        u.Email,
		"role":         u.Role,
		"departmentId": u.DepartmentID,
		"fixedIp":      u.FixedIP,
		"updatedAt":    u.UpdatedAt,
	}})
}

// DeleteUser 删除用户 (manager 仅本部门)
func (c *AdminUserController) DeleteUser(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	id := ctx.Param("id")
	var u model.User
	if err := database.DB.First(&u, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// RBAC for deletion
	if claims.Role == string(model.RoleManager) && (u.DepartmentID != claims.DeptID || u.ID == claims.UserID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only delete users in own department and cannot delete self"})
		return
	}
	if u.Role == model.RoleSuperAdmin && claims.Role != string(model.RoleSuperAdmin) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can delete superadmin user"})
		return
	}
	if u.ID == claims.UserID && claims.Role != string(model.RoleSuperAdmin) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you cannot delete your own account unless you are a superadmin"})
		return
	}

	// Remove Fixed IP config if it exists
	if u.FixedIP != "" {
		if err := openvpn.RemoveClientFixedIP(u.Name); err != nil {
			// log.Printf("Warning: failed to remove fixed IP for user %s during deletion: %v", u.Name, err)
		}
	}
	// Remove OpenVPN client certificate and other related configs
	if err := openvpn.DeleteClient(u.Name); err != nil {
		// log.Printf("Warning: failed to delete OpenVPN client data for user %s during deletion: %v", u.Name, err)
	}

	if err := database.DB.Delete(&model.User{}, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user from database: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
