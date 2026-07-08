package repo

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/utils"
	"github.com/leeryan2000/flashat/wire"
)

type PgxConversationRepo struct {
	Pool *pgxpool.Pool

	// In-process cache of conversation participants, keyed by conversation ID.
	// Avoids re-querying the DB on every chat message. Invalidated on any
	// membership change (Add/Remove/LeaveGroup). Returns a shared slice to
	// callers — safe today since the only consumer only reads it.
	mu    sync.RWMutex
	cache map[uuid.UUID][]*models.ConversationParticipant
}

func (r *PgxConversationRepo) CreateGroupConversation(ctx context.Context, conv *models.Conversation, msg *models.Message, creatorUID uuid.UUID, participantsUID []uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = tx.QueryRow(ctx, `
		INSERT INTO conversations (id, type, group_name, group_avatar_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, type, group_name, group_avatar_url, created_at`,
		conv.ID, conv.Type, conv.GroupName, conv.GroupAvatarUrl,
	).Scan(&conv.ID, &conv.Type, &conv.GroupName, &conv.GroupAvatarUrl, &conv.CreatedAt)

	if err != nil {
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

	err = utils.SaveMessageWithTx(ctx, tx, msg)
	if err != nil {
		slog.Error("save message failed", "error", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PgxConversationRepo) CreateDirectConversation(ctx context.Context, conv *models.Conversation, uid, targetUID uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = utils.CreateDirectConversationWithTx(ctx, conv, tx, uid, targetUID)

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *PgxConversationRepo) ListConversationByUID(ctx context.Context, uid uuid.UUID) ([]*models.Conversation, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT c.id, c.type, c.group_name, c.created_at
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
		if err := rows.Scan(&conv.ID, &conv.Type, &conv.GroupName, &conv.CreatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, conv)
	}

	return convs, nil
}

func (r *PgxConversationRepo) GetConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{}
	err := r.Pool.QueryRow(ctx, `
		SELECT id, type, group_name, created_at
		FROM conversations
		WHERE id = $1`,
		conversationID,
	).Scan(&conv.ID, &conv.Type, &conv.GroupName, &conv.CreatedAt)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *PgxConversationRepo) ListParticipantByID(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationParticipant, error) {
	r.mu.RLock()
	if cached, ok := r.cache[conversationID]; ok {
		r.mu.RUnlock()
		return cached, nil
	}
	r.mu.RUnlock()

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

	r.mu.Lock()
	if r.cache == nil {
		r.cache = make(map[uuid.UUID][]*models.ConversationParticipant)
	}
	r.cache[conversationID] = participants
	r.mu.Unlock()

	return participants, nil
}

func (r *PgxConversationRepo) invalidateParticipantCache(conversationID uuid.UUID) {
	r.mu.Lock()
	delete(r.cache, conversationID)
	r.mu.Unlock()
}

func (r *PgxConversationRepo) AddParticipant(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO conversation_participants (conversation_id, uid)
		VALUES ($1, $2)
		ON CONFLICT (conversation_id, uid) DO NOTHING`,
		conversationID, uid,
	)
	if err != nil {
		return err
	}
	r.invalidateParticipantCache(conversationID)
	return nil
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
	if err != nil {
		return err
	}
	r.invalidateParticipantCache(conversationID)
	return nil
}

func (r *PgxConversationRepo) LeaveGroup(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Check the leaving user's role
	var role string
	err = tx.QueryRow(ctx, `
		SELECT role FROM conversation_participants
		WHERE conversation_id = $1 AND uid = $2`,
		conversationID, uid,
	).Scan(&role)
	if err != nil {
		return err
	}

	if role == "creator" {
		// Find another member to promote
		var nextUID uuid.UUID
		err = tx.QueryRow(ctx, `
			SELECT uid FROM conversation_participants
			WHERE conversation_id = $1 AND uid != $2
			LIMIT 1`,
			conversationID, uid,
		).Scan(&nextUID)

		if err == pgx.ErrNoRows {
			// Last member — delete the whole conversation
			_, err = tx.Exec(ctx, `DELETE FROM conversations WHERE id = $1`, conversationID)
			if err != nil {
				return err
			}
			if err := tx.Commit(ctx); err != nil {
				return err
			}
			r.invalidateParticipantCache(conversationID)
			return nil
		}
		if err != nil {
			return err
		}

		// Promote the next member to creator
		_, err = tx.Exec(ctx, `
			UPDATE conversation_participants
			SET role = 'creator'
			WHERE conversation_id = $1 AND uid = $2`,
			conversationID, nextUID,
		)
		if err != nil {
			return err
		}
	}

	// Remove the leaving user
	_, err = tx.Exec(ctx, `
		DELETE FROM conversation_participants
		WHERE conversation_id = $1 AND uid = $2`,
		conversationID, uid,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	r.invalidateParticipantCache(conversationID)
	return nil
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
			c.id        AS conv_id,
			c.type      AS type,
			CASE 
				WHEN c.type = 'group' THEN c.group_name 
				ELSE u.name
			END AS title,

			CASE 
				WHEN c.type = 'group' THEN c.group_avatar_url
				ELSE u.user_avatar_url
			END AS avatar_url,

			lm.id        AS last_msg_id,
			lm.last_msg_text,
			lm.from_uid  AS last_msg_from,
  			COALESCE((EXTRACT(EPOCH FROM lm.last_msg_ts) * 1000)::bigint, NULL) AS last_msg_ts,

			COALESCE(cc.last_seq, lm.seq, 0) AS last_seq,
			cp.last_read_seq,
			GREATEST(COALESCE(cc.last_seq, lm.seq, 0) - cp.last_read_seq, 0) AS unread_count,

			COALESCE((
				SELECT json_agg(
					json_build_object('uid', p.uid, 'name', u2.name, 'role', p.role)
					ORDER BY CASE WHEN p.role IN ('creator', 'admin') THEN 0 ELSE 1 END, u2.name
				)
				FROM conversation_participants p
				JOIN users u2 ON u2.uid = p.uid
				WHERE p.conversation_id = c.id
			), '[]') AS participants
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
		LEFT JOIN LATERAL (
			SELECT p.uid AS uid
			FROM conversation_participants p
			WHERE p.conversation_id = c.id
				AND p.uid <> $1
			LIMIT 1
		) peer ON c.type = 'direct'
		LEFT JOIN users u
		  ON u.uid = peer.uid
		ORDER BY COALESCE(lm.last_msg_ts, c.created_at) DESC;
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*wire.ConversationSummary

	for rows.Next() {
		c := &wire.ConversationSummary{}
		var participantsJSON json.RawMessage

		if err := rows.Scan(
			&c.ConversationID,
			&c.ConvType,
			&c.Title,
			&c.AvatarURL,
			&c.LastMsgID,
			&c.LastMsgText,
			&c.LastMsgFrom,
			&c.LastMsgTs,
			&c.LastSeq,
			&c.LastReadSeq,
			&c.UnreadCount,
			&participantsJSON, // JSON array of participants
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(participantsJSON, &c.Participants); err != nil {
			return nil, err
		}

		out = append(out, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
