package models

import (
	"errors"
	"strings"
	"time"
)

// import (
// 	"time"
// )

type User struct {
	ID             uint   `gorm:"primaryKey"`
	UID            string `gorm:"type:uuid;uniqueIndex" json:"uid"`
	Email          string `gorm:"unique" json:"email" `
	HashedPassword string `gorm:"column:hashed_password" json:"-"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *User) Validate() error {
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email")
	}
	if len(u.HashedPassword) < 6 {
		return errors.New("password too short")
	}
	return nil
}
