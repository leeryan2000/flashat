package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/server"
)

func SetupRoutes(handlers *handlers.Handlers, server *server.Server) *gin.Engine {
	router := gin.Default()
	// cors config — production traffic is same-origin (nginx serves the
	// frontend and proxies /api/ under the same domain), so these origins
	// only matter for local dev and direct cross-origin API callers.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://flashatapp.com", "https://www.flashatapp.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		AllowWebSockets:  true,
		AllowFiles:       true,
	}))

	mainRouter := router.Group("/api")

	InitializeUserRoutes(mainRouter, handlers, server)
	InitializeAuthRoutes(mainRouter, handlers, server)
	InitializeWebsocketRoutes(mainRouter, handlers, server)
	InitializeConversationRoutes(mainRouter, handlers, server)
	InitializeMessageRoutes(mainRouter, handlers, server)
	InitializeFriendshipRoutes(mainRouter, handlers, server)

	mainRouter.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Ok",
		})
	})

	return router
}
