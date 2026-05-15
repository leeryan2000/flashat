package server

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/config"
	"github.com/leeryan2000/flashat/db"
	repo "github.com/leeryan2000/flashat/repo/impl"
	"github.com/leeryan2000/flashat/service"
)

type Server struct {
	Hub            *Hub
	Pool           *pgxpool.Pool
	RedisClient    *db.RedisClient
	RabbitMQClient *db.RabbitMQClient

	// Service
	MessageService *service.MessageService
	MessageWorker  *service.MessageWorker

	// Repo
	UserRepo         *repo.PgxUserRepo
	ConversationRepo *repo.PgxConversationRepo
	MessageRepo      *repo.PgxMessageRepo
	FriendshipRepo   *repo.PgxFriendshipRepo

	// Config
	RegisterCode string
}

// ***** Test
func RunPeriodicTask(task func()) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			task()
		}
	}()
}

func checkClients(s *Server) {
	log.Println("Checking connected clients...", len(s.Hub.Clients), "connections,", len(s.Hub.ClientsByUID), "unique users")
}

func StartServer() (*Server, error) {
	s := &Server{}

	cfg := config.LoadConfig()

	err := db.RunMigrations(cfg)
	if err != nil {
		return nil, err
	}

	pgx, err := db.NewPgxPool(cfg)
	if err != nil {
		return nil, err
	}
	s.Pool = pgx

	redisClient, err := db.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	s.RedisClient = redisClient

	rabbitMQClient, err := db.NewRabbitMQClient(cfg)
	if err != nil {
		return nil, err
	}
	s.RabbitMQClient = rabbitMQClient

	s.Hub = NewHub()
	go s.Hub.Run()

	s.UserRepo = &repo.PgxUserRepo{Pool: s.Pool}
	s.MessageRepo = &repo.PgxMessageRepo{Pool: s.Pool}
	s.ConversationRepo = &repo.PgxConversationRepo{Pool: s.Pool}
	s.FriendshipRepo = &repo.PgxFriendshipRepo{Pool: s.Pool}
	s.RegisterCode = cfg.REGISTER_CODE

	s.MessageService = service.NewMessageService(
		s.Hub,
		s.MessageRepo,
		s.ConversationRepo,
		s.RabbitMQClient,
	)

	s.MessageWorker = service.NewMessageWorker(s.RabbitMQClient, s.MessageService)
	s.MessageWorker.Start()

	RunPeriodicTask(func() {
		checkClients(s)
	})

	return s, nil
}
