package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/leeryan2000/flashat/server"
)

type WebsocketHandler struct{ Hub *server.Hub }

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // For dev; restrict in prod!
	},
}

func (wh WebsocketHandler) ServeWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		// UID retrive from user
		uid := c.GetString("uid")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("❌ WebSocket upgrade failed:", err)
			return
		}
		log.Println("✅ WebSocket connection established")

		client := &server.Client{
			UID:   uid,
			Hub:   wh.Hub,
			Conn:  conn,
			Send:  make(chan []byte, 256),
			Rooms: make(map[string]bool),
		}

		client.Hub.Register <- client
		go client.WritePump()
		go client.ReadPump()
	}
}
