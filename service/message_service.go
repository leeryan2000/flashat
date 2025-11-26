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

type MessageService struct {
	Hub              Hub
	MessageRepo      repo.MessageRepo
	ConversationRepo repo.ConversationRepo
}

func (s *MessageService) HandleEnvelope(ctx context.Context, env *wire.Msg) error {
	// Process the envelope based on its type
	switch env.Type {
	case wire.Chat:
		s.handleChat(ctx, env)
	case wire.Join:
		// Handle join message
	case wire.Leave:
		// Handle leave message
	default:
		return nil // Unknown type, ignore
	}

	return nil // Processed successfully
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

	ack := &wire.MsgAck{
		Status: "sent",
	}
	if err != nil {
		ack.Status = "failed"
		log.Println("❌ Failed to save message:", err)
		return err
	}
	// make ack json.rawMessage
	ackJSON, err := json.Marshal(ack)
	if err != nil {
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
		Body:           ackJSON,
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
		ClientMsgID:    env.ClientMsgID,
		FromUID:        env.FromUID,
		ServerMsgID:    msg.ID.String(),
		Seq:            msg.Seq, // client seq would be updated with server seq
		Ts:             msg.CreatedAt,
		// ***** could create wire message for the details of the content e.g. text, picture, is it reply to what message, or mentioned anyone
		Body: env.Body,
	}

	outEnvJSON, err := json.Marshal(outEnv)
	if err != nil {
		return err
	}
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

	return nil
}
