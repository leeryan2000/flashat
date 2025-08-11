package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID       `db:"id"`
	ConversationID uuid.UUID       `db:"conversation_id"`
	Seq            int64           `db:"seq"`
	FromUID        uuid.UUID       `db:"from_uid"` // UID string
	Body           json.RawMessage `db:"body"`     // flexible schema for now
	CreatedAt      time.Time       `db:"created_at"`
}
