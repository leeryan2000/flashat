package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

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

	handlers map[wire.WireType]MsgHandlerFunc
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
	log.Println("Handling read receipt from UID:", fromUID, "for conversation:", conversationID, "up to seq:", lastReadSeq)

	err = s.ConversationRepo.UpdateLastReadSeq(ctx, conversationID, fromUID, lastReadSeq)
	if err != nil {
		log.Println("❌ Failed to update last read seq:", err)
		return err
	}

	return nil
}

func (s *MessageService) handleChat(ctx context.Context, env *wire.Msg) error {
	log.Println("Handling chat message from UID:", env.FromUID)
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

	// --- Save to DB ---
	msg := &models.Message{
		ID:             msgID,
		ConversationID: conversationID,
		FromUID:        fromUID,
		Body:           env.Body,
	}

	if err := s.MessageRepo.SaveMessage(ctx, msg); err != nil {
		log.Println("❌ Failed to save message:", err)
		return err
	}

	// --- ACK the sender with server seq and id ---
	ackEnv := &wire.Msg{
		Type:           wire.Ack,
		ConversationID: env.ConversationID,
		ClientMsgID:    env.ClientMsgID,
		ServerMsgID:    msg.ID.String(),
		FromUID:        env.FromUID,
		Seq:            msg.Seq,
		Ts:             msg.CreatedAt,
		Body:           env.Body,
	}
	ackJSON, err := json.Marshal(ackEnv)
	if err != nil {
		return err
	}
	s.Hub.SendToUID(fromUID, ackJSON)

	// --- Broadcast to other participants ---
	outEnv := &wire.Msg{
		Type:           wire.Chat,
		ConversationID: env.ConversationID,
		ServerMsgID:    msg.ID.String(),
		FromUID:        env.FromUID,
		Seq:            msg.Seq,
		Ts:             msg.CreatedAt,
		Body:           env.Body,
	}
	outJSON, err := json.Marshal(outEnv)
	if err != nil {
		return err
	}

	participants, err := s.ConversationRepo.ListParticipantByID(ctx, conversationID)
	if err != nil {
		log.Println("❌ Failed to list participants:", err)
		return err
	}

	uids := make([]uuid.UUID, 0, len(participants))
	for _, p := range participants {
		uids = append(uids, p.UID)
	}

	s.Hub.BroadcastToParticipant(uids, fromUID, outJSON)
	s.Hub.SendToUID(fromUID, outJSON) // sync sender's other tabs/devices

	return nil
}
