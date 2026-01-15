package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/wire"
)

type ConversationRepo interface {
	CreateGroupConversation(ctx context.Context, conv *models.Conversation, creatorUID uuid.UUID, participantsUID []uuid.UUID) error
	// For direct chats: get existing or create with a canonical direct_key.
	CreateDirectConversation(ctx context.Context, conv *models.Conversation, uidA, uidB uuid.UUID) error

	ListConversationByUID(ctx context.Context, uid uuid.UUID) ([]*models.Conversation, error)
	GetConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error)
	ListParticipantByID(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationParticipant, error)

	AddParticipant(ctx context.Context, conversationID uuid.UUID, participantUID uuid.UUID) error
	ModifyParticipant(ctx context.Context, conversationID uuid.UUID, participantUID uuid.UUID, role string) error
	RemoveParticipant(ctx context.Context, conversationID uuid.UUID, participantUID uuid.UUID) error

	// Read-state
	UpdateLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID, seq int64) error
	GetLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) (int64, error)

	// Load conversations
	GetSummary(ctx context.Context, uid uuid.UUID) ([]*wire.Conversation, error)

	// Functions
	CreateDirectConversationWithTx(ctx context.Context, conv *models.Conversation, tx pgx.Tx, uid1, uid2 uuid.UUID) error
}
