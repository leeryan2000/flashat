package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/utils"
)

type UserHandler struct{ Repo repo.UserRepo }

type createUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h UserHandler) CreateUser(c *gin.Context) {
	// c.ShouldBindJSON can only read once for a struct body
	var createUserInput createUserInput

	if err := c.ShouldBindJSON(&createUserInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to bind create user input"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(createUserInput.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		UID:            uuid.NewString(),
		Name:           createUserInput.Name,
		Email:          createUserInput.Email,
		HashedPassword: hashedPassword,
	}

	// ***** after redis implemented, turn this into storing in redis first, then after verification store in main db through another handler function
	if err := h.Repo.CreateUser(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to server"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Repo.GetAllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h UserHandler) GetUserById(c *gin.Context) {
	uidStr := c.Param("uid")

	if uidStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.Repo.GetUserByUID(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) UpdateName(c *gin.Context) {
	uidStr := c.Param("uid")

	if uidStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to bind input"})
		return
	}

	user, err := h.Repo.GetUserByUID(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) SendVerification(c *gin.Context) {
	// otp := utils.SendVerification()
	// ***** store otp in redis with expiration associated with user email
	c.JSON(http.StatusOK, gin.H{"message": "Verification sent"})
}
