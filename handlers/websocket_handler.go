package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/service"
	"github.com/leeryan2000/flashat/utils"
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
	token := c.Query("token")
	if token == "" {
		log.Println("❌ Missing token in WebSocket connection")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Validate the token
	claims, err := utils.ParseToken(token)
	if err != nil {
		log.Println("❌ Invalid token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	// UID retrive from user
	uidStr := claims.UID

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade failed:", err)
		return
	}
	log.Println("✅ WebSocket connection established")

	client := &server.Client{
		UID:           uidStr,
		Hub:           wh.Hub,
		Conn:          conn,
		Send:          make(chan []byte, 256),
		Conversations: make(map[string]struct{}),
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		log.Println("❌ Failed to parse UID:", err)
		conn.Close()
		return
	}

	conversations, err := wh.ConversationRepo.ListConversationByUID(c.Request.Context(), uid)
	if err != nil {
		log.Println("❌ Failed to retrieve conversations:", err)
		conn.Close()
		return
	}

	// Add the user to the conversations in the Hub
	for _, conv := range conversations {
		convID := conv.ID.String()
		if _, ok := wh.Hub.Conversations[convID]; !ok {
			// create a new conversation if the conversation didn't exists in the slice
			wh.Hub.Conversations[convID] = make(map[*server.Client]struct{})
		}
		wh.Hub.Conversations[convID][client] = struct{}{}
		client.Conversations[convID] = struct{}{}
	}

	client.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump(c.Request.Context(), wh.MessageService.HandleEnvelope)
}
