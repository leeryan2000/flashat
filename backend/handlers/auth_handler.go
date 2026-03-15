package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/utils"
)

type AuthHandler struct{ Repo repo.UserRepo }

// `json:"email"` tells json encode/decoder how to map go struct fields to json key
type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h AuthHandler) Login(c *gin.Context) {
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.HashedPassword) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token, err := utils.GenerateToken(user.UID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// ***** should later implement session
	// ***** change the SameSite option for the cookie
	c.SetCookie("token", token, 3*24*60*60, "", "", false, true)

	c.JSON(http.StatusOK, user)

}

func (h AuthHandler) Logout(c *gin.Context) {
	// ***** see if there is a way to invalidate the JWT token
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h AuthHandler) GetCurrentUser(c *gin.Context) {
	uidStr := c.GetString("uid")

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UID format"})
		return
	}

	user, err := h.Repo.GetUserByUID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
