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
	MessageRepo      repo.MessageRepo
	ConversationRepo repo.ConversationRepo

	handlers map[wire.WireType]MsgHandlerFunc
}

func NewMessageService(hub Hub, msgRepo repo.MessageRepo, convRepo repo.ConversationRepo) *MessageService {
	s := &MessageService{
		Hub:              hub,
		MessageRepo:      msgRepo,
		ConversationRepo: convRepo,
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

	conversationID, err := uuid.Parse(env.ConversationID)
	if err != nil {
		return err
	}

	fromUID, err := uuid.Parse(env.FromUID)
	if err != nil {
		return err
	}
	// Create server generated msgID to replace the client generated
	msgID := uuid.New()
	msg := &models.Message{
		ID:             msgID,
		ConversationID: conversationID,
		FromUID:        fromUID,
		Body:           env.Body,
	}

	err = s.MessageRepo.SaveMessage(ctx, msg)

	if err != nil {
		log.Println("❌ Failed to save message:", err)
		return err
	}

	// Ack for the sender
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

	ackEnvJSON, err := json.Marshal(ackEnv)
	if err != nil {
		return err
	}

	// Send ack to tell the client if the message saved successfully
	s.Hub.SendToUID(fromUID, ackEnvJSON)

	// Pack the Envelope to broadcast to user
	outEnv := &wire.Msg{
		Type:           wire.Chat,
		ConversationID: env.ConversationID,
		ServerMsgID:    msg.ID.String(),
		FromUID:        env.FromUID,
		Seq:            msg.Seq, // client seq would be updated with server seq
		Ts:             msg.CreatedAt,
		// ***** could create wire message for the details of the content e.g. text, picture, is it reply to what message, or mentioned anyone
		Body: env.Body,
	}

	outEnvJSON, err := json.Marshal(outEnv)
	if err != nil {
		return err
	}

	// ***** cache: apply cache for the participantlist
	// retrieve participants from conversation using convID
	participants, err := s.ConversationRepo.ListParticipantByID(ctx, conversationID)
	if err != nil {
		return err
	}

	uids := make([]uuid.UUID, 0, len(participants))
	for _, p := range participants {
		uids = append(uids, p.UID)
	}

	s.Hub.BroadcastToParticipant(uids, fromUID, outEnvJSON)

	// Sync sender's other tabs/devices with the sent message
	s.Hub.SendToUID(fromUID, outEnvJSON)

	return nil
}
