package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		tokenString := parts[1]
		_ = tokenString

		// nanti kita verify JWT di sini
		// contoh sementara:

		userID := 1
		role := "admin"

		// set ke context
		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}