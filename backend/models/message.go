package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID       `db:"id" json:"id"`
	ConversationID uuid.UUID       `db:"conversation_id" json:"conversation_id"`
	Seq            int64           `db:"seq" json:"seq"`
	FromUID        uuid.UUID       `db:"from_uid" json:"from_uid"`
	FromName       string          `db:"from_name" json:"from_name"`
	Body           json.RawMessage `db:"body" json:"body"`
	CreatedAt      int64           `db:"created_at" json:"created_at"`
}
