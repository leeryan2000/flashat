package wire

import "github.com/google/uuid"

const (
	MsgJoin   MsgType = "join"
	MsgLeave  MsgType = "leave"
	MsgChat   MsgType = "chat"
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
