package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/leeryan2000/flashat/wire"
)

type Client struct {
	UID           string
	Hub           *Hub
	Conn          *websocket.Conn
	Send          chan []byte
	Conversations map[string]bool // rooms the client is subscribed to
}

type HandleEnvelope func(*wire.Envelope) error

func (c *Client) ReadPump(handle HandleEnvelope) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		// Read a full WS text frame
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// Expect JSON envelope from frontend
		var env wire.Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			continue // ignore malformed
		}

		handle(&env)
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
