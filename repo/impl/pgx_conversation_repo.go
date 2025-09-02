package repo

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/wire"
)

type PgxConversationRepo struct {
	Pool *pgxpool.Pool
}

func (r *PgxConversationRepo) CreateGroupConversation(ctx context.Context, conv *models.Conversation, creatorUID uuid.UUID, participantsUID []uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// conversations: id, type, direct_key(NULL), created_at (default)
	err = tx.QueryRow(ctx, `
		INSERT INTO conversations (id, type, group_name)
		VALUES ($1, $2, $3)
		RETURNING id, type, group_name, created_at`,
		conv.ID, conv.Type, conv.GroupName,
	).Scan(&conv.ID, &conv.Type, &conv.GroupName, &conv.CreatedAt)

	if err != nil {
		log.Println("Failed to create group conversation:", err)
		return err
	}

	// counters: start at 0 in default
	_, err = tx.Exec(ctx, `
		INSERT INTO conversation_counters (conversation_id)
		VALUES ($1)`,
		conv.ID,
	)
	if err != nil {
		return err
	}

	// Add creator in participants table with role 'creator'
	batch := &pgx.Batch{}
	batch.Queue(`
		INSERT INTO conversation_participants (conversation_id, uid, role)
		VALUES ($1, $2, 'creator')
		ON CONFLICT (conversation_id, uid) DO NOTHING`,
		conv.ID, creatorUID,
	)

	// Insert others as 'member', avoiding duplicates
	set := map[uuid.UUID]struct{}{creatorUID: {}}
	for _, puid := range participantsUID {
		if _, exists := set[puid]; !exists {
			batch.Queue(`
				INSERT INTO conversation_participants (conversation_id, uid)
				VALUES ($1, $2)
				ON CONFLICT (conversation_id, uid) DO NOTHING`,
				conv.ID, puid,
			)
			// put the uuid to the set to avoid duplicates
			set[puid] = struct{}{}
		}
	}
	br := tx.SendBatch(ctx, batch)
	if err = br.Close(); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PgxConversationRepo) GetOrCreateDirectConversation(ctx context.Context, conv *models.Conversation, uid, targetUID uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Find if the direct conversation already exists
	err = tx.QueryRow(ctx, `
		SELECT id, type, direct_key, created_at
		FROM conversations
		WHERE direct_key = $1`,
		conv.DirectKey,
	).Scan(&conv.ID, &conv.Type, &conv.DirectKey, &conv.CreatedAt)

	if err == pgx.ErrNoRows {
		// Not found, create a new direct conversation
		err = tx.QueryRow(ctx, `
			INSERT INTO conversations (id, type, direct_key)
			VALUES ($1, $2, $3)
			RETURNING id, type, direct_key, created_at`,
			conv.ID, conv.Type, conv.DirectKey,
		).Scan(&conv.ID, &conv.Type, &conv.DirectKey, &conv.CreatedAt)
		if err != nil {
			return err
		}

		// add conversation counter
		_, err = tx.Exec(ctx, `
			INSERT INTO conversation_counters (conversation_id)
			VALUES ($1)`,
			conv.ID,
		)
		if err != nil {
			return err
		}

		// Add participants to the conversation
		batch := &pgx.Batch{}
		batch.Queue(`
			INSERT INTO conversation_participants (conversation_id, uid)
			VALUES ($1, $2), ($1, $3)
			ON CONFLICT (conversation_id, uid) DO NOTHING`,
			conv.ID, uid, targetUID,
		)
		br := tx.SendBatch(ctx, batch)
		if err = br.Close(); err != nil {
			return err
		}
	} else if err != nil {
		return err // Some other error occurred
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *PgxConversationRepo) ListConversationByUID(ctx context.Context, uid uuid.UUID) ([]*models.Conversation, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT c.id, c.type, c.direct_key, c.group_name, c.created_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.uid = $1`,
		uid,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []*models.Conversation
	for rows.Next() {
		conv := &models.Conversation{}
		if err := rows.Scan(&conv.ID, &conv.Type, &conv.DirectKey, &conv.GroupName, &conv.CreatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, conv)
	}

	return convs, nil
}

func (r *PgxConversationRepo) GetConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{}
	err := r.Pool.QueryRow(ctx, `
		SELECT id, type, direct_key, group_name, created_at
		FROM conversations
		WHERE id = $1`,
		conversationID,
	).Scan(&conv.ID, &conv.Type, &conv.DirectKey, &conv.GroupName, &conv.CreatedAt)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *PgxConversationRepo) ListParticipantByID(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationParticipant, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT cp.conversation_id, cp.uid, cp.role, cp.last_read_seq
		FROM conversation_participants cp
		WHERE cp.conversation_id = $1`,
		conversationID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*models.ConversationParticipant
	for rows.Next() {
		p := &models.ConversationParticipant{}
		if err := rows.Scan(&p.ConversationID, &p.UID, &p.Role, &p.LastReadSeq); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

func (r *PgxConversationRepo) AddParticipant(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO conversation_participants (conversation_id, uid)
		VALUES ($1, $2)
		ON CONFLICT (conversation_id, uid) DO NOTHING`,
		conversationID, uid,
	)
	return err
}

func (r *PgxConversationRepo) ModifyParticipant(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID, role string) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE conversation_participants
		SET role = $1
		WHERE conversation_id = $2 AND uid = $3`,
		role, conversationID, uid,
	)
	return err
}

func (r *PgxConversationRepo) RemoveParticipant(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, `
		DELETE FROM conversation_participants
		WHERE conversation_id = $1 AND uid = $2`,
		conversationID, uid,
	)
	return err
}

func (r *PgxConversationRepo) UpdateLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID, seq int64) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE conversation_participants
		SET last_read_seq = $3
		WHERE conversation_id = $1 AND uid = $2`,
		conversationID, uid, seq,
	)
	return err
}

func (r *PgxConversationRepo) GetLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) (int64, error) {
	var seq int64
	err := r.Pool.QueryRow(ctx, `
		SELECT last_read_seq
		FROM conversation_participants
		WHERE conversation_id = $1 AND uid = $2`,
		conversationID, uid,
	).Scan(&seq)

	if err != nil {
		return 0, err
	}

	return seq, nil
}

func (r *PgxConversationRepo) GetSummary(ctx context.Context, uid uuid.UUID) ([]*wire.ConversationSummary, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT
			c.id        AS conversation_id,
			c.type      AS type,
			CASE WHEN c.type = 'group' THEN c.group_name END AS title,

			lm.id        AS last_message_id,
			lm.last_msg_text,
			lm.from_uid  AS last_message_from,
			lm.last_msg_ts,  -- timestamptz

			COALESCE(cc.last_seq, lm.seq, 0) AS last_seq,
			cp.last_read_seq,
			GREATEST(COALESCE(cc.last_seq, lm.seq, 0) - cp.last_read_seq, 0) AS unread_count
		FROM conversations c
		JOIN conversation_participants cp
		  ON cp.conversation_id = c.id
		 AND cp.uid = $1
		LEFT JOIN conversation_counters cc
		  ON cc.conversation_id = c.id
		LEFT JOIN LATERAL (
			SELECT
				m.id,
				m.seq,
				m.from_uid,
				(m.body->>'text') AS last_msg_text,
				m.created_at      AS last_msg_ts  -- timestamptz
			FROM messages m
			WHERE m.conversation_id = c.id
			ORDER BY m.seq DESC
			LIMIT 1
		) lm ON TRUE
		ORDER BY COALESCE(lm.last_msg_ts, c.created_at) DESC;
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*wire.ConversationSummary

	for rows.Next() {
		c := &wire.ConversationSummary{}

		if err := rows.Scan(
			&c.ConversationID,
			&c.ConvType,
			&c.Title,
			&c.LastMsgID,
			&c.LastMsgText,
			&c.LastMsgFrom,
			&c.LastMsgTs,
			&c.LastSeq,
			&c.LastReadSeq,
			&c.UnreadCount,
		); err != nil {
			return nil, err
		}

		out = append(out, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
