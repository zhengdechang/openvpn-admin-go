package controller

import (
   "net/http"

   "github.com/gin-gonic/gin"
   "openvpn-admin-go/database"
   "openvpn-admin-go/model"
)

// DepartmentController 管理部门
type DepartmentController struct{}

// CreateDepartment 创建部门
func (c *DepartmentController) CreateDepartment(ctx *gin.Context) {
   var dep model.Department
   if err := ctx.ShouldBindJSON(&dep); err != nil {
       ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
       return
   }
   if err := database.DB.Create(&dep).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, dep)
}

// ListDepartments 列出所有部门
func (c *DepartmentController) ListDepartments(ctx *gin.Context) {
   var deps []model.Department
   // 预加载负责人信息
   // 预加载负责人、上级部门及子部门信息
   if err := database.DB.
       Preload("Head").
       Preload("Parent").
       Preload("Children").
       Find(&deps).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, deps)
}

// GetDepartment 获取部门详情
func (c *DepartmentController) GetDepartment(ctx *gin.Context) {
   id := ctx.Param("id")
   var dep model.Department
   // 加载部门及负责人、上级和子部门信息
   if err := database.DB.
       Preload("Head").
       Preload("Parent").
       Preload("Children").
       First(&dep, "id = ?", id).Error; err != nil {
       ctx.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
       return
   }
   ctx.JSON(http.StatusOK, dep)
}

// UpdateDepartment 更新部门
func (c *DepartmentController) UpdateDepartment(ctx *gin.Context) {
   id := ctx.Param("id")
   var req model.Department
   if err := ctx.ShouldBindJSON(&req); err != nil {
       ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
       return
   }
   // 更新名称和负责人
   updates := map[string]interface{}{"name": req.Name, "head_id": req.HeadID, "parent_id": req.ParentID}
   if err := database.DB.Model(&model.Department{}).
       Where("id = ?", id).
       Updates(updates).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"message": "department updated"})
}

// DeleteDepartment 删除部门
func (c *DepartmentController) DeleteDepartment(ctx *gin.Context) {
   id := ctx.Param("id")
   if err := database.DB.Delete(&model.Department{}, "id = ?", id).Error; err != nil {
       ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }
   ctx.JSON(http.StatusOK, gin.H{"message": "department deleted"})
}