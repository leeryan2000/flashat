package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/wire"
)

type ConversationHandler struct {
	Repo repo.ConversationRepo
}

type groupInput struct {
	GroupName    string   `json:"name"`
	ParticipantIds []string `json:"participantIds"` // UIDs of participants
}

func (h *ConversationHandler) CreateGroupConversation(c *gin.Context) {
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
	participantsUID := make([]uuid.UUID, 0, len(input.ParticipantIds))
	for _, uidStr := range input.ParticipantIds {
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

	// set up the url of default group avatar, later changed by user
	avatar_url := "Default group avatar"

	conv := &models.Conversation{
		ID:             uuid.New(),
		Type:           "group",
		GroupName:      &input.GroupName,
		GroupAvatarUrl: &avatar_url,
	}

	initialText := "Group Conversation created, you can start chatting now!"

	msg := &models.Message{
		ID:             uuid.New(),
		ConversationID: conv.ID,
		FromUID:        creatorUID,
		Body:           json.RawMessage(`{"text":"` + initialText + `"}`),
	}

	err = h.Repo.CreateGroupConversation(c.Request.Context(), conv, msg, creatorUID, participantsUID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group conversation"})
		return
	}

	log.Println(conv.CreatedAt, " Group Conversation Created")
	log.Println(msg.CreatedAt, " Message Created")

	convSummary := &wire.ConversationSummary{
		ConversationID: conv.ID,
		ConvType:       conv.Type,
		Title:          *conv.GroupName,
		AvatarURL:      *conv.GroupAvatarUrl,
		LastMsgID:      msg.ID,
		LastMsgText:    initialText,
		LastMsgFrom:    creatorUID,
		LastMsgTs:      msg.CreatedAt,
		LastSeq:        msg.Seq,
		LastReadSeq:    0,
		UnreadCount:    1,
	}

	type response struct {
		Conversation wire.ConversationSummary `json:"conversation"`
		Message      models.Message           `json:"message"`
	}

	resp := &response{
		Conversation: *convSummary,
		Message:      *msg,
	}

	c.JSON(http.StatusOK, resp)
}

type directInput struct {
	TargetUID string `json:"target_uid"` // UID of the target user
}

func (h *ConversationHandler) CreateDirectConversation(c *gin.Context) {
	// get the uid saved from auth
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse UID failed"})
		return
	}

	var input directInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid target UID"})
		return
	}

	targetUID, err := uuid.Parse(input.TargetUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse target UID failed"})
		return
	}

	if uid == targetUID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot create direct conversation with self"})
		return
	}

	conv := &models.Conversation{
		ID:   uuid.New(),
		Type: "direct",
	}

	err = h.Repo.CreateDirectConversation(c.Request.Context(), conv, uid, targetUID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create direct conversation"})
		return
	}

	log.Println(conv.CreatedAt, " Conversation Created")

	c.JSON(http.StatusOK, conv)
}

func (h *ConversationHandler) ListConversationByUID(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse UID failed"})
		return
	}

	conversations, err := h.Repo.ListConversationByUID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to list conversations"})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

func (h *ConversationHandler) GetConversationByID(c *gin.Context) {
	conversationiIDStr := c.Param("conversation_id")
	conversationID, err := uuid.Parse(conversationiIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	conversation, err := h.Repo.GetConversationByID(c.Request.Context(), conversationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation"})
		return
	}

	if conversation == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func (h *ConversationHandler) ListParticipantByID(c *gin.Context) {
	conversationIDStr := c.Param("conversation_id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	participants, err := h.Repo.ListParticipantByID(c.Request.Context(), conversationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get participant"})
		return
	}

	if participants == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Participant not found"})
		return
	}

	c.JSON(http.StatusOK, participants)
}

type participantInput struct {
	ParticipantUID string `json:"participant_uid"` // UID of the participant to add or modify
	ConversationID string `json:"conversation_id"` // ID of the conversation
	Role           string `json:"role"`            // Role of the participant (e.g., "member", "admin")
}

func (h *ConversationHandler) AddParticipant(c *gin.Context) {
	participantInput := &participantInput{}
	if err := c.ShouldBindJSON(participantInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	participantUID, err := uuid.Parse(participantInput.ParticipantUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse participant UID failed"})
		return
	}

	conversationID, err := uuid.Parse(participantInput.ConversationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	err = h.Repo.AddParticipant(c.Request.Context(), conversationID, participantUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participant added successfully",
	})
}

func (h *ConversationHandler) ModifyParticipant(c *gin.Context) {
	participantInput := &participantInput{}
	if err := c.ShouldBindJSON(participantInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	participantUID, err := uuid.Parse(participantInput.ParticipantUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse participant UID failed"})
		return
	}

	conversationID, err := uuid.Parse(participantInput.ConversationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	role := participantInput.Role
	if role == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Role didn't provided"})
		return
	}

	err = h.Repo.ModifyParticipant(c.Request.Context(), conversationID, participantUID, role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to modify participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participant modified successfully",
	})
}

func (h *ConversationHandler) RemoveParticipant(c *gin.Context) {
	participantInput := &participantInput{}
	if err := c.ShouldBindJSON(participantInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	participantUID, err := uuid.Parse(participantInput.ParticipantUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse participant UID failed"})
		return
	}

	conversationID, err := uuid.Parse(participantInput.ConversationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	err = h.Repo.RemoveParticipant(c.Request.Context(), conversationID, participantUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participant removed successfully",
	})
}

type nextSeqInput struct {
	ConversationID string `json:"conversation_id"` // ID of the conversation
	LastReadSeq    int64  `json:"last_read_seq"`   // Last read sequence number
}

func (h *ConversationHandler) UpdateLastReadSeq(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse participant UID failed"})
		return
	}

	seqInput := &nextSeqInput{}
	if err := c.ShouldBindJSON(seqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	conversationID, err := uuid.Parse(seqInput.ConversationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	err = h.Repo.UpdateLastReadSeq(c.Request.Context(), conversationID, uid, seqInput.LastReadSeq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update last read sequence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Last read sequence updated successfully",
	})
}

func (h *ConversationHandler) GetLastReadSeq(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse participant UID failed"})
		return
	}

	conversationIDStr := c.Param("conversation_id")
	if conversationIDStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ConversationID is required"})
		return
	}
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	lastReadSeq, err := h.Repo.GetLastReadSeq(c.Request.Context(), conversationID, uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last read sequence"})
		return
	}

	c.JSON(http.StatusOK, lastReadSeq)
}

func (h *ConversationHandler) LeaveGroup(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse UID failed"})
		return
	}

	conversationIDStr := c.Param("conversation_id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse conversation ID failed"})
		return
	}

	err = h.Repo.LeaveGroup(c.Request.Context(), conversationID, uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left group successfully"})
}

func (h *ConversationHandler) GetSummary(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Parse participant UID failed"})
		return
	}

	summary, err := h.Repo.GetSummary(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
