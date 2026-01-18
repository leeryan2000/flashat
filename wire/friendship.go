package wire

import "github.com/google/uuid"

type Friendship struct {
	UID                  uuid.UUID  `db:"uid" json:"uid"`
	Name                 string     `db:"name" json:"name"`
	Email                string     `db:"email" json:"email"`
	AvatarURL            string     `db:"avatar_url" json:"avatar_url"`
	DirectConversationID *uuid.UUID `db:"direct_conversation_id" json:"direct_conversation_id"`
	Status               string     `db:"status" json:"status"`       // pending | accepted | blocked'
	Direction            string     `db:"direction" json:"direction"` // incoming | outgoing
}
