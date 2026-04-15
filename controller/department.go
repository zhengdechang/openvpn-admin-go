package controller

import (
	"openvpn-admin-go/common"
	"openvpn-admin-go/database"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DepartmentController 管理部门
type DepartmentController struct{}

// CreateDepartment 创建部门（事务：建部门 + 关联负责人）
func (c *DepartmentController) CreateDepartment(ctx *gin.Context) {
	var dep model.Department
	if err := ctx.ShouldBindJSON(&dep); err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dep).Error; err != nil {
			return err
		}
		if dep.HeadID != "" {
			return tx.Model(&model.User{}).
				Where("id = ?", dep.HeadID).
				Update("department_id", dep.ID).Error
		}
		return nil
	}); err != nil {
		common.InternalError(ctx, err.Error())
		return
	}

	common.OK(ctx, dep)
}

// ListDepartments 列出所有部门
func (c *DepartmentController) ListDepartments(ctx *gin.Context) {
	var deps []model.Department
	if err := database.DB.
		Preload("Head").
		Preload("Parent").
		Preload("Children").
		Find(&deps).Error; err != nil {
		common.InternalError(ctx, err.Error())
		return
	}
	common.OK(ctx, deps)
}

// GetDepartment 获取部门详情
func (c *DepartmentController) GetDepartment(ctx *gin.Context) {
	id := ctx.Param("id")
	var dep model.Department
	if err := database.DB.
		Preload("Head").
		Preload("Parent").
		Preload("Children").
		First(&dep, "id = ?", id).Error; err != nil {
		common.NotFound(ctx, "department not found")
		return
	}
	common.OK(ctx, dep)
}

// UpdateDepartment 更新部门（事务：更新部门 + 负责人变更）
func (c *DepartmentController) UpdateDepartment(ctx *gin.Context) {
	id := ctx.Param("id")
	var existing model.Department
	if err := database.DB.First(&existing, "id = ?", id).Error; err != nil {
		common.NotFound(ctx, "department not found")
		return
	}
	var req model.Department
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{"name": req.Name, "head_id": req.HeadID, "parent_id": req.ParentID}
		if err := tx.Model(&model.Department{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if req.HeadID != existing.HeadID {
			if existing.HeadID != "" {
				if err := tx.Model(&model.User{}).
					Where("id = ?", existing.HeadID).
					Update("department_id", "").Error; err != nil {
					return err
				}
			}
			if req.HeadID != "" {
				if err := tx.Model(&model.User{}).
					Where("id = ?", req.HeadID).
					Update("department_id", id).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		common.InternalError(ctx, err.Error())
		return
	}

	common.OKMsg(ctx, "department updated")
}

// DeleteDepartment 删除部门
func (c *DepartmentController) DeleteDepartment(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := database.DB.Delete(&model.Department{}, "id = ?", id).Error; err != nil {
		common.InternalError(ctx, err.Error())
		return
	}
	common.OKMsg(ctx, "department deleted")
}
