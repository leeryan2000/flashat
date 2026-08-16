package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat-posts/internal/genproto/authpb"
)

// Authenticate is posts' equivalent of backend/middleware.Authenticate,
// except the session lookup happens over gRPC instead of a direct Redis
// call, since the posts service doesn't have its own Redis access by
// design — the monolith remains the sole owner of session data.
func Authenticate(authClient authpb.AuthInternalClient) gin.HandlerFunc {
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

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		resp, err := authClient.ValidateSession(ctx, &authpb.ValidateSessionRequest{SessionId: sessionID})
		if err != nil || !resp.GetValid() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			return
		}

		c.Set("uid", resp.GetUid())
		c.Next()
	}
}
