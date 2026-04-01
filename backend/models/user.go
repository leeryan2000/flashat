package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UID            uuid.UUID `json:"uid"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	UserAvatarURL  string    `json:"user_avatar_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
