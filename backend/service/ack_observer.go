package service

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/wire"
)

type AckObserver struct {
	Hub Hub
}

func (o *AckObserver) OnMessageSaved(e *MessageEvent) {
	fromUID, err := uuid.Parse(e.OriginalEnv.FromUID)
	if err != nil {
		log.Println("❌ AckObserver: failed to parse fromUID:", err)
		return
	}

	ackEnv := &wire.Msg{
		Type:           wire.Ack,
		ConversationID: e.OriginalEnv.ConversationID,
		ClientMsgID:    e.OriginalEnv.ClientMsgID,
		ServerMsgID:    e.SavedMsg.ID.String(),
		FromUID:        e.OriginalEnv.FromUID,
		FromName:       e.SavedMsg.FromName,
		Seq:            e.SavedMsg.Seq,
		Ts:             e.SavedMsg.CreatedAt,
		Body:           e.OriginalEnv.Body,
	}

	payload, err := json.Marshal(ackEnv)
	if err != nil {
		log.Println("❌ AckObserver: failed to marshal ack:", err)
		return
	}

	o.Hub.SendToUID(fromUID, payload)
}
