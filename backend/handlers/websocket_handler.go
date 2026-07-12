package handlers

import (
	"context"
	"log/slog"
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

func (wh WebsocketHandler) ServeWs(c *gin.Context) {
	const maxConnectionsPerUser = 5

	uidStr := c.GetString("uid")

	// Fast-path rejection before paying for the handshake. The authoritative
	// check happens atomically with registration below, since a check here
	// followed by a separate Register call would let two concurrent connect
	// requests both pass the check and both register, exceeding the cap.
	if wh.Hub.ConnectionCountForUID(uidStr) >= maxConnectionsPerUser {
		slog.Warn("connection limit reached", "uid", uidStr)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many open connections"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}
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

	if !wh.Hub.Register(client, maxConnectionsPerUser) {
		slog.Warn("connection limit reached", "uid", uidStr)
		cancel()
		conn.Close()
		return
	}

	go client.WritePump()
	go client.ReadPump(wh.MessageService.HandleEnvelope)
}
