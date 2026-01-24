package utils

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/leeryan2000/flashat/models"
)

func SaveMessageWithTx(ctx context.Context, tx pgx.Tx, msg *models.Message) error {
	err := tx.QueryRow(ctx, `
		WITH s AS (
			Update conversation_counters
			SET last_seq = last_seq + 1
			WHERE conversation_id = $2
				AND EXISTS (
					SELECT 1
					FROM conversation_participants p
					WHERE p.conversation_id = $2
						AND p.uid = $3
					)
			RETURNING last_seq
		)
		INSERT INTO messages (id, conversation_id, seq, from_uid, body)
		SELECT $1, $2, s.last_seq, $3, $4
		FROM s
		RETURNING id, conversation_id, seq, from_uid, body, (EXTRACT(EPOCH FROM created_at) * 1000)::bigint AS created_at`,
		msg.ID, msg.ConversationID, msg.FromUID, msg.Body,
	).Scan(&msg.ID, &msg.ConversationID, &msg.Seq, &msg.FromUID, &msg.Body, &msg.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func CreateDirectConversationWithTx(ctx context.Context, conv *models.Conversation, tx pgx.Tx, uid1, uid2 uuid.UUID) error {
	// Implementation for creating direct conversation using pgx
	batch := &pgx.Batch{}

	batch.Queue(`
		INSERT INTO conversations (id, type)
		VALUES ($1, $2)
		RETURNING created_at`,
		conv.ID, conv.Type,
	)

	batch.Queue(`
		INSERT INTO conversation_counters (conversation_id)
		VALUES ($1)`,
		conv.ID,
	)

	batch.Queue(`
		INSERT INTO conversation_participants (conversation_id, uid)
		VALUES ($1, $2), ($1, $3)`,
		conv.ID, uid1, uid2,
	)

	br := tx.SendBatch(ctx, batch)
	defer br.Close()
	return br.Close()
}

func RemoveConversationWithTx(ctx context.Context, tx pgx.Tx, uid1, uid2 uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		DELETE FROM conversations 
		WHERE id IN(
			SELECT p1.conversation_id
			FROM conversation_participants p1
			JOIN conversation_participants p2 ON p1.conversation_id = p2.conversation_id
			JOIN conversations c ON p1.conversation_id = c.id
			WHERE p1.uid = $1        -- User 1 is in it
			AND p2.uid = $2         -- User 2 is in it
			AND c.type = 'direct'   -- It is a 1-on-1 chat
			LIMIT 1
		);`,
		uid1, uid2,
	)
	return err
}
