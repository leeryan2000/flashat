package wire

import "github.com/google/uuid"

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

type Ack struct {
	ServerMsgID uuid.UUID `json:"server_msg_id"`
	ServerTS    int64     `json:"server_ts"`
	Status      string    `json:"status"` // success, failed
}
