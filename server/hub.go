package server

import "github.com/google/uuid"

type Hub struct {
	Clients       map[*Client]bool
	ClientsByUID  map[string]*Client
	Conversations map[string]map[*Client]bool

	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:       make(map[*Client]bool),
		ClientsByUID:  make(map[string]*Client),
		Conversations: make(map[string]map[*Client]bool),
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

}

func (hub *Hub) BroadcastToConversation(convID uuid.UUID, payload []byte) {

}

func (hub *Hub) addClient(c *Client) {
	hub.Clients[c] = true
	hub.ClientsByUID[c.UID] = c
}

func (hub *Hub) removeClient(c *Client) {
	if _, ok := hub.Clients[c]; ok {
		delete(hub.Clients, c)
		delete(hub.ClientsByUID, c.UID)
		for conv := range c.Conversations {
			if set, ok := hub.Conversations[conv]; ok {
				delete(set, c)
				if len(set) == 0 {
					delete(hub.Conversations, conv)
				}
			}
		}
		close(c.Send)
	}
}
