package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
)

func InitializeUserRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	uh := h.User

	userRoutes := router.Group("/user")
	userRoutes.POST("/", uh.CreateUser)

	// User authentication setup here
	userAuthRoutes := router.Group("/user/auth")
	userAuthRoutes.Use(middleware.Authenticate())

	userAuthRoutes.GET("/all", uh.GetAllUsers)
	userAuthRoutes.GET("/:id", uh.GetUserById)
}
