package models

import "github.com/google/uuid"

// LikesCount/CommentsCount/LikedByMe mirror how backend/models/Message
// carries FromName — a value composed by the query (join/counter column),
// not a literal 1:1 column mapping, but still returned directly as JSON
// rather than through a separate wire/response type.
type Post struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Seq           int64     `db:"seq" json:"seq"`
	AuthorUID     uuid.UUID `db:"author_uid" json:"author_uid"`
	Body          string    `db:"body" json:"body"`
	LikesCount    int       `db:"likes_count" json:"likes_count"`
	CommentsCount int       `db:"comments_count" json:"comments_count"`
	LikedByMe     bool      `db:"liked_by_me" json:"liked_by_me"`
	CreatedAt     int64     `db:"created_at" json:"created_at"`
}
