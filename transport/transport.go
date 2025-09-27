package transport

import (
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport/routes"
)

func InitializeServer(s *server.Server) {
	handlers := handlers.Handlers{
		User:         handlers.UserHandler{Repo: s.UserRepo},
		Auth:         handlers.AuthHandler{Repo: s.UserRepo},
		Conversation: handlers.ConversationHandler{Repo: s.ConversationRepo},
		Message:      handlers.MessageHandler{Repo: s.MessageRepo},

		Websocket: handlers.WebsocketHandler{
			Hub:              s.Hub,
			MessageService:   s.MessageService,
			ConversationRepo: s.ConversationRepo,
		},
	}

	transport := routes.SetupRoutes(&handlers, s)
	transport.Run()
}
