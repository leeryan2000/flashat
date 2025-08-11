package server

import (
	"encoding/json"
	"log"
)

type Hub struct {
	Clients      map[*Client]bool
	ClientsByUID map[string]*Client
	Rooms        map[string]map[*Client]bool

	Register   chan *Client
	Unregister chan *Client

	// Routed messages from clients
	Incoming chan *Envelope
}

func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*Client]bool),
		ClientsByUID: make(map[string]*Client),
		Rooms:        make(map[string]map[*Client]bool),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Incoming:     make(chan *Envelope, 1024),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case c := <-hub.Register:
			hub.Clients[c] = true
			hub.ClientsByUID[c.UID] = c

		case c := <-hub.Unregister:
			// see if c exists in the hub
			if _, ok := hub.Clients[c]; ok {
				// remove from global set
				delete(hub.Clients, c)
				delete(hub.ClientsByUID, c.UID)
				// remove from rooms
				for room := range c.Rooms {
					if set, ok := hub.Rooms[room]; ok {
						delete(set, c)
						if len(set) == 0 {
							delete(hub.Rooms, room)
						}
					}
				}
				close(c.Send)
			}

		case env := <-hub.Incoming:
			switch env.Type {
			case MsgJoin:
				hub.join(env.FromUID, env.RoomID)
			case MsgLeave:
				hub.leave(env.FromUID, env.RoomID)
			case MsgChat:
				hub.broadcastToRoom(env.RoomID, env)
			case MsgDirect:
				hub.direct(env.ToUID, env)
			default:
				// ignore or log unknown types
				log.Println("❌ Unknown message type:", env.Type)
			}
		}
	}
}

func (hub *Hub) join(uid, roomid string) {
	if uid == "" || roomid == "" {
		log.Println("❌ Join failed: missing UID or room ID")
		return
	}
	c := hub.ClientsByUID[uid]
	if c == nil {
		log.Println("❌ Join failed: no client found for UID", uid)
		return
	}
	if _, ok := hub.Rooms[roomid]; !ok {
		hub.Rooms[roomid] = make(map[*Client]bool)
	}
	hub.Rooms[roomid][c] = true
	if c.Rooms == nil {
		c.Rooms = make(map[string]bool)
	}
	c.Rooms[roomid] = true
}

func (hub *Hub) leave(uid, roomid string) {
	if uid == "" || roomid == "" {
		return
	}
	c := hub.ClientsByUID[uid]
	if c == nil {
		return
	}
	if set, ok := hub.Rooms[roomid]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(hub.Rooms, roomid)
		}
	}
	delete(c.Rooms, roomid)
}

func (hub *Hub) broadcastToRoom(roomid string, env *Envelope) {
	room := hub.Rooms[roomid]
	if room == nil {
		return
	}
	payload, _ := json.Marshal(env)
	for cli := range room {
		select {
		case cli.Send <- payload:
		default:
			// slow consumer: drop
			close(cli.Send)
			delete(room, cli)
			delete(hub.Clients, cli)
			delete(hub.ClientsByUID, cli.UID)
			for r := range cli.Rooms {
				if s, ok := hub.Rooms[r]; ok {
					delete(s, cli)
					if len(s) == 0 {
						delete(hub.Rooms, r)
					}
				}
			}
		}
	}
}

func (hub *Hub) direct(targetUID string, env *Envelope) {
	if targetUID == "" {
		return
	}
	// Find the target client by UID
	target := hub.ClientsByUID[targetUID]
	if target == nil {
		return
	}
	payload, _ := json.Marshal(env)
	select {
	case target.Send <- payload:
	default:
		// drop slow target
		close(target.Send)
		delete(hub.Clients, target)
		delete(hub.ClientsByUID, target.UID)
		for r := range target.Rooms {
			if s, ok := hub.Rooms[r]; ok {
				delete(s, target)
				if len(s) == 0 {
					delete(hub.Rooms, r)
				}
			}
		}
	}
}
