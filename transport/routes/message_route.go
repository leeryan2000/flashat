package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
)

func InitializeMessageRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	mh := h.Message
	messageRoutes := router.Group("/message")
	messageRoutes.Use(middleware.Authenticate())

	messageRoutes.GET("/latest/:conversation_id", mh.ListLatest)
}
