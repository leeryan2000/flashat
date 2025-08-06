package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/utils"
	"gorm.io/gorm"
)

type User struct {
	ID             uint   `gorm:"primaryKey"`
	UID            string `gorm:"type:uuid;uniqueIndex" json:"uid"`
	Email          string `gorm:"unique" json:"email" `
	HashedPassword string `gorm:"column:hashed_password" json:"-"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	// create UID
	u.UID = uuid.NewString()
	// Hash the password
	hashedPassword, err := utils.HashPassword(u.HashedPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}
	u.HashedPassword = hashedPassword

	return nil
}
