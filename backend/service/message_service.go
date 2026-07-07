package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/wire"
)

type MsgHandlerFunc func(ctx context.Context, env *wire.Msg) error

type MessageService struct {
	Hub              Hub
	Publisher        Publisher
	MessageRepo      repo.MessageRepo
	ConversationRepo repo.ConversationRepo

	handlers  map[wire.WireType]MsgHandlerFunc
	observers []MessageObserver
}

func NewMessageService(hub Hub, msgRepo repo.MessageRepo, convRepo repo.ConversationRepo, publisher Publisher) *MessageService {
	s := &MessageService{
		Hub:              hub,
		MessageRepo:      msgRepo,
		ConversationRepo: convRepo,
		Publisher:        publisher,
		handlers:         make(map[wire.WireType]MsgHandlerFunc),
	}

	s.handlers[wire.Chat] = s.handleChat
	s.handlers[wire.Read] = s.handleRead

	return s
}

func (s *MessageService) RegisterObserver(obs MessageObserver) {
	s.observers = append(s.observers, obs)
}

func (s *MessageService) HandleEnvelope(ctx context.Context, env *wire.Msg) error {
	handler, ok := s.handlers[env.Type]
	if !ok {
		return nil // Unknown type, ignore
	}
	return handler(ctx, env)
}

// ---- Handlers ----

func (s *MessageService) handleRead(ctx context.Context, env *wire.Msg) error {
	conversationID, err := uuid.Parse(env.ConversationID)
	if err != nil {
		return err
	}

	fromUID, err := uuid.Parse(env.FromUID)
	if err != nil {
		return err
	}

	lastReadSeq := env.LastReadSeq
	slog.Info("handling read receipt", "uid", fromUID, "conv_id", conversationID, "seq", lastReadSeq)

	err = s.ConversationRepo.UpdateLastReadSeq(ctx, conversationID, fromUID, lastReadSeq)
	if err != nil {
		slog.Error("failed to update last read seq", "error", err)
		return err
	}

	return nil
}

func (s *MessageService) handleChat(ctx context.Context, env *wire.Msg) error {
	slog.Debug("handling chat message", "from_uid", env.FromUID)
	if env.ConversationID == "" {
		return errors.New("missing conversation id")
	}

	if _, err := uuid.Parse(env.ConversationID); err != nil {
		return errors.New("invalid conversation id")
	}
	if _, err := uuid.Parse(env.FromUID); err != nil {
		return errors.New("invalid from uid")
	}

	if env.ServerMsgID == "" {
		env.ServerMsgID = uuid.New().String()
	}

	return s.Publisher.Publish(env)
}

func (s *MessageService) ProcessChat(ctx context.Context, env *wire.Msg) error {
	conversationID, err := uuid.Parse(env.ConversationID)
	if err != nil {
		return err
	}

	fromUID, err := uuid.Parse(env.FromUID)
	if err != nil {
		return err
	}

	msgID, err := uuid.Parse(env.ServerMsgID)
	if err != nil {
		return err
	}

	// --- Save to DB first, Stops if this fails ---
	msg := &models.Message{
		ID:             msgID,
		ConversationID: conversationID,
		FromUID:        fromUID,
		Body:           env.Body,
	}
	if err := s.MessageRepo.SaveMessage(ctx, msg); err != nil {
		slog.Error("failed to save message", "error", err)
		return err
	}

	// --- Fetch participants ---
	participants, err := s.ConversationRepo.ListParticipantByID(ctx, conversationID)
	if err != nil {
		slog.Error("failed to list participants", "error", err)
		return err
	}
	uids := make([]uuid.UUID, 0, len(participants))
	for _, p := range participants {
		uids = append(uids, p.UID)
	}

	// --- Notify observers ---
	event := &MessageEvent{
		Ctx:          ctx,
		OriginalEnv:  env,
		SavedMsg:     msg,
		Participants: uids,
	}
	for _, obs := range s.observers {
		obs.OnMessageSaved(event)
	}

	return nil
}
