package models

import (
	"time"

	"github.com/google/uuid"
)

type ConversationParticipant struct {
	ConversationID uuid.UUID `db:"conversation_id" json:"conversation_id"`
	UID            uuid.UUID `db:"uid" json:"uid"`                     // UID string
	Role           string    `db:"role" json:"role"`                   // "admin", "member", etc.
	LastReadSeq    int64     `db:"last_read_seq" json:"last_read_seq"` // last read message sequence number
	JoinedAt       time.Time `db:"joined_at" json:"joined_at"`         // timestamp when the user joined the conversation
}
