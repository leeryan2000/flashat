package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
	"github.com/leeryan2000/flashat/server"
)

// ws://localhost:8080/api/websocket/ws
func InitializeWebsocketRoutes(router *gin.RouterGroup, h *handlers.Handlers, s *server.Server) {
	wh := h.Websocket

	websocketRoutes := router.Group("/websocket")
	websocketRoutes.Use(middleware.Authenticate(s)) // Ensure authentication middleware is applied

	websocketRoutes.GET("/ws", wh.ServeWs)
}
