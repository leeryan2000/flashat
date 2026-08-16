package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat-posts/handlers"
	"github.com/leeryan2000/flashat-posts/server"
	"github.com/leeryan2000/flashat-posts/transport/routes"
)

// BuildRouter wires up all handlers and returns the router without
// starting it — main.go owns the server lifecycle, same split as
// backend/transport/transport.go.
func BuildRouter(s *server.Server) *gin.Engine {
	h := handlers.Handlers{
		Post: handlers.PostHandler{Repo: s.PostRepo, AuthClient: s.AuthClient},
	}

	return routes.SetupRoutes(&h, s.AuthClient)
}
