package repo

import (
	"context"

	"github.com/google/uuid"
	// "github.com/leeryan2000/flashat/models"
)

type ConversationRepo interface {
	CreateGroupConversation(ctx context.Context, creatorUID uuid.UUID, participantsUID []uuid.UUID, groupName string) error
	// For direct chats: get existing or create with a canonical direct_key.
	// GetOrCreateDirectConversation(ctx context.Context, uidA, uidB uuid.UUID) (*models.Conversation, error)

	// ListByUID(ctx context.Context, uid uuid.UUID, page, pageSize int) ([]*models.Conversation, error)

	// GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)

	// AddParticipant(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) error

	// RemoveParticipant(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) error

	// // Read-state
	// UpdateLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID, seq int64) error
	// GetLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) (int64, error)

	// // Next Sequence
	// NextSeq(ctx context.Context, conversationID uuid.UUID) (int64, error)
}
