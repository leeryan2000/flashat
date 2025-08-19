package wire

import (
	"encoding/json"
)

type MsgType string

type Envelope struct {
	Type           MsgType         `json:"type"`
	ConversationID string          `json:"conversation_id,omitempty"`
	FromUID        string          `json:"from_uid,omitempty"`
	Seq            int64           `json:"seq,omitempty"` // optional for acks
	Ts             int64           `json:"ts,omitempty"`
	Body           json.RawMessage `json:"body,omitempty"`
}
