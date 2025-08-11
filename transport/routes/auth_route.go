package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
)

func InitializeAuthRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	ah := h.Auth
	authRoutes := router.Group("/auth")

	authRoutes.POST("/login", ah.Login)
}
