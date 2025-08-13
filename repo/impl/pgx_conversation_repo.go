package repo

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/utils"
)

type PgxConversationRepo struct {
	Pool *pgxpool.Pool
}

func (r *PgxConversationRepo) CreateGroupConversation(ctx context.Context, creatorUID uuid.UUID, participantsUID []uuid.UUID, groupName string) (*models.Conversation, error) {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	conv := &models.Conversation{}
	convID := uuid.New()
	// conversations: id, type, direct_key(NULL), created_at (default)
	err = tx.QueryRow(ctx, `
		INSERT INTO conversations (id, type, group_name)
		VALUES ($1, 'group', $2)
		RETURNING id, type, group_name, created_at`,
		convID, groupName,
	).Scan(&conv.ID, &conv.Type, &conv.GroupName, &conv.CreatedAt)

	if err != nil {
		log.Println("Failed to create group conversation:", err)
		return nil, err
	}

	// counters: start at 0 in default
	_, err = tx.Exec(ctx, `
		INSERT INTO conversation_counters (conversation_id)
		VALUES ($1)`,
		convID,
	)
	if err != nil {
		return nil, err
	}

	// Add creator in participants table with role 'creator'
	batch := &pgx.Batch{}
	batch.Queue(`
		INSERT INTO conversation_participants (conversation_id, uid, role)
		VALUES ($1, $2, 'creator')
		ON CONFLICT (conversation_id, uid) DO NOTHING`,
		convID, creatorUID,
	)

	// Insert others as 'member', avoiding duplicates
	set := map[uuid.UUID]struct{}{creatorUID: {}}
	for _, u := range participantsUID {
		if _, exists := set[u]; !exists {
			batch.Queue(`
				INSERT INTO conversation_participants (conversation_id, uid)
				VALUES ($1, $2)
				ON CONFLICT (conversation_id, uid) DO NOTHING`,
				convID, u,
			)
			// put the uuid to the set to avoid duplicates
			set[u] = struct{}{}
		}
	}
	br := tx.SendBatch(ctx, batch)
	if err = br.Close(); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *PgxConversationRepo) GetOrCreateDirectConversation(ctx context.Context, uidA, uidB uuid.UUID) (*models.Conversation, error) {
	directKey := utils.CanonDirectKey(uidA.String(), uidB.String())
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	conv := &models.Conversation{}
	// Find if the direct conversation already exists
	err = tx.QueryRow(ctx, `
		SELECT id, type, direct_key, created_at
		FROM conversations
		WHERE direct_key = $1`,
		directKey,
	).Scan(&conv.ID, &conv.Type, &conv.DirectKey, &conv.CreatedAt)

	if err == pgx.ErrNoRows {
		// Not found, create a new direct conversation
		convID := uuid.New()

		err = tx.QueryRow(ctx, `
			INSERT INTO conversations (id, type, direct_key)
			VALUES ($1, 'direct', $2)
			RETURNING id, type, direct_key, created_at`,
			convID, directKey,
		).Scan(&conv.ID, &conv.Type, &conv.DirectKey, &conv.CreatedAt)
		if err != nil {
			return nil, err
		}

		// add conversation counter
		_, err = tx.Exec(ctx, `
			INSERT INTO conversation_counters (conversation_id)
			VALUES ($1)`,
			convID,
		)
		if err != nil {
			return nil, err
		}

		// Add participants to the conversation
		batch := &pgx.Batch{}
		batch.Queue(`
			INSERT INTO conversation_participants (conversation_id, uid)
			VALUES ($1, $2), ($1, $3)
			ON CONFLICT (conversation_id, uid) DO NOTHING`,
			conv.ID, uidA, uidB,
		)
		br := tx.SendBatch(ctx, batch)
		if err = br.Close(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err // Some other error occurred
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return conv, nil
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
