package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/repo"
)

type ConversationHandler struct {
	Repo repo.ConversationRepo
}

type groupInput struct {
	GroupName    string   `json:"group_name"`
	Participants []string `json:"participants"` // UIDs of participants
}

func (h ConversationHandler) CreateGroupConversation(c *gin.Context) {
	// get the uid saved from auth
	creatorStr := c.GetString("uid")
	creatorUID, err := uuid.Parse(creatorStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse UID failed"})
		return
	}

	var input groupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	set := map[uuid.UUID]struct{}{creatorUID: {}}
	participantsUID := make([]uuid.UUID, 0, len(input.Participants))
	for _, uidStr := range input.Participants {
		uid, err := uuid.Parse(uidStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse UID failed"})
			return
		}
		if _, exists := set[uid]; !exists {
			participantsUID = append(participantsUID, uid)
			set[uid] = struct{}{}
		}
	}

	if err := h.Repo.CreateGroupConversation(c.Request.Context(), creatorUID, participantsUID, input.GroupName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group conversation created successfully"})
}
