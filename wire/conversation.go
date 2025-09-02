package wire

import (
	"time"

	"github.com/google/uuid"
)

type ConversationSummary struct {
	ConvType       string    `json:"type"` // group | direct
	ConversationID uuid.UUID `json:"conversation_id"`
	// fields that are possible being null
	// ***** find a way to display direct conversations title
	Title       *string    `json:"title,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	LastMsgID   *uuid.UUID `json:"last_message_id,omitempty"`
	LastMsgText *string    `json:"last_message_text,omitempty"`
	LastMsgFrom *string    `json:"last_message_from,omitempty"`
	LastMsgTs   *time.Time `json:"last_message_ts,omitempty"`

	LastSeq     int64 `json:"last_seq"`
	LastReadSeq int64 `json:"last_read_seq"`
	UnreadCount int64 `json:"unread_count"`
}
