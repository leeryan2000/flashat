package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/mocks"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/wire"
	"github.com/stretchr/testify/mock"
)

// Compile-time checks that the mocks satisfy these package-local interfaces.
// Must live here (package service), not in package mocks, to avoid mocks
// importing service — see plan for the import-cycle reasoning.
var _ Hub = (*mocks.HubMock)(nil)
var _ Publisher = (*mocks.PublisherMock)(nil)

// fakeObserver is a plain recording spy, not a testify mock: OnMessageSaved
// is a single void method with nothing to stub, so a mock buys nothing over
// just recording what it received.
type fakeObserver struct {
	events []*MessageEvent
}

func (f *fakeObserver) OnMessageSaved(event *MessageEvent) {
	f.events = append(f.events, event)
}

func newTestMessageService() (*MessageService, *mocks.HubMock, *mocks.PublisherMock, *mocks.MessageRepoMock, *mocks.ConversationRepoMock) {
	hub := new(mocks.HubMock)
	pub := new(mocks.PublisherMock)
	msgRepo := new(mocks.MessageRepoMock)
	convRepo := new(mocks.ConversationRepoMock)
	svc := NewMessageService(hub, msgRepo, convRepo, pub)
	return svc, hub, pub, msgRepo, convRepo
}

func TestHandleChat(t *testing.T) {
	ctx := context.Background()

	t.Run("missing conversation id", func(t *testing.T) {
		svc, _, pub, _, _ := newTestMessageService()
		env := &wire.Msg{ConversationID: "", FromUID: uuid.New().String()}

		err := svc.handleChat(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		pub.AssertNotCalled(t, "Publish", mock.Anything)
	})

	t.Run("invalid conversation id", func(t *testing.T) {
		svc, _, pub, _, _ := newTestMessageService()
		env := &wire.Msg{ConversationID: "not-a-uuid", FromUID: uuid.New().String()}

		err := svc.handleChat(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		pub.AssertNotCalled(t, "Publish", mock.Anything)
	})

	t.Run("invalid from uid", func(t *testing.T) {
		svc, _, pub, _, _ := newTestMessageService()
		env := &wire.Msg{ConversationID: uuid.New().String(), FromUID: "not-a-uuid"}

		err := svc.handleChat(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		pub.AssertNotCalled(t, "Publish", mock.Anything)
	})

	t.Run("auto-generates ServerMsgID when empty and publishes it", func(t *testing.T) {
		svc, _, pub, _, _ := newTestMessageService()
		env := &wire.Msg{
			ConversationID: uuid.New().String(),
			FromUID:        uuid.New().String(),
			ServerMsgID:    "",
		}
		pub.On("Publish", mock.AnythingOfType("*wire.Msg")).Return(nil)

		err := svc.handleChat(ctx, env)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if env.ServerMsgID == "" {
			t.Fatal("expected ServerMsgID to be auto-generated, still empty")
		}
		if _, parseErr := uuid.Parse(env.ServerMsgID); parseErr != nil {
			t.Fatalf("auto-generated ServerMsgID is not a valid uuid: %v", parseErr)
		}
		pub.AssertCalled(t, "Publish", env)
	})

	t.Run("preserves already-set ServerMsgID", func(t *testing.T) {
		svc, _, pub, _, _ := newTestMessageService()
		env := &wire.Msg{
			ConversationID: uuid.New().String(),
			FromUID:        uuid.New().String(),
			ServerMsgID:    "existing-id",
		}
		pub.On("Publish", env).Return(nil)

		err := svc.handleChat(ctx, env)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if env.ServerMsgID != "existing-id" {
			t.Fatalf("ServerMsgID changed, got %q, want %q", env.ServerMsgID, "existing-id")
		}
	})

	t.Run("propagates Publisher.Publish error", func(t *testing.T) {
		svc, _, pub, _, _ := newTestMessageService()
		env := &wire.Msg{
			ConversationID: uuid.New().String(),
			FromUID:        uuid.New().String(),
		}
		pub.On("Publish", mock.AnythingOfType("*wire.Msg")).Return(errors.New("publish boom"))

		err := svc.handleChat(ctx, env)

		if err == nil || err.Error() != "publish boom" {
			t.Fatalf("expected publish error to propagate, got %v", err)
		}
	})
}

func TestHandleRead(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid conversation id", func(t *testing.T) {
		svc, _, _, _, convRepo := newTestMessageService()
		env := &wire.Msg{ConversationID: "bad", FromUID: uuid.New().String()}

		err := svc.handleRead(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		convRepo.AssertNotCalled(t, "UpdateLastReadSeq", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("invalid from uid", func(t *testing.T) {
		svc, _, _, _, convRepo := newTestMessageService()
		env := &wire.Msg{ConversationID: uuid.New().String(), FromUID: "bad"}

		err := svc.handleRead(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		convRepo.AssertNotCalled(t, "UpdateLastReadSeq", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("calls UpdateLastReadSeq with correct args on success", func(t *testing.T) {
		svc, _, _, _, convRepo := newTestMessageService()
		convID := uuid.New()
		fromUID := uuid.New()
		env := &wire.Msg{ConversationID: convID.String(), FromUID: fromUID.String(), LastReadSeq: 42}
		convRepo.On("UpdateLastReadSeq", ctx, convID, fromUID, int64(42)).Return(nil)

		err := svc.handleRead(ctx, env)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		convRepo.AssertExpectations(t)
	})

	t.Run("propagates UpdateLastReadSeq error", func(t *testing.T) {
		svc, _, _, _, convRepo := newTestMessageService()
		convID := uuid.New()
		fromUID := uuid.New()
		env := &wire.Msg{ConversationID: convID.String(), FromUID: fromUID.String(), LastReadSeq: 1}
		convRepo.On("UpdateLastReadSeq", ctx, convID, fromUID, int64(1)).Return(errors.New("db down"))

		err := svc.handleRead(ctx, env)

		if err == nil || err.Error() != "db down" {
			t.Fatalf("expected db error to propagate, got %v", err)
		}
	})
}

func TestProcessChat(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid conversation id", func(t *testing.T) {
		svc, _, _, msgRepo, _ := newTestMessageService()
		env := &wire.Msg{ConversationID: "bad", FromUID: uuid.New().String(), ServerMsgID: uuid.New().String()}

		err := svc.ProcessChat(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		msgRepo.AssertNotCalled(t, "SaveMessage", mock.Anything, mock.Anything)
	})

	t.Run("invalid from uid", func(t *testing.T) {
		svc, _, _, msgRepo, _ := newTestMessageService()
		env := &wire.Msg{ConversationID: uuid.New().String(), FromUID: "bad", ServerMsgID: uuid.New().String()}

		err := svc.ProcessChat(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		msgRepo.AssertNotCalled(t, "SaveMessage", mock.Anything, mock.Anything)
	})

	t.Run("invalid server msg id", func(t *testing.T) {
		svc, _, _, msgRepo, _ := newTestMessageService()
		env := &wire.Msg{ConversationID: uuid.New().String(), FromUID: uuid.New().String(), ServerMsgID: "bad"}

		err := svc.ProcessChat(ctx, env)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		msgRepo.AssertNotCalled(t, "SaveMessage", mock.Anything, mock.Anything)
	})

	t.Run("SaveMessage failure short-circuits before ListParticipantByID and before observers fire", func(t *testing.T) {
		svc, _, _, msgRepo, convRepo := newTestMessageService()
		env := &wire.Msg{ConversationID: uuid.New().String(), FromUID: uuid.New().String(), ServerMsgID: uuid.New().String()}
		msgRepo.On("SaveMessage", ctx, mock.AnythingOfType("*models.Message")).Return(errors.New("save failed"))
		observer := &fakeObserver{}
		svc.RegisterObserver(observer)

		err := svc.ProcessChat(ctx, env)

		if err == nil || err.Error() != "save failed" {
			t.Fatalf("expected save error to propagate, got %v", err)
		}
		convRepo.AssertNotCalled(t, "ListParticipantByID", mock.Anything, mock.Anything)
		if len(observer.events) != 0 {
			t.Fatalf("expected no observer events, got %d", len(observer.events))
		}
	})

	t.Run("ListParticipantByID failure propagates and observers do not fire", func(t *testing.T) {
		svc, _, _, msgRepo, convRepo := newTestMessageService()
		convID := uuid.New()
		env := &wire.Msg{ConversationID: convID.String(), FromUID: uuid.New().String(), ServerMsgID: uuid.New().String()}
		msgRepo.On("SaveMessage", ctx, mock.AnythingOfType("*models.Message")).Return(nil)
		convRepo.On("ListParticipantByID", ctx, convID).Return(nil, errors.New("lookup failed"))
		observer := &fakeObserver{}
		svc.RegisterObserver(observer)

		err := svc.ProcessChat(ctx, env)

		if err == nil || err.Error() != "lookup failed" {
			t.Fatalf("expected lookup error to propagate, got %v", err)
		}
		if len(observer.events) != 0 {
			t.Fatalf("expected no observer events, got %d", len(observer.events))
		}
	})

	t.Run("success path notifies every registered observer exactly once", func(t *testing.T) {
		svc, _, _, msgRepo, convRepo := newTestMessageService()

		convID, fromUID, msgID := uuid.New(), uuid.New(), uuid.New()
		body := json.RawMessage(`"hello"`)
		env := &wire.Msg{
			ConversationID: convID.String(),
			FromUID:        fromUID.String(),
			ServerMsgID:    msgID.String(),
			Body:           body,
		}
		msgRepo.On("SaveMessage", ctx, mock.AnythingOfType("*models.Message")).Return(nil)

		p1, p2 := uuid.New(), uuid.New()
		participants := []*models.ConversationParticipant{
			{ConversationID: convID, UID: p1},
			{ConversationID: convID, UID: p2},
		}
		convRepo.On("ListParticipantByID", ctx, convID).Return(participants, nil)

		obs1 := &fakeObserver{}
		obs2 := &fakeObserver{}
		svc.RegisterObserver(obs1)
		svc.RegisterObserver(obs2)

		err := svc.ProcessChat(ctx, env)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for name, obs := range map[string]*fakeObserver{"obs1": obs1, "obs2": obs2} {
			if len(obs.events) != 1 {
				t.Fatalf("%s: expected exactly 1 event, got %d", name, len(obs.events))
			}
			event := obs.events[0]
			if event.SavedMsg.ID != msgID || event.SavedMsg.ConversationID != convID || event.SavedMsg.FromUID != fromUID {
				t.Fatalf("%s: SavedMsg mismatch: %+v", name, event.SavedMsg)
			}
			if string(event.SavedMsg.Body) != string(body) {
				t.Fatalf("%s: SavedMsg.Body mismatch: got %s, want %s", name, event.SavedMsg.Body, body)
			}
			if len(event.Participants) != 2 || event.Participants[0] != p1 || event.Participants[1] != p2 {
				t.Fatalf("%s: Participants mismatch: %+v", name, event.Participants)
			}
			if event.OriginalEnv != env {
				t.Fatalf("%s: OriginalEnv is not the same pointer passed in", name)
			}
			if event.Ctx != ctx {
				t.Fatalf("%s: Ctx is not the same context passed in", name)
			}
		}
	})
}

func TestHandleEnvelope(t *testing.T) {
	ctx := context.Background()

	t.Run("unknown wire type returns nil without touching any dependency", func(t *testing.T) {
		svc, _, pub, msgRepo, convRepo := newTestMessageService()
		env := &wire.Msg{Type: wire.System, ConversationID: "anything"}

		err := svc.HandleEnvelope(ctx, env)

		if err != nil {
			t.Fatalf("expected nil for unknown wire type, got %v", err)
		}
		pub.AssertNotCalled(t, "Publish", mock.Anything)
		msgRepo.AssertNotCalled(t, "SaveMessage", mock.Anything, mock.Anything)
		convRepo.AssertNotCalled(t, "UpdateLastReadSeq", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		convRepo.AssertNotCalled(t, "ListParticipantByID", mock.Anything, mock.Anything)
	})
}
