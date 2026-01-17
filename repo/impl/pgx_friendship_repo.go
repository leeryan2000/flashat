package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
)

type PgxFriendshipRepo struct {
	Pool *pgxpool.Pool

	convRepo *PgxConversationRepo
	msgRepo  *PgxMessageRepo
}

func (r *PgxFriendshipRepo) RequestFriendship(ctx context.Context, requesterUID uuid.UUID, email string) error {
	// Implementation for requesting friendship using pgx
	tag, err := r.Pool.Exec(ctx, `
		INSERT INTO friendships (requester_uid, receiver_uid, status)
		SELECT
			$1,
			u.uid,
			'pending'
		FROM users u
		WHERE u.email = $2
		AND u.uid <> $1
		AND NOT EXISTS (
			SELECT 1 FROM friendships f
			WHERE (f.requester_uid = $1 AND f.receiver_uid = u.uid)
			   OR (f.requester_uid = u.uid AND f.receiver_uid = $1)
		)
		`,
		requesterUID, email,
	)

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return err
}

// Transaction: for making sure the conversation is created only if friendship is accepted
func (r *PgxFriendshipRepo) AcceptFriendship(ctx context.Context, conv *models.Conversation, msg *models.Message, requesterUID, receiverUID uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = r.Pool.Exec(ctx, `
		UPDATE friendships
		SET status = 'accepted'
		WHERE requester_uid = $1 AND receiver_uid = $2 AND status = 'pending'`,
		requesterUID, receiverUID,
	)
	if err != nil {
		return err
	}

	err = r.convRepo.CreateDirectConversationWithTx(ctx, conv, tx, requesterUID, receiverUID)
	if err != nil {
		return err
	}

	err = r.msgRepo.SaveMessageWithTx(ctx, tx, msg)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// ***** continue: remove the conversation after deletion of friendship
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

	err = r.convRepo.RemoveConversationWithTx(ctx, tx, uid1, uid2)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PgxFriendshipRepo) RejectFriendship(ctx context.Context, requesterUID, receiverUID uuid.UUID) error {
	// Implementation for rejecting friendship using pgx
	_, err := r.Pool.Exec(ctx, `
		DELETE FROM friendships
		WHERE requester_uid = $1 AND receiver_uid = $2 AND status = 'pending'`,
		requesterUID, receiverUID,
	)
	return err
}

func (r *PgxFriendshipRepo) BlockUser(ctx context.Context, requesterUID, receiverUID uuid.UUID) error {
	// Implementation for blocking user using pgx
	_, err := r.Pool.Exec(ctx, `
		UPDATE friendships
		SET status = 'blocked'
		WHERE requester_uid = $1 AND receiver_uid = $2`,
		requesterUID, receiverUID,
	)

	return err
}

func (r *PgxFriendshipRepo) ListFriendships(ctx context.Context, uid uuid.UUID) ([]models.Friendship, error) {
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

	var friends []models.Friendship
	for rows.Next() {
		var friend models.Friendship
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
