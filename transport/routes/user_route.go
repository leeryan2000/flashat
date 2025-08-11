package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middlewares"
)

func InitializeUserRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	uh := h.User

	userRoutes := router.Group("/user")
	userRoutes.POST("/createUser", uh.CreateUser)

	// User authentication setup here
	userAuthRoutes := router.Group("/user/auth")
	userAuthRoutes.Use(middlewares.Authenticate())

	userAuthRoutes.GET("/getAllUsers", uh.GetAllUsers)
	userAuthRoutes.GET("/getUserById/:id", uh.GetUserById)
}
