package server

type MsgType string

const (
	MsgJoin   MsgType = "join"
	MsgLeave  MsgType = "leave"
	MsgChat   MsgType = "chat"   // room broadcast
	MsgDirect MsgType = "direct" // 1:1
	MsgSystem MsgType = "system"
	MsgAck    MsgType = "ack"
	MsgPing   MsgType = "ping"
	MsgPong   MsgType = "pong"
)

type Envelope struct {
	Type    MsgType `json:"type"`
	RoomID  string  `json:"roomId,omitempty"`
	FromUID string  `json:"fromUid,omitempty"`
	ToUID   string  `json:"toUid,omitempty"`
	Body    string  `json:"body,omitempty"`
	Seq     int64   `json:"seq,omitempty"` // optional for acks
	Ts      int64   `json:"ts,omitempty"`
}
