package server

import (
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	mu           sync.RWMutex
	Clients      map[*Client]struct{}
	ClientsByUID map[string]map[*Client]struct{} // supports multiple tabs/devices per user
}

// connect the user as long as the webpage is open
func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*Client]struct{}),
		ClientsByUID: make(map[string]map[*Client]struct{}),
	}
}

// Register registers c only if uid has fewer than maxConns existing connections
func (hub *Hub) Register(c *Client, maxConns int) bool {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if len(hub.ClientsByUID[c.UID]) >= maxConns {
		return false
	}
	hub.addClientLocked(c)
	return true
}

// Unregister removes c from the hub synchronously.
func (hub *Hub) Unregister(c *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.removeClientLocked(c)
}

func (hub *Hub) SendToUID(uid uuid.UUID, payload []byte) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	if clients, ok := hub.ClientsByUID[uid.String()]; ok {
		for client := range clients {
			client.Send <- payload
		}
	}
}

func (hub *Hub) BroadcastToParticipant(uids []uuid.UUID, fromUID uuid.UUID, payload []byte) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	for _, uid := range uids {
		if uid == fromUID {
			continue
		}
		if clients, ok := hub.ClientsByUID[uid.String()]; ok {
			for client := range clients {
				client.Send <- payload
			}
		}
	}
}

// Counts returns the current connection/user counts. Safe for concurrent use —
// use this instead of reading Clients/ClientsByUID directly from outside the Hub.
func (hub *Hub) Counts() (clients int, uniqueUsers int) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	return len(hub.Clients), len(hub.ClientsByUID)
}

// ConnectionCountForUID returns how many open connections a given user
// currently has. Safe for concurrent use — use this instead of reading
// ClientsByUID directly from outside the Hub.
func (hub *Hub) ConnectionCountForUID(uid string) int {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	return len(hub.ClientsByUID[uid])
}

// addClientLocked mutates hub state and must be called with hub.mu held.
func (hub *Hub) addClientLocked(c *Client) {
	hub.Clients[c] = struct{}{}
	if hub.ClientsByUID[c.UID] == nil {
		hub.ClientsByUID[c.UID] = make(map[*Client]struct{})
	}
	hub.ClientsByUID[c.UID][c] = struct{}{}
	slog.Info("client connected", "uid", c.UID)
}

// removeClientLocked mutates hub state and must be called with hub.mu held.
func (hub *Hub) removeClientLocked(c *Client) {
	if _, ok := hub.Clients[c]; ok {
		delete(hub.Clients, c)
		delete(hub.ClientsByUID[c.UID], c)
		if len(hub.ClientsByUID[c.UID]) == 0 {
			delete(hub.ClientsByUID, c.UID)
		}
		close(c.Send)
	}
}
