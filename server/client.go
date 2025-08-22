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

type HandleEnvelope func(context.Context, *wire.Envelope) error

func (c *Client) ReadPump(handle HandleEnvelope) {
	defer func() {
		c.Hub.Unregister <- c
		c.Cancel()
		c.Conn.Close()
	}()

	for {
		// Read a full WS text frame
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// ***** delete line
		log.Println("Received message:", string(raw))
		// Expect JSON envelope from frontend
		var env wire.Envelope
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
	defer func() {
		c.Hub.Unregister <- c
		c.Cancel()
		c.Conn.Close()
	}()

	for env := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, env)
		if err != nil {
			break
		}
	}
}
