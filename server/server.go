package server

import (
	"log"
	"time"

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
	FriendshipRepo   *repo.PgxFriendshipRepo
}

// ***** Test
func RunPeriodicTask(task func()) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			task()
		}
	}()
}

func checkClients(s *Server) {
	log.Println("Checking connected clients...", len(s.Hub.ClientsByUID))
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
		Hub:              s.Hub,
		MessageRepo:      &repo.PgxMessageRepo{Pool: s.Pool},
		ConversationRepo: &repo.PgxConversationRepo{Pool: s.Pool},
	}

	// Repo
	s.UserRepo = &repo.GormUserRepo{DB: s.DB}
	s.MessageRepo = &repo.PgxMessageRepo{Pool: s.Pool}
	s.ConversationRepo = &repo.PgxConversationRepo{Pool: s.Pool}
	s.FriendshipRepo = &repo.PgxFriendshipRepo{Pool: s.Pool}

	RunPeriodicTask(func() {
		checkClients(s)
	})

	return s, nil
}
