package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/wire"
)

// MessageEvent carries everything observers need after a message is saved
type MessageEvent struct {
	Ctx          context.Context
	OriginalEnv  *wire.Msg
	SavedMsg     *models.Message
	Participants []uuid.UUID
}

// MessageObserver is notified after a message is successfully saved
type MessageObserver interface {
	OnMessageSaved(event *MessageEvent)
}
