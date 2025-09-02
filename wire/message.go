package wire

import "encoding/json"

const (
	MsgJoin   MsgType = "join"
	MsgLeave  MsgType = "leave"
	MsgChat   MsgType = "chat"
	MsgSystem MsgType = "system"
	MsgAck    MsgType = "ack"
	MsgPing   MsgType = "ping"
	MsgPong   MsgType = "pong"
)

type MsgType string

type MsgEnvelope struct {
	Type           MsgType `json:"type"`
	ConversationID string  `json:"conversation_id,omitempty"`
	ClientMsgID    string  `json:"client_msg_id,omitempty"`
	ServerMsgID    string  `json:"server_msg_id,omitempty"`
	FromUID        string  `json:"from_uid,omitempty"`
	// discover whether seq is needed from the client
	Seq  int64           `json:"seq,omitempty"`
	Ts   int64           `json:"ts,omitempty"`
	Body json.RawMessage `json:"body,omitempty"`
}

type Ack struct {
	Status string `json:"status"` // sent, failed
}
