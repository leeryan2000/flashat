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
	Hub  Hub
	Repo repo.MessageRepo
}

func (s *MessageService) HandleEnvelope(ctx context.Context, env *wire.Envelope) error {
	// Process the envelope based on its type
	switch env.Type {
	case wire.MsgChat:
		s.handleConversation(ctx, env)
	case wire.MsgJoin:
		// Handle join message
	case wire.MsgLeave:
		// Handle leave message
	default:
		return nil // Unknown type, ignore
	}

	return nil // Processed successfully
}

func (s *MessageService) handleConversation(ctx context.Context, env *wire.Envelope) error {
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

	err = s.Repo.SaveMessage(ctx, msg)

	ack := &wire.Ack{
		Status: "success",
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
	ackEnv := &wire.Envelope{
		Type:           wire.MsgAck,
		ConversationID: env.ConversationID,
		ClientMsgID:    env.ClientMsgID,
		ServerMsgID:    msg.ID.String(),
		FromUID:        env.FromUID,
		Seq:            msg.Seq,
		Ts:             msg.CreatedAt.UnixMilli(),
		Body:           ackJSON,
	}

	ackEnvJSON, err := json.Marshal(ackEnv)
	if err != nil {
		return err
	}

	// Send ack to tell the client if the message saved successfully
	s.Hub.SendToUID(fromUID, ackEnvJSON)

	// Pack the Envelope to broadcast to user
	outEnv := &wire.Envelope{
		Type:           wire.MsgChat,
		ConversationID: env.ConversationID,
		ClientMsgID:    env.ClientMsgID,
		ServerMsgID:    msg.ID.String(),
		FromUID:        env.FromUID,
		Seq:            msg.Seq, // client seq would be updated with server seq
		Ts:             msg.CreatedAt.UnixMilli(),
		Body:           env.Body, // ***** could create a wire message for the details of the content e.g. text, picture, is it reply to what message, or mentioned anyone
	}

	outEnvJSON, err := json.Marshal(outEnv)
	if err != nil {
		return err
	}

	s.Hub.BroadcastToConversation(conversationID, outEnvJSON)

	return nil
}
