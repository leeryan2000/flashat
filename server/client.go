package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/leeryan2000/flashat/wire"
)

type Client struct {
	UID           string
	Hub           *Hub
	Conn          *websocket.Conn
	Ctx           context.Context
	Cancel        context.CancelFunc
	Send          chan []byte
	Conversations map[string]struct{} // rooms the client is subscribed to
}

func (c *Client) Cleanup() {
	c.Hub.Unregister <- c // Unregister the client from the Hub
	c.Cancel()            // Cancel the context
	c.Conn.Close()        // Close the WebSocket connection
	close(c.Send)         // Close the send channel
}

type HandleEnvelope func(context.Context, *wire.MsgEnvelope) error

func (c *Client) ReadPump(handle HandleEnvelope) {
	defer c.Cleanup()

	for {
		// Read a full WS text frame
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// ***** delete line
		log.Println("Received message:", string(raw))
		// Expect JSON envelope from frontend
		var env wire.MsgEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			log.Println("❌ Failed to unmarshal envelope:", err)
			break // ignore malformed
		}

		if err = handle(c.Ctx, &env); err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	defer c.Cleanup()
	for env := range c.Send {
		log.Print("Sending message:", string(env))
		err := c.Conn.WriteMessage(websocket.TextMessage, env)
		if err != nil {
			break
		}
	}
}
