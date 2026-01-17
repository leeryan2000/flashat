package handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/server"
)

type FriendshipHandler struct {
	Hub      server.Hub
	Repo     repo.FriendshipRepo
	UserRepo repo.UserRepo
}

func (h *FriendshipHandler) RequestFriendship(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parse UID failed"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err = h.Repo.RequestFriendship(c.Request.Context(), uid, input.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to request friendship"})
		return
	}

	c.JSON(200, gin.H{"message": "Friendship request sent"})
}

func (h *FriendshipHandler) AcceptFriendship(c *gin.Context) {
	var input struct {
		RequesterUID string `json:"requester_uid" binding:"required"`
	}

	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parse UID failed"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	requesterUID, err := uuid.Parse(input.RequesterUID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid requester UID"})
		return
	}

	conv := &models.Conversation{
		ID:   uuid.New(),
		Type: "direct",
	}

	initialText := "I have accepted your friend request."

	msg := &models.Message{
		ID:             uuid.New(),
		ConversationID: conv.ID,
		FromUID:        uid,
		Body:           json.RawMessage(`{"text":"` + initialText + `"}`),
	}

	err = h.Repo.AcceptFriendship(c.Request.Context(), conv, msg, requesterUID, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to accept friendship"})
		return
	}

	// Retreive friend's profile info
	friendProfile, err := h.UserRepo.GetUserByUID(requesterUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get friend profile"})
		return
	}

	convSummary := &models.ConversationSummary{
		ConversationID: conv.ID,
		ConvType:       conv.Type,
		Title:          &friendProfile.Name,
		AvatarURL:      friendProfile.UserAvatarURL,
		LastMsgID:      &msg.ID,
		LastMsgText:    &initialText,
		LastMsgFrom:    &friendProfile.UID,
		LastMsgTs:      &msg.CreatedAt,
		LastSeq:        msg.Seq,
		LastReadSeq:    0,
		UnreadCount:    1,
	}

	// ***** return the created direct conversation information to the client
	c.JSON(200, convSummary)
}

func (h *FriendshipHandler) DeleteFriendship(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parse UID failed"})
		return
	}

	friendUIDstr := c.Param("friend_uid")
	friendUID, err := uuid.Parse(friendUIDstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid friend UID"})
		return
	}

	err = h.Repo.DeleteFriendship(c.Request.Context(), uid, friendUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete friendship"})
		return
	}

	c.JSON(200, gin.H{"message": "Friendship deleted"})
}

func (h *FriendshipHandler) RejectFriendship(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parse UID failed"})
		return
	}

	friendUIDstr := c.Param("friend_uid")
	friendUID, err := uuid.Parse(friendUIDstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid friend UID"})
		return
	}

	err = h.Repo.RejectFriendship(c.Request.Context(), friendUID, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to reject friendship"})
		return
	}

	c.JSON(200, gin.H{"message": "Friendship request rejected"})
}

func (h *FriendshipHandler) BlockUser(c *gin.Context) {
}

func (h *FriendshipHandler) ListFriendships(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parse UID failed"})
		return
	}

	friendUIDs, err := h.Repo.ListFriendships(c.Request.Context(), uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list friendships"})
		return
	}

	c.JSON(200, friendUIDs)
}

func (h *FriendshipHandler) GetFriendshipStatus(c *gin.Context) {
}
