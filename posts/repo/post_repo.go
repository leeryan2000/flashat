package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat-posts/models"
)

type PostRepo interface {
	CreatePost(ctx context.Context, authorUID uuid.UUID, body string) (*models.Post, error)

	GetPost(ctx context.Context, requesterUID, postID uuid.UUID) (*models.Post, error)

	// ListFeedLatest/ListFeedBefore return posts authored by anyone in
	// visibleUIDs (the requester's accepted friends + themself, per the
	// handler), newest first — feed semantics, unlike message pagination
	// which returns ascending (oldest-first, chat-thread reading order).
	ListFeedLatest(ctx context.Context, requesterUID uuid.UUID, visibleUIDs []uuid.UUID, limit int) ([]models.Post, error)

	ListFeedBefore(ctx context.Context, requesterUID uuid.UUID, visibleUIDs []uuid.UUID, beforeSeq int64, limit int) ([]models.Post, error)

	// ToggleLike flips the like state for (postID, userID) and returns
	// the resulting state. Idempotent under post_likes' composite PK.
	ToggleLike(ctx context.Context, postID, userID uuid.UUID) (liked bool, likesCount int, err error)

	AddComment(ctx context.Context, postID, authorUID uuid.UUID, body string) (*models.Comment, error)

	ListComments(ctx context.Context, postID uuid.UUID, limit int) ([]models.Comment, error)
}
