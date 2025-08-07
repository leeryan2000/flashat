package http

import (
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/http/routes"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/server"
)

func InitializeServer(s *server.Server) {
	userRepo := &repo.GormUserRepo{DB: s.DB}

	handlers := handlers.Handlers{
		User:      handlers.UserHandler{Repo: userRepo},
		Auth:      handlers.AuthHandler{Repo: userRepo},
		Websocket: handlers.WebsocketHandler{Hub: s.Hub},
	}

	http := routes.SetupRoutes(&handlers)
	http.Run()
}
