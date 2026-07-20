package repo

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/utils"
	"github.com/leeryan2000/flashat/wire"
)

type PgxFriendshipRepo struct {
	Pool *pgxpool.Pool
}

func (r *PgxFriendshipRepo) RequestFriendship(ctx context.Context, requesterUID uuid.UUID, email string) (*wire.Friendship, error) {
	var f wire.Friendship
	err := r.Pool.QueryRow(ctx, `
		WITH target AS (
			SELECT u.uid, u.name, u.email, u.user_avatar_url
			FROM users u
			WHERE u.email = $2
			AND u.uid <> $1
			AND NOT EXISTS (
				SELECT 1 FROM friendships f
				WHERE (f.requester_uid = $1 AND f.receiver_uid = u.uid)
				   OR (f.requester_uid = u.uid AND f.receiver_uid = $1)
			)
		),
		ins AS (
			INSERT INTO friendships (requester_uid, receiver_uid, status)
			SELECT $1, target.uid, 'pending' FROM target
		)
		SELECT uid, name, email, user_avatar_url FROM target
		`,
		requesterUID, email,
	).Scan(&f.UID, &f.Name, &f.Email, &f.AvatarURL)

	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	f.Status = "pending"
	f.Direction = "outgoing"
	return &f, nil
}

// Transaction: for making sure the conversation is created only if friendship is accepted
func (r *PgxFriendshipRepo) AcceptFriendship(ctx context.Context, conv *models.Conversation, msg *models.Message, requesterUID, receiverUID uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	tag, err := r.Pool.Exec(ctx, `
		UPDATE friendships
		SET status = 'accepted'
		WHERE requester_uid = $1 AND receiver_uid = $2 AND status = 'pending'`,
		requesterUID, receiverUID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		slog.Warn("no pending friendship found to accept")
		return pgx.ErrNoRows
	}

	err = utils.CreateDirectConversationWithTx(ctx, conv, tx, requesterUID, receiverUID)
	if err != nil {
		return err
	}

	err = utils.SaveMessageWithTx(ctx, tx, msg)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PgxFriendshipRepo) DeleteFriendship(ctx context.Context, uid1, uid2 uuid.UUID) error {
	// Implementation for deleting friendship using pgx
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		DELETE FROM friendships
		WHERE (requester_uid = $1 AND receiver_uid = $2)
		   OR (requester_uid = $2 AND receiver_uid = $1)`,
		uid1, uid2,
	)

	err = utils.RemoveConversationWithTx(ctx, tx, uid1, uid2)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PgxFriendshipRepo) RemovePendingFriendship(ctx context.Context, requesterUID, receiverUID uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, `
		DELETE FROM friendships
		WHERE requester_uid = $1 AND receiver_uid = $2 AND status = 'pending'`,
		requesterUID, receiverUID,
	)
	return err
}

func (r *PgxFriendshipRepo) BlockUser(ctx context.Context, blockerUID, targetUID uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Remove any existing friendship row
	_, err = tx.Exec(ctx, `
		DELETE FROM friendships
		WHERE (requester_uid = $1 AND receiver_uid = $2)
		   OR (requester_uid = $2 AND receiver_uid = $1)`,
		blockerUID, targetUID,
	)
	if err != nil {
		return err
	}

	// Insert blocked row with blocker as requester so direction is clear
	_, err = tx.Exec(ctx, `
		INSERT INTO friendships (requester_uid, receiver_uid, status)
		VALUES ($1, $2, 'blocked')`,
		blockerUID, targetUID,
	)
	if err != nil {
		return err
	}

	// Remove the shared direct conversation if one exists
	if err = utils.RemoveConversationWithTx(ctx, tx, blockerUID, targetUID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PgxFriendshipRepo) ListFriendships(ctx context.Context, uid uuid.UUID) ([]wire.Friendship, error) {
	// Implementation for listing friendships using pgx
	rows, err := r.Pool.Query(ctx, `
		SELECT 
			u.uid,
			u.name AS name,
			u.user_avatar_url,
			u.email,
			-- Finds the shared conversation ID
			(
				SELECT p1.conversation_id
				FROM conversation_participants p1
				JOIN conversation_participants p2 ON p1.conversation_id = p2.conversation_id
				JOIN conversations c ON p1.conversation_id = c.id
				WHERE p1.uid = $1        -- You are in it
				AND p2.uid = u.uid     -- Friend is in it
				AND c.type = 'direct'      -- It is a 1-on-1 chat
				LIMIT 1
			) AS conversation_id,
			f.status,
			CASE
				WHEN f.requester_uid = $1 THEN 'outgoing'
				ELSE 'incoming'
			END AS direction
		FROM users u
		JOIN friendships f ON 
			(f.requester_uid = $1 AND f.receiver_uid = u.uid)
			OR 
			(f.receiver_uid = $1 AND f.requester_uid = u.uid);`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []wire.Friendship
	for rows.Next() {
		var friend wire.Friendship
		if err := rows.Scan(&friend.UID, &friend.Name, &friend.AvatarURL, &friend.Email, &friend.DirectConversationID, &friend.Status, &friend.Direction); err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (r *PgxFriendshipRepo) GetFriendshipStatus(ctx context.Context, userAUID, userBUID uuid.UUID) (string, error) {
	// Implementation for getting friendship status using pgx
	return "", nil
}
