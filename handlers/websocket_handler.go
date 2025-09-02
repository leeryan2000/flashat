package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/service"
	"github.com/leeryan2000/flashat/wire"
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
	uidStr := c.GetString("uid")

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

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		log.Println("❌ Failed to parse UID:", err)
		client.Cleanup()
		return
	}

	summaries, err := wh.ConversationRepo.GetSummary(c.Request.Context(), uid)
	if err != nil {
		log.Println("❌ Failed to retrieve conversations:", err)
		client.Cleanup()
		return
	}

	outEnv := &wire.Envelope{
		Type: "conv_summary",
		Body: summaries,
	}

	outEnvJson, err := json.Marshal(outEnv)
	if err != nil {
		log.Println("❌ Failed to marshal envelope:", err)
		client.Cleanup()
		return
	}

	client.Send <- outEnvJson

	// Add the user to the conversations in the Hub
	for _, conv := range summaries {
		convID := conv.ConversationID.String()
		if _, ok := wh.Hub.Conversations[convID]; !ok {
			// create a new conversation if the conversation didn't exists in the slice
			wh.Hub.Conversations[convID] = make(map[*server.Client]struct{})
		}
		wh.Hub.Conversations[convID][client] = struct{}{}
		client.Conversations[convID] = struct{}{}
	}

	client.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump(wh.MessageService.HandleEnvelope)
}
