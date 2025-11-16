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

func (s *MessageService) HandleEnvelope(ctx context.Context, env *wire.MsgEnvelope) error {
	// Process the envelope based on its type
	switch env.Type {
	case wire.Chat:
		s.handleChat(ctx, env)
	case wire.Subscribe:
		s.handleSubscription(ctx, env)
	case wire.Join:
		// Handle join message
	case wire.Leave:
		// Handle leave message
	default:
		return nil // Unknown type, ignore
	}

	return nil // Processed successfully
}

func (s *MessageService) handleChat(ctx context.Context, env *wire.MsgEnvelope) error {
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
	ackEnv := &wire.MsgEnvelope{
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
	outEnv := &wire.MsgEnvelope{
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

	s.Hub.BroadcastToConversation(conversationID, outEnvJSON)

	return nil
}

// ***** continue here: implement subscription of user, sending wire from client to tell the server that which conversation the user should be subscribed to
func (s *MessageService) handleSubscription(ctx context.Context, env *wire.MsgEnvelope) {

}
