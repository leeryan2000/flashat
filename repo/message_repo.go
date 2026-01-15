package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/leeryan2000/flashat/models"
)

type MessageRepo interface {
	SaveMessage(ctx context.Context, msg *models.Message) error

	ListLatest(ctx context.Context, uid uuid.UUID, conversationID uuid.UUID, limit int) ([]models.Message, error)

	ListBefore(ctx context.Context, uid uuid.UUID, conversationID uuid.UUID, beforeSeq int64, limit int) ([]models.Message, error)

	ListAfter(ctx context.Context, conversationID uuid.UUID, afterSeq int64, limit int) ([]models.Message, error)

	// Could implement read count in future

	// Functions
	SaveMessageWithTx(ctx context.Context, tx pgx.Tx, msg *models.Message) error
}
