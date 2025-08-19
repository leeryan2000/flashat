package models

import (
	"time"
)

type User struct {
	ID             uint   `gorm:"primaryKey"`
	UID            string `gorm:"type:uuid;uniqueIndex" json:"uid"`
	Email          string `gorm:"unique" json:"email" `
	HashedPassword string `gorm:"column:hashed_password" json:"-"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
