package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/service"
)

type WebsocketHandler struct {
	Hub              *server.Hub
	MessageService   *service.MessageService
	ConversationRepo repo.ConversationRepo
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // For dev; restrict in prod!
	},
}

// ***** make it more bulletproof
func (wh WebsocketHandler) ServeWs(c *gin.Context) {
	uidStr := c.GetString("uid")

	// see if the user are connected already
	if _, ok := wh.Hub.ClientsByUID[uidStr]; ok {
		log.Println("User already connected:", uidStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already connected"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade failed:", err)
		return
	}
	log.Println("✅ WebSocket connection established")

	// Create a context that lives with the client lifecycle
	ctx, cancel := context.WithCancel(context.Background())

	client := &server.Client{
		UID:           uidStr,
		Hub:           wh.Hub,
		Conn:          conn,
		Ctx:           ctx,
		Cancel:        cancel,
		Send:          make(chan []byte, 256),
		Conversations: make(map[string]struct{}),
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump(wh.MessageService.HandleEnvelope)
}
