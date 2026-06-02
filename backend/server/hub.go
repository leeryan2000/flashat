package server

import (
	"log/slog"

	"github.com/google/uuid"
)

type Hub struct {
	Clients      map[*Client]struct{}
	ClientsByUID map[string]map[*Client]struct{} // supports multiple tabs/devices per user

	Register   chan *Client
	Unregister chan *Client
}

// connect the user as long as the webpage is open
func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*Client]struct{}),
		ClientsByUID: make(map[string]map[*Client]struct{}),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case c := <-hub.Register:
			hub.addClient(c)

		case c := <-hub.Unregister:
			hub.removeClient(c)
		}
	}
}

func (hub *Hub) SendToUID(uid uuid.UUID, payload []byte) {
	if clients, ok := hub.ClientsByUID[uid.String()]; ok {
		for client := range clients {
			client.Send <- payload
		}
	}
}

func (hub *Hub) BroadcastToParticipant(uids []uuid.UUID, fromUID uuid.UUID, payload []byte) {
	for _, uid := range uids {
		if uid == fromUID {
			continue
		}
		if clients, ok := hub.ClientsByUID[uid.String()]; ok {
			for client := range clients {
				slog.Info("broadcasting to UID", "uid", uid)
				client.Send <- payload
			}
		}
	}
}

func (hub *Hub) addClient(c *Client) {
	hub.Clients[c] = struct{}{}
	if hub.ClientsByUID[c.UID] == nil {
		hub.ClientsByUID[c.UID] = make(map[*Client]struct{})
	}
	hub.ClientsByUID[c.UID][c] = struct{}{}
	slog.Info("client registered", "uid", c.UID)
}

func (hub *Hub) removeClient(c *Client) {
	if _, ok := hub.Clients[c]; ok {
		delete(hub.Clients, c)
		delete(hub.ClientsByUID[c.UID], c)
		if len(hub.ClientsByUID[c.UID]) == 0 {
			delete(hub.ClientsByUID, c.UID)
		}
		close(c.Send)
	}
}
