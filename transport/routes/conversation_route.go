package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middlewares"
)

func InitializeConversationRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	ch := h.Conversation
	conversationRoutes := router.Group("/conversation")
	conversationRoutes.Use(middlewares.Authenticate())

	conversationRoutes.POST("/create", ch.CreateGroupConversation)
}
