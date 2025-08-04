package repo

import (
	"github.com/leeryan2000/flashat/models"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func (r *GormUserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}