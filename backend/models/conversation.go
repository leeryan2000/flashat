package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Type           string    `db:"type" json:"type"`             // "direct" or "group"
	GroupName      *string   `db:"group_name" json:"group_name"` // nil if direct
	GroupAvatarUrl *string   `db:"group_avatar_url" json:"group_avatar_url"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
