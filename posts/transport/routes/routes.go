package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat-posts/handlers"
	"github.com/leeryan2000/flashat-posts/internal/genproto/authpb"
)

func SetupRoutes(h *handlers.Handlers, authClient authpb.AuthInternalClient) *gin.Engine {
	router := gin.Default()
	// Same CORS shape as backend/transport/routes/routes.go — production
	// traffic is same-origin (nginx proxies /api/posts/ under the same
	// domain), these origins only matter for local dev/direct callers.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://flashatapp.com", "https://www.flashatapp.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		AllowFiles:       true,
	}))

	mainRouter := router.Group("/api/posts")

	mainRouter.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Ok"})
	})

	InitializePostRoutes(mainRouter, h, authClient)

	return router
}
