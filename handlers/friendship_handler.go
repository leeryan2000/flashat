package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/repo"
)

type FriendshipHandler struct {
	Repo repo.FriendshipRepo
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

	err = h.Repo.AcceptFriendship(c.Request.Context(), requesterUID, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to accept friendship"})
		return
	}

	// ***** return the created direct conversation information to the client
	c.JSON(200, gin.H{"message": "Friendship accepted"})
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

func (h *FriendshipHandler) ListFriendshipRequests(c *gin.Context) {
	uidStr := c.GetString("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parse UID failed"})
		return
	}

	requests, err := h.Repo.ListFriendshipRequests(c.Request.Context(), uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list friendship requests"})
		return
	}

	c.JSON(200, requests)
}

func (h *FriendshipHandler) GetFriendshipStatus(c *gin.Context) {
}
