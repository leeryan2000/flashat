package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/db"
	"gorm.io/gorm"
)

type Server struct {
	DB  *gorm.DB
	Hub *Hub
	Pgx *pgxpool.Pool
}

func StartServer() (*Server, error) {
	s := &Server{}

	dbConnection, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	s.DB = dbConnection

	pgx, err := db.NewPgxPool()
	if err != nil {
		return nil, err
	}
	s.Pgx = pgx

	s.Hub = NewHub()
	// start hub
	go s.Hub.Run()

	return s, nil
}
