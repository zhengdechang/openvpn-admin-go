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

		// Client Config Download (accessible by user for their own config, and admins/managers)
		// Path changed from /config/:username to /:userId/config
		client.GET("/:userId/config", clientCtrl.GetClientConfig) // Controller logic should enforce user can only get own

		// Client Status Routes
		statusGroup := client.Group("/status")
		// GET /client/status -> clientCtrl.GetAllClientStatuses (summary for all clients)
		statusGroup.GET("", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.GetAllClientStatuses)
		// GET /client/status/live -> clientCtrl.GetLiveConnections
		statusGroup.GET("/live", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.GetLiveConnections)
		// GET /client/status/detailed -> clientCtrl.GetDetailedClientStatuses
		statusGroup.GET("/detailed", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.GetDetailedClientStatuses)
		// GET /client/status/summary -> clientCtrl.GetClientList (assuming GetClientList provides a summary list)
        statusGroup.GET("/summary", middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)), clientCtrl.GetClientList)


		// Individual Client Status
		// Path changed from /status/:username to /:userId/status
		client.GET("/:userId/status", clientCtrl.GetClientStatus) // Controller logic should enforce user can only get own if not admin/manager

		// Client OpenVPN Configuration Management Actions (for existing users)
		// These routes are for managing OpenVPN specific configurations for a user (client)
		clientActions := client.Group("/:userId/actions")
		clientActions.Use(middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)))
		{
			// POST /client/:userId/actions/create-config -> clientCtrl.AddClient (generates certs for an existing user)
			clientActions.POST("/create-config", clientCtrl.AddClient)
			// PUT /client/:userId/actions/update-config -> clientCtrl.UpdateClient (regenerates certs)
			clientActions.PUT("/update-config", clientCtrl.UpdateClient) // Note: UpdateClient in controller might need adjustment if it takes username from body
			// DELETE /client/:userId/actions/delete-config -> clientCtrl.DeleteClient (deletes certs)
			clientActions.DELETE("/delete-config", clientCtrl.DeleteClient)
			// POST /client/:userId/actions/revoke-config -> clientCtrl.RevokeClient (revokes certs)
			clientActions.POST("/revoke-config", clientCtrl.RevokeClient)
			// POST /client/:userId/actions/renew-config -> clientCtrl.RenewClient (renews/regenerates certs)
			clientActions.POST("/renew-config", clientCtrl.RenewClient)
			// POST /client/:userId/actions/pause -> clientCtrl.PauseClient
			clientActions.POST("/pause", clientCtrl.PauseClient)
			// POST /client/:userId/actions/resume -> clientCtrl.ResumeClient
			clientActions.POST("/resume", clientCtrl.ResumeClient)
		}
	}
}