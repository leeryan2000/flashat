package wire

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
	Status string `json:"status"` // sent, failed
}
