package repo

import (
	"github.com/leeryan2000/flashat/model"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func (r *GormUserRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := r.DB.Find(&users).Error
	return users, err
}