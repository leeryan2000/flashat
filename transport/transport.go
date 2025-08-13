package transport

import (
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/repo/impl"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport/routes"
)

func InitializeServer(s *server.Server) {
	userRepo := &repo.GormUserRepo{DB: s.DB}
	conversationRepo := &repo.PgxConversationRepo{Pool: s.Pool}

	handlers := handlers.Handlers{
		User:         handlers.UserHandler{Repo: userRepo},
		Auth:         handlers.AuthHandler{Repo: userRepo},
		Conversation: handlers.ConversationHandler{Repo: conversationRepo},

		Websocket: handlers.WebsocketHandler{Hub: s.Hub},
	}

	transport := routes.SetupRoutes(&handlers)
	transport.Run()
}
