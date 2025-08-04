package repo

import "github.com/leeryan2000/flashat/models"

type UserRepo interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]models.User, error)
	GetUserById(id uint) (*models.User, error)
}
