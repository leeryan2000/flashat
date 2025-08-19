package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Type      string    `db:"type" json:"type"`             // "direct" or "group"
	DirectKey *string   `db:"direct_key" json:"direct_key"` // nil if group
	GroupName *string   `db:"group_name" json:"group_name"` // nil if direct
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
