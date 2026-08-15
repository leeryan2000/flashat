package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport/routes"
)

// BuildRouter wires up all handlers and returns the router without starting it.
// main.go owns the server lifecycle and calls Shutdown() for graceful cleanup.
func BuildRouter(s *server.Server) *gin.Engine {
	h := handlers.Handlers{
		User: handlers.UserHandler{
			Repo:         s.UserRepo,
			RegisterCode: s.RegisterCode,
			S3Client:     s.S3Client,
			S3Presigner:  s.S3Presigner,
			S3Bucket:     s.S3Bucket,
			S3Region:     s.S3Region,
		},
		Auth:         handlers.AuthHandler{Repo: s.UserRepo, RedisClient: s.RedisClient, GoEnv: s.GoEnv},
		Conversation: handlers.ConversationHandler{Repo: s.ConversationRepo},
		Message:      handlers.MessageHandler{Repo: s.MessageRepo},
		Friendship: handlers.FriendshipHandler{
			Repo:     s.FriendshipRepo,
			UserRepo: s.UserRepo,
		},
		Websocket: handlers.WebsocketHandler{
			Hub:              s.Hub,
			MessageService:   s.MessageService,
			ConversationRepo: s.ConversationRepo,
			GoEnv:            s.GoEnv,
		},
	}

	return routes.SetupRoutes(&h, s)
}
