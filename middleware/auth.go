package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/utils"
)

// In front end set the "token" value in localstorage, send request and set header with key: Authorization, value: <token>
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid session token"})
			return
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// Store UID in context for downstream handlers
		c.Set("uid", claims.UID)

		c.Next()
	}
}
