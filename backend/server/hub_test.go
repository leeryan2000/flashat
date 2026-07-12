package server

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
)

// newTestClient builds a Client with no real websocket.Conn — Hub's
// register/unregister/send paths never touch Conn, so a bare Send
// channel is enough to exercise them.
func newTestClient(uid string) *Client {
	return &Client{UID: uid, Send: make(chan []byte, 256)}
}

// noCap is passed to TryRegister in tests that aren't exercising the
// connection cap and just want a plain, always-succeeding registration.
const noCap = math.MaxInt

// drain simulates WritePump: it keeps consuming Send so SendToUID /
// BroadcastToParticipant never block on a full channel while holding
// the hub's RLock. It exits on its own once removeClient closes Send.
func drain(c *Client) {
	go func() {
		for range c.Send {
		}
	}()
}

// TestHub_ConcurrentRegisterSendUnregister hammers the Hub the way
// production traffic does: clients connecting/disconnecting on one set
// of goroutines while messages are broadcast to them on another.
//
// Run with the race detector — this is what actually catches the bug
// class fixed in 22f4142 / 3ee4a50:
//
//	go test ./server/... -race -run TestHub -count=20
func TestHub_ConcurrentRegisterSendUnregister(t *testing.T) {
	hub := NewHub()

	const (
		numUsers          = 10
		registrarsPerUser = 5
		itersPerRegistrar = 200
		senders           = 20
		msgsPerSender     = 500
	)

	uids := make([]uuid.UUID, numUsers)
	for i := range uids {
		uids[i] = uuid.New()
	}

	var wg sync.WaitGroup

	// Registrars: continuously connect and disconnect clients for each
	// user, concurrently with the senders below.
	for _, uid := range uids {
		for j := 0; j < registrarsPerUser; j++ {
			wg.Add(1)
			go func(uid uuid.UUID) {
				defer wg.Done()
				for i := 0; i < itersPerRegistrar; i++ {
					c := newTestClient(uid.String())
					drain(c)
					hub.Register(c, noCap)
					hub.Unregister(c)
				}
			}(uid)
		}
	}

	// Senders: hammer SendToUID / BroadcastToParticipant / the read-only
	// count helpers while clients are being registered and unregistered.
	for i := 0; i < senders; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload := []byte(fmt.Sprintf("payload-%d", i))
			for j := 0; j < msgsPerSender; j++ {
				uid := uids[j%len(uids)]
				hub.SendToUID(uid, payload)
				hub.BroadcastToParticipant(uids, uuid.Nil, payload)
				hub.Counts()
				hub.ConnectionCountForUID(uid.String())
			}
		}(i)
	}

	wg.Wait()

	// Every registrar unregistered everything it registered, so the hub
	// should be back to empty — if it isn't, addClient/removeClient are
	// racing each other and losing updates.
	clients, users := hub.Counts()
	if clients != 0 || users != 0 {
		t.Fatalf("hub not empty after all clients unregistered: clients=%d users=%d", clients, users)
	}
}

// TestHub_MultipleConnectionsPerUser checks the "multiple tabs/devices
// per user" bookkeeping in ClientsByUID: registering N clients under the
// same UID and unregistering one at a time should shrink the set by
// exactly one each time, not drop the whole user early.
func TestHub_MultipleConnectionsPerUser(t *testing.T) {
	hub := NewHub()

	uid := uuid.New()
	const numConns = 4

	clients := make([]*Client, numConns)
	for i := range clients {
		c := newTestClient(uid.String())
		drain(c)
		clients[i] = c
		hub.Register(c, noCap)
	}

	if got := hub.ConnectionCountForUID(uid.String()); got != numConns {
		t.Fatalf("ConnectionCountForUID = %d, want %d", got, numConns)
	}

	for i, c := range clients {
		hub.Unregister(c)
		want := numConns - i - 1
		if got := hub.ConnectionCountForUID(uid.String()); got != want {
			t.Fatalf("after unregistering client %d: ConnectionCountForUID = %d, want %d", i, got, want)
		}
	}

	_, users := hub.Counts()
	if users != 0 {
		t.Fatalf("expected 0 unique users after all connections closed, got %d", users)
	}
}

// TestHub_Register_EnforcesCapUnderConcurrency is the regression test for
// the race in websocket_handler.go: many simultaneous connection attempts
// for the same user must never leave more than maxConns registered, even
// though each attempt only decides based on hub state observed at call
// time. A separate ConnectionCountForUID-then-insert (the old code path)
// fails this every time; Register's single-lock check-and-insert does not.
func TestHub_Register_EnforcesCapUnderConcurrency(t *testing.T) {
	hub := NewHub()
	uid := uuid.New()
	const (
		maxConns = 5
		attempts = 50
	)

	var wg sync.WaitGroup
	var accepted int32
	var mu sync.Mutex
	var accepteds []*Client

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := newTestClient(uid.String())
			drain(c)
			if hub.Register(c, maxConns) {
				atomic.AddInt32(&accepted, 1)
				mu.Lock()
				accepteds = append(accepteds, c)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if accepted != maxConns {
		t.Fatalf("accepted %d registrations, want exactly %d (cap=%d, attempts=%d)", accepted, maxConns, maxConns, attempts)
	}
	if got := hub.ConnectionCountForUID(uid.String()); got != maxConns {
		t.Fatalf("ConnectionCountForUID = %d, want %d", got, maxConns)
	}

	for _, c := range accepteds {
		hub.Unregister(c)
	}
}
