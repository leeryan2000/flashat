package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/utils"
)

type UserHandler struct{ Repo repo.UserRepo }

type createUserInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (uh UserHandler) CreateUser(c *gin.Context) {
	// c.ShouldBindJSON can only read once for a struct body
	var createUserInput createUserInput

	if err := c.ShouldBindJSON(&createUserInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind model user"})
		return
	}
	log.Println(createUserInput.Password)

	user := models.User {
		Email: createUserInput.Email,
		HashedPassword: createUserInput.Password,
	}

	hashedPassword, err := utils.HashPassword(user.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	log.Println(hashedPassword)
	user.HashedPassword = hashedPassword

	// ***** failed to add a user into database would cause increment of in primary id
	if err := uh.Repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("Failed to add user to server")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uh UserHandler) GetAllUsers(c *gin.Context) {
	users, err := uh.Repo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uh UserHandler) GetUserById(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uh.Repo.GetUserById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
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
