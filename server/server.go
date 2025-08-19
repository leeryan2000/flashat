package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/db"
	repo "github.com/leeryan2000/flashat/repo/impl"
	"github.com/leeryan2000/flashat/service"
	"gorm.io/gorm"
)

type Server struct {
	DB   *gorm.DB
	Hub  *Hub
	Pool *pgxpool.Pool

	// Service
	MessageService *service.MessageService

	// Repo
	UserRepo         *repo.GormUserRepo
	ConversationRepo *repo.PgxConversationRepo
	MessageRepo      *repo.PgxMessageRepo
}

func StartServer() (*Server, error) {
	s := &Server{}

	// Database connection
	dbConnection, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	s.DB = dbConnection

	// Pgxpool setup
	pgx, err := db.NewPgxPool()
	if err != nil {
		return nil, err
	}
	s.Pool = pgx

	// Hub creation
	s.Hub = NewHub()
	go s.Hub.Run()

	// Service
	// Message service creation
	s.MessageService = &service.MessageService{
		Hub:  s.Hub,
		Repo: &repo.PgxMessageRepo{Pool: s.Pool},
	}

	// Repo
	s.UserRepo = &repo.GormUserRepo{DB: s.DB}
	s.ConversationRepo = &repo.PgxConversationRepo{Pool: s.Pool}
	s.MessageRepo = &repo.PgxMessageRepo{Pool: s.Pool}

	return s, nil
}
