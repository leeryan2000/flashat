package models

import (
	"time"
)

type User struct {
	ID             uint      `json:"id"`
	UID            string    `json:"uid"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	UserAvatarURL  string    `json:"user_avatar_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}