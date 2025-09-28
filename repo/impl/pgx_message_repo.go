package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
)

type PgxMessageRepo struct {
	Pool *pgxpool.Pool
}

const (
	defaultLimit = 50
	maxLimit     = 100
	minLimit     = 1
)

func clampLimit(n int) int {
	if n <= 0 {
		n = defaultLimit
	}
	if n < minLimit {
		n = minLimit
	}
	if n > maxLimit {
		n = maxLimit
	}
	return n
}

func reverse(ms []models.Message) {
	for i, j := 0, len(ms)-1; i < j; i, j = i+1, j-1 {
		ms[i], ms[j] = ms[j], ms[i]
	}
}

func (r *PgxMessageRepo) SaveMessage(ctx context.Context, msg *models.Message) error {
	// Implementation for saving a message using pgx
	err := r.Pool.QueryRow(ctx, `
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
		RETURNING id, conversation_id, seq, from_uid, body, created_at`,
		msg.ID, msg.ConversationID, msg.FromUID, msg.Body,
	).Scan(&msg.ID, &msg.ConversationID, &msg.Seq, &msg.FromUID, &msg.Body, &msg.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgxMessageRepo) ListLatest(ctx context.Context, conversationID uuid.UUID, limit int, uid uuid.UUID) ([]models.Message, error) {
	limit = clampLimit(limit)
	rows, err := r.Pool.Query(ctx, `
		SELECT 
			m.id, 
			m.conversation_id, 
			m.seq, 
			m.from_uid, 
			m.body, 
			(EXTRACT(EPOCH FROM m.created_at) * 1000)::bigint AS created_at
		FROM messages m
		JOIN conversation_participants p
			ON p.conversation_id = m.conversation_id
		AND p.uid = $1
		WHERE m.conversation_id = $2
		ORDER BY m.seq DESC
		LIMIT $3`,
		uid, conversationID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgs := make([]models.Message, 0, limit)
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Seq, &msg.FromUID, &msg.Body, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	reverse(msgs) // Reverse to get the latest messages first, since the message order is descending
	return msgs, nil
}

func (r *PgxMessageRepo) ListBefore(ctx context.Context, conversationID uuid.UUID, beforeSeq int64, limit int) ([]models.Message, error) {
	limit = clampLimit(limit)
	rows, err := r.Pool.Query(ctx, `
		SELECT id, conversation_id, seq, from_uid, body, created_at
		FROM messages
		WHERE conversation_id = $1 AND seq < $2
		ORDER BY seq DESC
		LIMIT $3`,
		conversationID, beforeSeq, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Seq, &msg.FromUID, &msg.Body, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	reverse(messages) // Reverse to get the latest messages first
	return messages, nil
}

func (r *PgxMessageRepo) ListAfter(ctx context.Context, conversationID uuid.UUID, afterSeq int64, limit int) ([]models.Message, error) {
	limit = clampLimit(limit)
	rows, err := r.Pool.Query(ctx, `
		SELECT id, conversation_id, seq, from_uid, body, created_at
		FROM messages
		WHERE conversation_id = $1 AND seq > $2
		ORDER BY seq ASC
		LIMIT $3`,
		conversationID, afterSeq, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Seq, &msg.FromUID, &msg.Body, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
