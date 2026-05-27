package service

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/wire"
)

type BroadcastObserver struct {
	Hub Hub
}

func (o *BroadcastObserver) OnMessageSaved(e *MessageEvent) {
	fromUID, err := uuid.Parse(e.OriginalEnv.FromUID)
	if err != nil {
		log.Println("❌ BroadcastObserver: failed to parse fromUID:", err)
		return
	}

	outEnv := &wire.Msg{
		Type:           wire.Chat,
		ConversationID: e.OriginalEnv.ConversationID,
		ServerMsgID:    e.SavedMsg.ID.String(),
		FromUID:        e.OriginalEnv.FromUID,
		FromName:       e.SavedMsg.FromName,
		Seq:            e.SavedMsg.Seq,
		Ts:             e.SavedMsg.CreatedAt,
		Body:           e.OriginalEnv.Body,
	}

	payload, err := json.Marshal(outEnv)
	if err != nil {
		log.Println("❌ BroadcastObserver: failed to marshal broadcast:", err)
		return
	}

	o.Hub.BroadcastToParticipant(e.Participants, fromUID, payload)
	o.Hub.SendToUID(fromUID, payload) // sync sender's other tabs/devices
}
