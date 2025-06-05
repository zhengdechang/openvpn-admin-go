/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-04 10:37:43
 */
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleRequired 基于角色的访问控制中间件
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		claims := val.(*Claims)
		for _, r := range roles {
			if claims.Role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not allowed"})
	}
}
