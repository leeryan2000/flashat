package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/model"
	"github.com/leeryan2000/flashat/server"
	"log"
	"net/http"
)

type UserHandler struct{ S *server.Server }

func (uh UserHandler) CreateUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.S.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal("Failed to add user to server")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh UserHandler) GetUsers(c *gin.Context) {
	var users []model.User
	fmt.Println("here")
	uh.S.DB.Find(&users)
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
