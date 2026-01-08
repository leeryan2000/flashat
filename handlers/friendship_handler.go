package handlers

import (
	"log"

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
		log.Println(err)
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

	c.JSON(200, gin.H{"message": "Friendship accepted"})
}

func (h *FriendshipHandler) DeleteFriendship(c *gin.Context) {
}

func (h *FriendshipHandler) BlockUser(c *gin.Context) {
}

func (h *FriendshipHandler) ListFriendships(c *gin.Context) {
}

func (h *FriendshipHandler) GetFriendshipStatus(c *gin.Context) {
}
