package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
	"github.com/leeryan2000/flashat/server"
)

func InitializeAuthRoutes(router *gin.RouterGroup, h *handlers.Handlers, s *server.Server) {
	ah := h.Auth
	router.POST("/login", ah.Login)

	authRoutes := router.Group("/auth")
	authRoutes.Use(middleware.Authenticate(s))
	authRoutes.DELETE("/logout", ah.Logout)
	authRoutes.GET("/me", ah.GetCurrentUser)
}
