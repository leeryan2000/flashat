package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
)

type PgxMessageRepo struct {
	Pool *pgxpool.Pool
}

func (r *PgxMessageRepo) SaveMessage(ctx context.Context, msg *models.Message) error {
	// Implementation for saving a message using pgx
	err := r.Pool.QueryRow(ctx, `
		WITH s AS (
			Update conversation_counters
			SET last_seq = last_seq + 1
			WHERE conversation_id = $1
			RETURNING last_seq
		)
		INSERT INTO messages (id, conversation_id, seq, from_uid, body)
		SELECT $1, $2, s.last_seq, $3, $4
		FROM s
		RETURNING id, conversation_id, seq, from_uid, body, created_at
		`, msg.ID, msg.ConversationID, msg.FromUID, msg.Body,
	).Scan(&msg.ID, &msg.ConversationID, &msg.Seq, &msg.FromUID, &msg.Body, &msg.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
