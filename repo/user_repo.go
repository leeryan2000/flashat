package repo

import (
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
)

type UserRepo interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]models.User, error)
	GetUserByUID(uid uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
