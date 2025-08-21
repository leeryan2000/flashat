package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middleware"
)

func InitializeWebsocketRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	wh := h.Websocket

	websocketRoutes := router.Group("/websocket")
	router.Use(middleware.Authenticate()) // Ensure authentication middleware is applied

	websocketRoutes.GET("/ws", wh.ServeWs)
}
