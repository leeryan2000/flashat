package routes

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/server"
)

func SetupRoutes(handlers *handlers.Handlers, server *server.Server) *gin.Engine {
	router := gin.Default()
	// cors config
	// ***** modify cors settings for production
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			log.Println("CORS Origin:", origin)
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		AllowWebSockets:  true,
		AllowFiles:       true,
	}))

	router.Static("/static", "./static")
	router.StaticFile("/chat", "./static/chat.html")

	mainRouter := router.Group("/api")

	InitializeUserRoutes(mainRouter, handlers, server)
	InitializeAuthRoutes(mainRouter, handlers, server)
	InitializeWebsocketRoutes(mainRouter, handlers, server)
	InitializeConversationRoutes(mainRouter, handlers, server)
	InitializeMessageRoutes(mainRouter, handlers, server)

	mainRouter.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Ok",
		})
	})

	return router
}
