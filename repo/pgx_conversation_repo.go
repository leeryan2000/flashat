package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConversationRepo struct {
	Pool *pgxpool.Pool
}

func (r *PgxConversationRepo) CreateGroupConversation(ctx context.Context, creatorUID uuid.UUID, participantsUID []uuid.UUID, groupName string) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	convID := uuid.New()
	// conversations: id, type, direct_key(NULL), created_at (default)
	_, err = tx.Exec(ctx, `
		INSERT INTO conversations (id, type, group_name)
		VALUES ($1, 'group', $2)`,
		convID, groupName,
	)
	if err != nil {
		return err
	}

	// counters: start at 0 in default
	_, err = tx.Exec(ctx, `
		INSERT INTO conversation_counters (conversation_id)
		VALUES ($1)`,
		convID,
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
		return err
	}

	return tx.Commit(ctx)
}
