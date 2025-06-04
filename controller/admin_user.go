package controller

import (
   "net/http"
   "github.com/gin-gonic/gin"
   "openvpn-admin-go/common"
   "openvpn-admin-go/database"
   "openvpn-admin-go/middleware"
   "openvpn-admin-go/model"
   "openvpn-admin-go/openvpn"
)

// AdminUserController 管理用户
type AdminUserController struct{}

// CreateUser 创建用户 (superadmin/admin 全量，manager 仅限本部门、仅 user 角色)
func (c *AdminUserController) CreateUser(ctx *gin.Context) {
   claims := ctx.MustGet("claims").(*middleware.Claims)
   var req struct {
       Name         string `json:"name" binding:"required"`
       Email        string `json:"email" binding:"required,email"`
       Password     string `json:"password" binding:"required,min=6"`
       Role         string `json:"role" binding:"required,oneof=superadmin admin manager user"`
       DepartmentID string `json:"departmentId"`
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
   }
   if err := database.DB.Create(&user).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   // 同步创建 OpenVPN 客户端配置
   if err := openvpn.CreateClient(user.ID); err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建 OpenVPN 客户端失败: " + err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"data": gin.H{
       "id":           user.ID,
       "name":         user.Name,
       "email":        user.Email,
       "role":         user.Role,
       "departmentId": user.DepartmentID,
   }})
}

// ListUsers 列出用户列表 (manager 仅本部门)
func (c *AdminUserController) ListUsers(ctx *gin.Context) {
   claims := ctx.MustGet("claims").(*middleware.Claims)
   var users []model.User
   db := database.DB
   if claims.Role == string(model.RoleManager) {
       db = db.Where("department_id = ?", claims.DeptID)
   }
   if err := db.Find(&users).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   var resp []gin.H
   for _, u := range users {
       resp = append(resp, gin.H{
           "id":           u.ID,
           "name":         u.Name,
           "email":        u.Email,
           "role":         u.Role,
           "departmentId": u.DepartmentID,
       })
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
   ctx.JSON(http.StatusOK, gin.H{"data": gin.H{
       "id":           u.ID,
       "name":         u.Name,
       "email":        u.Email,
       "role":         u.Role,
       "departmentId": u.DepartmentID,
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
   if claims.Role == string(model.RoleManager) && u.DepartmentID != claims.DeptID {
       ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only update own department users"})
       return
   }
   var req struct {
       Name         *string `json:"name"`
       Email        *string `json:"email" binding:"omitempty,email"`
       Role         *string `json:"role" binding:"omitempty,oneof=superadmin admin manager user"`
       DepartmentID *string `json:"departmentId"`
       Password     *string `json:"password" binding:"omitempty,min=6"`
   }
   if err := ctx.ShouldBindJSON(&req); err != nil {
       ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
       return
   }
   updates := map[string]interface{}{}
   if req.Name != nil {
       updates["name"] = *req.Name
   }
   if req.Email != nil {
       updates["email"] = *req.Email
   }
   if req.Role != nil {
       if claims.Role == string(model.RoleManager) {
           ctx.JSON(http.StatusForbidden, gin.H{"error": "manager cannot update role"})
           return
       }
       updates["role"] = *req.Role
   }
   if req.DepartmentID != nil {
       if claims.Role == string(model.RoleManager) && *req.DepartmentID != claims.DeptID {
           ctx.JSON(http.StatusForbidden, gin.H{"error": "manager cannot change department"})
           return
       }
       updates["department_id"] = *req.DepartmentID
   }
   if req.Password != nil {
       hash, err := common.HashPassword(*req.Password)
       if err != nil {
           ctx.JSON(http.StatusInternalServerError, gin.H{"error": "hash password failed"})
           return
       }
       updates["password_hash"] = hash
   }
   if len(updates) == 0 {
       ctx.JSON(http.StatusBadRequest, gin.H{"error": "no data to update"})
       return
   }
   if err := database.DB.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"message": "user updated"})
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
   if claims.Role == string(model.RoleManager) && u.DepartmentID != claims.DeptID {
       ctx.JSON(http.StatusForbidden, gin.H{"error": "manager can only delete own department users"})
       return
   }
   if err := database.DB.Delete(&model.User{}, "id = ?", id).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}