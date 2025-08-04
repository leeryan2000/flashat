package models

// import (
// 	"time"
// )

type User struct {
	Id              uint   `gorm:"primaryKey"`
	Hashed_Password string `json:"name"`
	Email           string `json:"email" gorm:"unique"`
}
