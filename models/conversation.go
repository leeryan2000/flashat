package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID         uuid.UUID `db:"id"`
	Type       string    `db:"type"`        // "direct" or "group"
	DirectKey  *string   `db:"direct_key"`  // nil if group
	GroupName *string   `db:"group_name"`  // nil if direct
	CreatedAt  time.Time `db:"created_at"`
}