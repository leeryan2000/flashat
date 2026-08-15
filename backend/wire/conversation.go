package wire

import (
	"github.com/google/uuid"
)

type ParticipantSummary struct {
	UID       uuid.UUID `json:"uid"`
	Name      string    `json:"name"`
	Role      string    `json:"role"` // "admin" | "member"
	AvatarURL string    `json:"avatar_url"`
}

// conversation summary
type ConversationSummary struct {
	ConversationID uuid.UUID `json:"id"`
	ConvType       string    `json:"type"` // group | direct

	Title       string    `json:"title,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	LastMsgID   uuid.UUID `json:"last_msg_id,omitempty"`
	LastMsgText string    `json:"last_msg_text,omitempty"`
	LastMsgFrom uuid.UUID `json:"last_msg_from,omitempty"`
	LastMsgTs   int64     `json:"last_msg_ts,omitempty"`

	LastSeq     int64 `json:"last_seq"`
	LastReadSeq int64 `json:"last_read_seq"`
	UnreadCount int64 `json:"unread_count"`

	Participants []ParticipantSummary `json:"participants"`
}
