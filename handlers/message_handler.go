package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/repo"
)

type MessageHandler struct {
	Repo repo.MessageRepo
}

func (h *MessageHandler) ListLatest(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UID"})
		return
	}

	conversationIDStr := c.Param("conversation_id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid conversation ID"})
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	msgs, err := h.Repo.ListLatest(c.Request.Context(), uid, conversationID, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, msgs)
}

func (h *MessageHandler) ListBefore(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UID"})
		return
	}

	conversationIDStr := c.Param("conversation_id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid conversation ID"})
		return
	}

	seqStr := c.Query("seq")
	seq, err := strconv.ParseInt(seqStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid sequence number"})
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	msgs, err := h.Repo.ListBefore(c.Request.Context(), uid, conversationID, seq, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, msgs)
}
