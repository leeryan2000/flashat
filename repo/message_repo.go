package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
)


type MessageRepo interface {
	SaveMessage(ctx context.Context, msg *models.Message) (int64, error)

	ListByConversationID(ctx context.Context, conversationID uuid.UUID, afterSeq int64, limit int) ([]models.Message, error)

	MarkDelivered(ctx context.Context, messageID uuid.UUID, uid uuid.UUID) error

	MarkRead(ctx context.Context, messageID uuid.UUID, uid uuid.UUID) error
}