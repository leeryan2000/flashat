package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
)

type UserHandler struct{ Repo repo.UserRepo }

func (uh UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.Repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal("Failed to add user to server")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uh UserHandler) GetAllUsers(c *gin.Context) {
	var users []models.User
	users, err := uh.Repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// func CreateUserLocal(db *gorm.DB) {
//     var user = models.User {
// 		Email:"ryan1",
// 		Hashed_Password:"testpassword",
// 	}

//     if err := db.Create(&user).Error; err != nil {
// 		log.Fatal("Failed to add user to server")
// 		return
//     }
// }
