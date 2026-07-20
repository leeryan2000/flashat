package server

import (
	"context"
	"log/slog"
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
	GoEnv        string
}

func RunPeriodicTask(ctx context.Context, task func()) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				task()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func checkClients(s *Server) {
	clients, uniqueUsers := s.Hub.Counts()
	slog.Info("connected clients", "connections", clients, "unique_users", uniqueUsers)
}

// ctx is a shutdown signal, threaded down to MessageWorker so it can stop
// picking up new deliveries and bound in-flight ones on shutdown.
func StartServer(ctx context.Context) (*Server, error) {
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

	s.UserRepo = &repo.PgxUserRepo{Pool: s.Pool}
	s.MessageRepo = &repo.PgxMessageRepo{Pool: s.Pool}
	s.ConversationRepo = &repo.PgxConversationRepo{Pool: s.Pool}
	s.FriendshipRepo = &repo.PgxFriendshipRepo{Pool: s.Pool}
	s.RegisterCode = cfg.REGISTER_CODE
	s.GoEnv = cfg.GO_ENV

	s.MessageService = service.NewMessageService(
		s.Hub,
		s.MessageRepo,
		s.ConversationRepo,
		s.RabbitMQClient,
	)
	s.MessageService.RegisterObserver(&service.AckObserver{Hub: s.Hub})
	s.MessageService.RegisterObserver(&service.BroadcastObserver{Hub: s.Hub})

	s.MessageWorker = service.NewMessageWorker(s.RabbitMQClient, s.MessageService)
	s.MessageWorker.Start(ctx)

	RunPeriodicTask(ctx, func() {
		checkClients(s)
	})

	return s, nil
}
