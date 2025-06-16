package controller

import (
   "net/http"

   "github.com/gin-gonic/gin"
   "openvpn-admin-go/database"
   "openvpn-admin-go/model"
   "math"    // For pagination
   "strconv" // For pagination
)

// DepartmentController 管理部门
type DepartmentController struct{}

// Helper function to recursively load children for a department
func loadChildren(db *gorm.DB, department *model.Department) error {
	if department == nil {
		return nil
	}
	var children []model.Department
	// Preload Head for children as well
	if err := db.Order("name asc").Preload("Head").Where("parent_id = ?", department.ID).Find(&children).Error; err != nil {
		return err
	}
	department.Children = children // Assign loaded children

	for i := range department.Children {
		if err := loadChildren(db, &department.Children[i]); err != nil {
			return err
		}
	}
	return nil
}

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
   // 将负责人用户关联到此部门
   if dep.HeadID != "" {
       if err := database.DB.Model(&model.User{}).
           Where("id = ?", dep.HeadID).
           Update("department_id", dep.ID).Error; err != nil {
           ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新负责人部门关联失败"})
           return
       }
   }
   ctx.JSON(http.StatusOK, dep)
}

// ListDepartments 列出所有部门 (paginated for top-level)
func (c *DepartmentController) ListDepartments(ctx *gin.Context) {
	// Pagination parameters
	pageQuery := ctx.DefaultQuery("page", "1")
	pageSizeQuery := ctx.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize < 1 || pageSize > 100 { // Max pageSize 100
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var topLevelDepartments []model.Department
	var totalItems int64

	db := database.DB

	// Count total top-level departments
	if err := db.Model(&model.Department{}).Where("parent_id IS NULL OR parent_id = ?", "").Count(&totalItems).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count top-level departments: " + err.Error()})
		return
	}

	// Fetch paginated top-level departments
	// Ordered by name for consistency
	if err := db.Order("name asc").Preload("Head").Where("parent_id IS NULL OR parent_id = ?", "").Offset(offset).Limit(pageSize).Find(&topLevelDepartments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list top-level departments: " + err.Error()})
		return
	}

	// For each top-level department, load its children recursively
	for i := range topLevelDepartments {
		if err := loadChildren(db, &topLevelDepartments[i]); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load department hierarchy: " + err.Error()})
			return
		}
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	if totalPages == 0 && totalItems > 0 {
		totalPages = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"totalItems":  totalItems,
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
		"departments": topLevelDepartments,
	})
}

// GetDepartment 获取部门详情
func (c *DepartmentController) GetDepartment(ctx *gin.Context) {
   id := ctx.Param("id")
   var dep model.Department
   // 加载部门及负责人、上级和子部门信息
   // For GetDepartment, we want the full hierarchy including parents, so existing preloads are fine.
   // However, to be consistent with how ListDepartments now loads children, we can use our loadChildren helper.
   if err := database.DB.Preload("Head").Preload("Parent").First(&dep, "id = ?", id).Error; err != nil {
       ctx.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
       return
   }
   // Manually load children for this specific department
   if err := loadChildren(database.DB, &dep); err != nil {
	   ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load department children: " + err.Error()})
	   return
   }
   ctx.JSON(http.StatusOK, dep)
}

// UpdateDepartment 更新部门
func (c *DepartmentController) UpdateDepartment(ctx *gin.Context) {
   id := ctx.Param("id")
   // 读取现有部门，用于判断负责人变更
   var existing model.Department
   if err := database.DB.First(&existing, "id = ?", id).Error; err != nil {
       ctx.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
       return
   }
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
   // 如果部门负责人发生变更，同步更新用户的部门关联
   if req.HeadID != existing.HeadID {
       // 清除旧负责人关联
       if existing.HeadID != "" {
           _ = database.DB.Model(&model.User{}).
               Where("id = ?", existing.HeadID).
               Update("department_id", "").Error
       }
       // 设置新负责人关联
       if req.HeadID != "" {
           if err := database.DB.Model(&model.User{}).
               Where("id = ?", req.HeadID).
               Update("department_id", id).Error; err != nil {
               ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新新负责人部门关联失败"})
               return
           }
       }
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