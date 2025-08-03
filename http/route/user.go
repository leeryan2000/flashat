package route

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handler"
)

func InitializeUserRoutes(router *gin.RouterGroup, h *handler.Handlers) {
	uh := h.User
	userRoutes := router.Group("/user")
	userRoutes.GET("/getUsers", uh.GetUsers)
	userRoutes.POST("/createUser", uh.CreateUser)
}
