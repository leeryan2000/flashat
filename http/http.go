package http

import (
	"github.com/leeryan2000/flashat/handler"
	"github.com/leeryan2000/flashat/http/route"
	"github.com/leeryan2000/flashat/server"
)

func InitializeServer(s *server.Server) {
	handlers := handler.Handlers{
		User: handler.UserHandler{S: s},
	}

	http := route.SetupRoutes(&handlers)
	http.Run()
}
