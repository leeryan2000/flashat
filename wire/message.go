package wire

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	MsgJoin   MsgType = "join"
	MsgLeave  MsgType = "leave"
	MsgGroup  MsgType = "group"  // room broadcast
	MsgDirect MsgType = "direct" // 1:1
	MsgSystem MsgType = "system"
	MsgAck    MsgType = "ack"
	MsgPing   MsgType = "ping"
	MsgPong   MsgType = "pong"
)

type MsgSend struct {
	ConversationID uuid.UUID       `json:"conversation_id"`
	Body           json.RawMessage `json:"body"`
}
