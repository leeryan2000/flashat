package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/utils"
)

// In front end set the "token" value in localstorage, send request and set header with key: Authorization, value: <token>
func Authenticate(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid session token"})
			// return

			// ***** test version which read from query to test through postman
			tokenStr = c.Query("token")
			if tokenStr == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid session token"})
				return
			}
			// -------------------------
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// Store UID in context for downstream handlers
		uid, err := uuid.Parse(claims.UID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid UID in token"})
			return
		}

		user, err := s.UserRepo.GetUserByUID(uid)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("uid", claims.UID)

		c.Next()
	}
}
