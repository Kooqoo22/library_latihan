package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			return
		}

		c.Next()
	}
}