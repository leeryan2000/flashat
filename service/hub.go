package service

import "github.com/google/uuid"

type Hub interface {
	SendToUID(uid uuid.UUID, payload []byte)
	BroadcastToParticipant(uids []uuid.UUID, fromUID uuid.UUID, payload []byte)
}
