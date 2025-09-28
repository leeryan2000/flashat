package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID       `db:"id" json:"id"`
	ConversationID uuid.UUID       `db:"conversation_id" json:"conversation_id"`
	Seq            int64           `db:"seq" json:"seq"`
	FromUID        uuid.UUID       `db:"from_uid" json:"from_uid"`     // UID string
	Body           json.RawMessage `db:"body" json:"body"`             // flexible schema for now
	CreatedAt      time.Time       `db:"created_at" json:"created_at"` // ***** type incorrect should reveive number
}
