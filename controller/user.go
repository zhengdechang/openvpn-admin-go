package controller

import (
	"net/http"

	"openvpn-admin-go/common"
	"openvpn-admin-go/database"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// 注册请求结构
type registerRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	hash, err := common.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "hash password failed"})
		return
	}
	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         model.RoleUser,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "register success"})
}

// 登录请求结构
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	var user model.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid credentials"})
		return
	}
	if !common.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid credentials"})
		return
	}
	token, err := middleware.GenerateToken(user.ID, string(user.Role), user.DepartmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "generate token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"user": gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}, "token": token}})
}

// VerifyEmail 邮箱验证 (暂未实现)
func VerifyEmail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "email verified"})
}

// ForgotPassword 忘记密码 (暂未实现)
func ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "reset link sent"})
}

// ResetPassword 重置密码 (暂未实现)
func ResetPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "password reset"})
}

// GetMe 获取当前用户信息
func GetMe(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	var user model.User
	if err := database.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}})
}

// UpdateMe 更新当前用户信息
type updateMeRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email" binding:"omitempty,email"`
}

func UpdateMe(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "no data to update"})
		return
	}
	if err := database.DB.Model(&model.User{}).Where("id = ?", claims.UserID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "update success"})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// RefreshToken 刷新 JWT
func RefreshToken(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	token, err := middleware.GenerateToken(claims.UserID, claims.Role, claims.DeptID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "generate token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"token": token}})
}

// GetRoles 获取角色列表
func GetRoles(c *gin.Context) {
	roles := []string{
		string(model.RoleSuperAdmin),
		string(model.RoleAdmin),
		string(model.RoleManager),
		string(model.RoleUser),
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": roles})
}

// GetUserInfo 查询指定用户
func GetUserInfo(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"dept":  user.DepartmentID,
	}})
}
