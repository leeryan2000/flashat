package models

import (
	"time"
)

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UID            string    `gorm:"type:uuid;uniqueIndex" json:"uid"`
	Name           string    `gorm:"type:varchar(100)" json:"name"`
	Email          string    `gorm:"unique" json:"email"`
	HashedPassword string    `gorm:"column:hashed_password" json:"-"`
	UserAvatarURL  *string   `gorm:"column:user_avatar_url" json:"user_avatar_url"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}
