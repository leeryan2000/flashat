package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/server"
)

func Authenticate(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			// ***** test version which reads from query to test through postman
			sessionID = c.Query("session_id")
			if sessionID == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid session"})
				return
			}
		}

		uid, err := s.RedisClient.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			// err is redis.Nil when session expired or not found
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			return
		}

		c.Set("uid", uid)
		c.Next()
	}
}
