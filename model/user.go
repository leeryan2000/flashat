package model

// import (
// 	"time"
// )

type User struct {
	Hashed_Password string `json:"name"`
	Email           string `json:"email" gorm:"unique"`
}
