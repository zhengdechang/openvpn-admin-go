package router

import (
	"openvpn-admin-go/controller"
	"openvpn-admin-go/middleware"
	"openvpn-admin-go/model"

	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes registers notification endpoints (superadmin only)
func SetupNotificationRoutes(r *gin.RouterGroup) {
	ctrl := &controller.NotificationController{}
	g := r.Group("/notifications")
	g.Use(middleware.JWTAuthMiddleware())
	g.Use(middleware.RoleRequired(string(model.RoleSuperAdmin)))
	{
		g.GET("", ctrl.List)
		g.GET("/unread-count", ctrl.UnreadCount)
		// IMPORTANT: "read-all" must be registered before "/:id/read"
		// to prevent Gin from matching "read-all" as the :id parameter.
		g.PATCH("/read-all", ctrl.MarkAllRead)
		g.PATCH("/:id/read", ctrl.MarkRead)
	}
}
