package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
)

func InitializeUserRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	uh := h.User
	userRoutes := router.Group("/user")
	
	userRoutes.GET("/getAllUsers", uh.GetAllUsers)
	userRoutes.GET("/getUserById/:id", uh.GetUserById)

	userRoutes.POST("/createUser", uh.CreateUser)
}
