package service

import "github.com/google/uuid"

type Hub interface {
	SendToUID(uid uuid.UUID, payload []byte)
	BroadcastToConversation(conversationID uuid.UUID, payload []byte)
}
