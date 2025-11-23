package server

import (
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	Clients      map[*Client]struct{} // Changed from bool to struct{}
	ClientsByUID map[string]*Client

	Register   chan *Client
	Unregister chan *Client
}

// connect the user as long as the webpage is open
func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]struct{}),
		// ***** should change to map[string]map[*Client]struct{} if allowing multiple connections per user from different devices
		ClientsByUID: make(map[string]*Client),
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
	if client, ok := hub.ClientsByUID[uid.String()]; ok {
		client.Send <- payload
	}
}

func (hub *Hub) BroadcastToParticipant(uids []uuid.UUID, fromUID uuid.UUID, payload []byte) {
	// check if the participants are online and send the payload
	for _, uid := range uids {
		if client, ok := hub.ClientsByUID[uid.String()]; ok && uid != fromUID {
			log.Println("Broadcasting to UID:", uid)
			client.Send <- payload
		}
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
	}
}
