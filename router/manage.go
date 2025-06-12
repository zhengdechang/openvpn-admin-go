/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-04 10:37:43
 */
package router

import (
	"openvpn-admin-go/controller"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// SetupManageRoutes 设置部门和用户管理路由
func SetupManageRoutes(r *gin.RouterGroup) {
	depCtrl := &controller.DepartmentController{}
	// 部门管理: superadmin, admin
	dep := r.Group("/departments")
	dep.GET("", depCtrl.ListDepartments)
	dep.Use(middleware.JWTAuthMiddleware(), middleware.RoleRequired(
		string(model.RoleSuperAdmin), string(model.RoleAdmin)))
	{
		dep.POST("", depCtrl.CreateDepartment)
		dep.GET("/:id", depCtrl.GetDepartment)
		dep.PUT("/:id", depCtrl.UpdateDepartment)
		dep.DELETE("/:id", depCtrl.DeleteDepartment)
	}

	userCtrl := &controller.AdminUserController{}
	// 用户管理 routes have been moved to router/client.go
}
