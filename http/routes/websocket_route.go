package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leeryan2000/flashat/handlers"
	"github.com/leeryan2000/flashat/middlewares"
)

func InitializeWebsocketRoutes(router *gin.RouterGroup, h *handlers.Handlers) {
	wh := h.Websocket
	
	websocketRoutes := router.Group("/websocket")
	websocketRoutes.Use(middlewares.Authenticate())
	

	websocketRoutes.GET("/ws", wh.ServeWs())
}
