package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
	"github.com/leeryan2000/flashat/server"
)

func InitializeConversationRoutes(router *gin.RouterGroup, h *handlers.Handlers, s *server.Server) {
	ch := h.Conversation
	conversationRoutes := router.Group("/conversation")
	conversationRoutes.Use(middleware.Authenticate(s))

	conversationRoutes.POST("/group", ch.CreateGroupConversation)
	conversationRoutes.POST("/direct", ch.CreateDirectConversation)

	conversationRoutes.GET("/", ch.ListConversationByUID)
	conversationRoutes.GET("/:conversation_id", ch.GetConversationByID)
	conversationRoutes.GET("/:conversation_id/participant", ch.ListParticipantByID)

	conversationRoutes.POST("/participant", ch.AddParticipant)
	conversationRoutes.PUT("/participant", ch.ModifyParticipant)
	conversationRoutes.DELETE("/participant", ch.RemoveParticipant)

	conversationRoutes.PUT("/read_seq", ch.UpdateLastReadSeq)
	conversationRoutes.GET("/read_seq", ch.GetLastReadSeq)

	conversationRoutes.GET("/summary", ch.GetSummary)
}
