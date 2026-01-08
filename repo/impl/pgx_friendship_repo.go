package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxFriendshipRepo struct {
	Pool *pgxpool.Pool
}

func (r *PgxFriendshipRepo) RequestFriendship(ctx context.Context, requesterUID uuid.UUID, email string) error {
	// Implementation for requesting friendship using pgx
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO friendships (requester_uid, receiver_uid, status)
		SELECT
			$1,
			u.uid,
			'pending'
		FROM users u
		WHERE u.email = $2
		AND NOT EXISTS (
			SELECT 1 FROM friendships f
			WHERE (f.requester_uid = $1 AND f.receiver_uid = u.uid)
			   OR (f.requester_uid = u.uid AND f.receiver_uid = $1)
		)
		`,
		requesterUID, email,
	)

	return err
}

func (r *PgxFriendshipRepo) AcceptFriendship(ctx context.Context, requesterUID, receiverUID uuid.UUID) error {
	// Implementation for accepting friendship using pgx
	_, err := r.Pool.Exec(ctx, `
		UPDATE friendships
		SET status = 'accepted'
		WHERE requester_uid = $1 AND receiver_uid = $2 AND status = 'pending'`,
		requesterUID, receiverUID,
	)

	return err
}

func (r *PgxFriendshipRepo) DeleteFriendship(ctx context.Context, requesterUID, receiverUID uuid.UUID) error {
	// Implementation for deleting friendship using pgx
	_, err := r.Pool.Exec(ctx, `
		DELETE FROM friendships
		WHERE (requester_uid = $1 AND receiver_uid = $2)
		   OR (requester_uid = $2 AND receiver_uid = $1)`,
		requesterUID, receiverUID,
	)
	return err
}

func (r *PgxFriendshipRepo) BlockUser(ctx context.Context, requesterUID, receiverUID uuid.UUID) error {
	// Implementation for blocking user using pgx
	return nil
}

func (r *PgxFriendshipRepo) ListFriendships(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	// Implementation for listing friendships using pgx
	return nil, nil
}

func (r *PgxFriendshipRepo) GetFriendshipStatus(ctx context.Context, userAUID, userBUID uuid.UUID) (string, error) {
	// Implementation for getting friendship status using pgx
	return "", nil
}
