package router

import (
   "openvpn-admin-go/controller"
   "openvpn-admin-go/middleware"
   "openvpn-admin-go/model"

   "github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户和权限相关路由
func SetupUserRoutes(r *gin.RouterGroup) {
   user := r.Group("/user")
   {
       user.POST("/register", controller.Register)
       user.POST("/login", controller.Login)
       user.GET("/verify-email/:token", controller.VerifyEmail)
       user.POST("/forgot-password", controller.ForgotPassword)
       user.PATCH("/reset-password/:token", controller.ResetPassword)
       user.GET("/me", middleware.JWTAuthMiddleware(), controller.GetMe)
       user.PATCH("/me", middleware.JWTAuthMiddleware(), controller.UpdateMe)
       user.POST("/logout", middleware.JWTAuthMiddleware(), controller.Logout)
       user.GET("/refresh", middleware.JWTAuthMiddleware(), controller.RefreshToken)
       user.GET("/roles", middleware.JWTAuthMiddleware(), controller.GetRoles)
       // 查询用户信息: superadmin/admin/manager
       user.GET("/info/:id",
           middleware.JWTAuthMiddleware(),
           middleware.RoleRequired(string(model.RoleSuperAdmin), string(model.RoleAdmin), string(model.RoleManager)),
           controller.GetUserInfo,
       )
   }
}