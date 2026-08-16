package models

import "github.com/google/uuid"

type Comment struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PostID    uuid.UUID `db:"post_id" json:"post_id"`
	AuthorUID uuid.UUID `db:"author_uid" json:"author_uid"`
	Body      string    `db:"body" json:"body"`
	CreatedAt int64     `db:"created_at" json:"created_at"`
}
