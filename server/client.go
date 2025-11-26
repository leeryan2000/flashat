package server

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
}

const (
	writeWait = 10 * time.Second // per-write deadline (text, ping, close)
	pongWait  = 60 * time.Second // must receive PONG within this after last read/ping
	pingEvery = 25 * time.Second // send PINGs this often; must be < pongWait and < infra idle timeout
)

type HandleEnvelope func(context.Context, *wire.Msg) error

func (c *Client) ReadPump(handle HandleEnvelope) {
	defer c.Cleanup()

	// add pong handler
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Println("Pong received from UID:", c.UID)
		return nil
	})

	for {
		// Read a full WS text frame
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// ***** delete line
		log.Println("Received message:", string(raw))
		// Expect JSON envelope from frontend
		var env wire.Msg
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
	ticker := time.NewTicker(pingEvery)
	defer func() {
		c.Cleanup()
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "hub closed"),
					time.Now().Add(writeWait),
				)
				return
			}
			log.Print("Sending message:", string(msg))
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			log.Println("Sending ping to UID:", c.UID)
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("❌ Ping failed, closing connection:", err)
				c.Cleanup()
				return
			}
		case <-c.Ctx.Done():
			return
		}
	}

}
