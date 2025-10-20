package server

import (
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	Clients       map[*Client]struct{} // Changed from bool to struct{}
	ClientsByUID  map[string]*Client
	Conversations map[string]map[*Client]struct{} // Changed from bool to struct{}

	Register   chan *Client
	Unregister chan *Client
}

// connect the user as long as the webpage is open
func NewHub() *Hub {
	return &Hub{
		Clients:       make(map[*Client]struct{}),
		ClientsByUID:  make(map[string]*Client),
		Conversations: make(map[string]map[*Client]struct{}),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
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
	if client, ok := hub.ClientsByUID[uid.String()]; ok {
		client.Send <- payload
	}
}

func (hub *Hub) BroadcastToConversation(convID uuid.UUID, payload []byte) {
	if clients, ok := hub.Conversations[convID.String()]; ok {
		for client := range clients {
			client.Send <- payload
		}
	}
}

// ***** implement subscribeToConversation
func (hub *Hub) subscribeToConversation(uid uuid.UUID, convID string) {
	// create the conversation map if it doesn't exist
	if _, ok := hub.Conversations[convID]; !ok {
		hub.Conversations[convID] = make(map[*Client]struct{})
	}
	// add the client to the conversation if the client exists
	if client, ok := hub.ClientsByUID[uid.String()]; ok {
		hub.Conversations[convID][client] = struct{}{}
		client.Conversations[convID] = struct{}{}
	}
}

func (hub *Hub) addClient(c *Client) {
	hub.Clients[c] = struct{}{}
	hub.ClientsByUID[c.UID] = c
	log.Println("Client registered:", c.UID)
}

func (hub *Hub) removeClient(c *Client) {
	if _, ok := hub.Clients[c]; ok {
		delete(hub.Clients, c)
		delete(hub.ClientsByUID, c.UID)
		close(c.Send)
		for conv := range c.Conversations {
			if set, ok := hub.Conversations[conv]; ok {
				delete(set, c)
				if len(set) == 0 {
					delete(hub.Conversations, conv)
				}
			}
		}
	}
}
