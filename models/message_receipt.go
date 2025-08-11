package models

import (
	"time"

	"github.com/google/uuid"
)


type MessageReceipt struct {
	MessageID     uuid.UUID    `db:"message_id"` // ID of the message
	UID		   uuid.UUID    `db:"uid"`        // User ID of the recipient
	DeliveredAt     *time.Time `db:"delivered_at"` // Timestamp when the message was delivered
	ReadAt		 *time.Time `db:"read_at"` // Timestamp when the message was read, nil if not read
}