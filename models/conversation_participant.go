package models

import (
	"time"

	"github.com/google/uuid"
)

type ConversationParticipant struct {
	ConversationID uuid.UUID `db:"conversation_id"`
	UID            uuid.UUID `db:"uid"` // UID string
	Role           string    `db:"role"`
	LastReadSeq    int64     `db:"last_read_seq"`
	JoinedAt       time.Time `db:"joined_at"`
}
