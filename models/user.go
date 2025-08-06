package models

import (
	"errors"
	"strings"
)

// import (
// 	"time"
// )

type User struct {
	Id             uint   `gorm:"primaryKey"`
	Email          string `gorm:"unique" json:"email" `
	HashedPassword string `gorm:"column:hashed_password" json:"-"`
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
