package repo

import (
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"gorm.io/gorm"
)

type GormUserRepo struct {
	DB *gorm.DB
}

func (r *GormUserRepo) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *GormUserRepo) GetUserByUID(uid uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.DB.First(&user, "uid = ?", uid).Error
	return user, err
}

func (r *GormUserRepo) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
