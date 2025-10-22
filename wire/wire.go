package wire

import "encoding/json"

const (
	Join      WireType = "join"
	Leave     WireType = "leave"
	Chat      WireType = "chat"
	System    WireType = "system"
	Ack       WireType = "ack"
	Ping      WireType = "ping"
	Pong      WireType = "pong"
	Subscribe WireType = "subscribe"
)

type WireType string

type MsgEnvelope struct {
	Type           WireType `json:"type"`
	ConversationID string   `json:"conversation_id,omitempty"`
	ClientMsgID    string   `json:"client_msg_id,omitempty"`
	ServerMsgID    string   `json:"server_msg_id,omitempty"`
	FromUID        string   `json:"from_uid,omitempty"`
	// discover whether seq is needed from the client
	Seq  int64           `json:"seq,omitempty"`
	Ts   int64           `json:"ts,omitempty"`
	Body json.RawMessage `json:"body,omitempty"`
}

type MsgAck struct {
	Status string `json:"status"` // sent, failed
}
