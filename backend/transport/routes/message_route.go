package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
	"github.com/leeryan2000/flashat/server"
)

func InitializeMessageRoutes(router *gin.RouterGroup, h *handlers.Handlers, s *server.Server) {
	mh := h.Message
	messageRoutes := router.Group("/message")
	messageRoutes.Use(middleware.Authenticate(s))

	messageRoutes.GET("/latest/:conversation_id", mh.ListLatest)
	messageRoutes.GET("/before/:conversation_id", mh.ListBefore)
}
