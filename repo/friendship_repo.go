package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/wire"
)

type FriendshipRepo interface {
	// Methods for managing friendships would go here
	RequestFriendship(ctx context.Context, requesterUID uuid.UUID, email string) error

	AcceptFriendship(ctx context.Context, conv *models.Conversation, msg *models.Message, requesterUID, receiverUID uuid.UUID) error

	DeleteFriendship(ctx context.Context, uid1, uid2 uuid.UUID) error

	RemovePendingFriendship(ctx context.Context, requesterUID, receiverUID uuid.UUID) error

	BlockUser(ctx context.Context, uid1, uid2 uuid.UUID) error

	ListFriendships(ctx context.Context, uid uuid.UUID) ([]wire.Friendship, error)

	GetFriendshipStatus(ctx context.Context, uid1, uid2 uuid.UUID) (string, error)
}
