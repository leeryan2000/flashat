package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
)

func SetupRoutes(handlers *handlers.Handlers) *gin.Engine {
	router := gin.Default()
	mainRouter := router.Group("/api")

	InitializeUserRoutes(mainRouter, handlers)
	
	mainRouter.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Ok",
		})
	})

	return router
}
