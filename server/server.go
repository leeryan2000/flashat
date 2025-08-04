package server

import (
	"github.com/leeryan2000/flashat/db"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
}

func StartServer() (*Server, error) {
	s := &Server{}

	dbConnection, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	s.DB = dbConnection

	return s, nil
}

