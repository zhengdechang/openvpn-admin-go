package router

import (
	"openvpn-admin-go/controller"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// SetupClientRoutes 设置客户端相关路由 (now includes user management)
func SetupClientRoutes(r *gin.RouterGroup) {
	clientCtrl := &controller.ClientController{}
	client := r.Group("/client")
	client.Use(middleware.JWTAuthMiddleware())
	{
		// User Management Routes (formerly in manage.go)
		// POST /client -> clientCtrl.CreateUser
		client.POST("", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.CreateUser)
		// GET /client -> clientCtrl.ListUsers
		client.GET("", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager), string(model.RoleUser)), clientCtrl.ListUsers)
		// GET /client/:id -> clientCtrl.GetUser
		client.GET("/:id", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager), string(model.RoleUser)), clientCtrl.GetUser)
		// PUT /client/:id -> clientCtrl.UpdateUser
		client.PUT("/:id", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.UpdateUser)
		// DELETE /client/:id -> clientCtrl.DeleteUser
		client.DELETE("/:id", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.DeleteUser)

		// Pause and Resume client routes
		client.POST("/:username/pause", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.PauseClient)
		client.POST("/:username/resume", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.ResumeClient)

		// Client Config Download (accessible by user for their own config, and admins/managers)
		// Path changed from /config/:username to /:id/config
		// Note: The GetClientConfig route uses /:username, matching Pause/Resume. The :id param is used for other user operations.
		client.GET("/config/:username", clientCtrl.GetClientConfig) // Controller logic should enforce user can only get own
	}
}