package handlers

import (
	"net/http"
	"strconv"

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

	// ***** failed to add a user into database would cause increment of in primary id
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
	idParam := c.Param("id")

	if idParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.Repo.GetUserById(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
