package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/leeryan2000/flashat/repo"
)

type Client struct {
	UID   string
	Hub   *Hub
	Repo  repo.MessageRepo // Message repository for message persistence
	Conn  *websocket.Conn
	Send  chan []byte
	Rooms map[string]bool // rooms the client is subscribed to
}

func (c *Client) ReadPump() {
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
		var env Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			continue // ignore malformed
		}
		// stamp sender UID
		env.FromUID = c.UID

		// hand off to hub for routing
		c.Hub.Incoming <- &env
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
