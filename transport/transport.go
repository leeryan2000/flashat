package transport

import (
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport/routes"
)

func InitializeServer(s *server.Server) {
	userRepo := &repo.GormUserRepo{DB: s.DB}
	conversationRepo := &repo.PgxConversationRepo{Pool: s.Pool}

	handlers := handlers.Handlers{
		User:         handlers.UserHandler{Repo: userRepo},
		Auth:         handlers.AuthHandler{Repo: userRepo},
		Websocket:    handlers.WebsocketHandler{Hub: s.Hub},
		Conversation: handlers.ConversationHandler{Repo: conversationRepo},
	}

	transport := routes.SetupRoutes(&handlers)
	transport.Run()
}
