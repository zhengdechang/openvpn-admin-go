package controller

import (
	"net/http"

	"openvpn-admin-go/common"
	"openvpn-admin-go/database"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 注册请求结构
type registerRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	DepartmentID    string `json:"departmentId" binding:"required"`
}

// Register 用户注册：必须选择部门；新用户默认待批准，需管理员审批后方可登录。
func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, err.Error())
		return
	}

	// 校验部门存在
	var dept model.Department
	if err := database.DB.First(&dept, "id = ?", req.DepartmentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid department"})
		return
	}

	hash, err := common.HashPassword(req.Password)
	if err != nil {
		common.InternalError(c, "hash password failed")
		return
	}
	user := model.User{
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   hash,
		Role:           model.RoleUser,
		DepartmentID:   req.DepartmentID,
		ApprovalStatus: model.ApprovalPending, // 待批准
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return tx.Model(&user).Update("CreatorID", user.ID).Error
	}); err != nil {
		common.InternalError(c, err.Error())
		return
	}

	// Set CreatorID to the user's own ID after creation
	if err := database.DB.Model(&user).Update("CreatorID", user.ID).Error; err != nil {
		// Log this error, but don't fail the registration because of it
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        "register success, pending approval",
		"approvalStatus": string(model.ApprovalPending),
	})
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
		common.BadRequest(c, err.Error())
		return
	}
	var user model.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		common.Unauthorized(c, "invalid credentials")
		return
	}
	if !common.CheckPasswordHash(req.Password, user.PasswordHash) {
		common.Unauthorized(c, "invalid credentials")
		return
	}
	// 审批门控：未批准的用户不能登录
	if user.ApprovalStatus != model.ApprovalApproved {
		msg := "account pending approval"
		if user.ApprovalStatus == model.ApprovalRejected {
			msg = "account registration rejected"
		}
		c.JSON(http.StatusForbidden, gin.H{
			"success":        false,
			"error":          msg,
			"approvalStatus": string(user.ApprovalStatus),
		})
		return
	}
	token, err := middleware.GenerateToken(user.ID, string(user.Role), user.DepartmentID)
	if err != nil {
		common.InternalError(c, "generate token failed")
		return
	}
	common.OK(c, gin.H{"user": gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}, "token": token})
}

// VerifyEmail 邮箱验证 (暂未实现)
func VerifyEmail(c *gin.Context) {
	common.OKMsg(c, "email verified")
}

// ForgotPassword 忘记密码 (暂未实现)
func ForgotPassword(c *gin.Context) {
	common.OKMsg(c, "reset link sent")
}

// ResetPassword 重置密码 (暂未实现)
func ResetPassword(c *gin.Context) {
	common.OKMsg(c, "password reset")
}

// GetMe 获取当前用户信息
func GetMe(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	var user model.User
	if err := database.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
		common.InternalError(c, "user not found")
		return
	}

	responseData := gin.H{
		"id":           user.ID,
		"name":         user.Name,
		"email":        user.Email,
		"role":         user.Role,
		"departmentId": user.DepartmentID,
		"isOnline":     user.IsOnline,
		"creatorId":    user.CreatorID,
	}

	if user.LastConnectionTime != nil {
		responseData["lastConnectionTime"] = *user.LastConnectionTime
	} else {
		responseData["lastConnectionTime"] = nil
	}

	common.OK(c, responseData)
}

// UpdateMe 更新当前用户信息
type updateMeRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
}

func UpdateMe(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, err.Error())
		return
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Password != nil {
		hashed, err := common.HashPassword(*req.Password)
		if err != nil {
			common.InternalError(c, "hash password failed")
			return
		}
		updates["password_hash"] = hashed
	}
	if len(updates) == 0 {
		common.BadRequest(c, "no data to update")
		return
	}
	if err := database.DB.Model(&model.User{}).Where("id = ?", claims.UserID).Updates(updates).Error; err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.OKMsg(c, "update success")
}

// Logout 用户登出
func Logout(c *gin.Context) {
	common.OK(c, nil)
}

// RefreshToken 刷新 JWT
func RefreshToken(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)
	token, err := middleware.GenerateToken(claims.UserID, claims.Role, claims.DeptID)
	if err != nil {
		common.InternalError(c, "generate token failed")
		return
	}
	common.OK(c, gin.H{"token": token})
}

// GetRoles 获取角色列表
func GetRoles(c *gin.Context) {
	roles := []string{
		string(model.RoleSuperAdmin),
		string(model.RoleAdmin),
		string(model.RoleManager),
		string(model.RoleUser),
	}
	common.OK(c, roles)
}

// GetUserInfo 查询指定用户
func GetUserInfo(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		common.NotFound(c, "user not found")
		return
	}
	responseData := gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"dept":      user.DepartmentID,
		"isOnline":  user.IsOnline,
		"creatorId": user.CreatorID,
	}
	if user.LastConnectionTime != nil {
		responseData["lastConnectionTime"] = *user.LastConnectionTime
	} else {
		responseData["lastConnectionTime"] = nil
	}
	common.OK(c, responseData)
}
