package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat-posts/models"
)

type PgxPostRepo struct {
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

func (r *PgxPostRepo) CreatePost(ctx context.Context, authorUID uuid.UUID, body string) (*models.Post, error) {
	var p models.Post
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO posts (author_uid, body)
		VALUES ($1, $2)
		RETURNING id, seq, author_uid, body, likes_count, comments_count,
			(EXTRACT(EPOCH FROM created_at) * 1000)::bigint`,
		authorUID, body,
	).Scan(&p.ID, &p.Seq, &p.AuthorUID, &p.Body, &p.LikesCount, &p.CommentsCount, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	// A brand new post can't have been liked yet.
	p.LikedByMe = false
	return &p, nil
}

func (r *PgxPostRepo) GetPost(ctx context.Context, requesterUID, postID uuid.UUID) (*models.Post, error) {
	var p models.Post
	err := r.Pool.QueryRow(ctx, `
		SELECT p.id, p.seq, p.author_uid, p.body, p.likes_count, p.comments_count,
			(EXTRACT(EPOCH FROM p.created_at) * 1000)::bigint,
			(pl.user_id IS NOT NULL) AS liked_by_me
		FROM posts p
		LEFT JOIN post_likes pl ON pl.post_id = p.id AND pl.user_id = $2
		WHERE p.id = $1`,
		postID, requesterUID,
	).Scan(&p.ID, &p.Seq, &p.AuthorUID, &p.Body, &p.LikesCount, &p.CommentsCount, &p.CreatedAt, &p.LikedByMe)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PgxPostRepo) ListFeedLatest(ctx context.Context, requesterUID uuid.UUID, visibleUIDs []uuid.UUID, limit int) ([]models.Post, error) {
	limit = clampLimit(limit)
	rows, err := r.Pool.Query(ctx, `
		SELECT p.id, p.seq, p.author_uid, p.body, p.likes_count, p.comments_count,
			(EXTRACT(EPOCH FROM p.created_at) * 1000)::bigint,
			(pl.user_id IS NOT NULL) AS liked_by_me
		FROM posts p
		LEFT JOIN post_likes pl ON pl.post_id = p.id AND pl.user_id = $1
		WHERE p.author_uid = ANY($2)
		ORDER BY p.seq DESC
		LIMIT $3`,
		requesterUID, visibleUIDs, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

func (r *PgxPostRepo) ListFeedBefore(ctx context.Context, requesterUID uuid.UUID, visibleUIDs []uuid.UUID, beforeSeq int64, limit int) ([]models.Post, error) {
	limit = clampLimit(limit)
	rows, err := r.Pool.Query(ctx, `
		SELECT p.id, p.seq, p.author_uid, p.body, p.likes_count, p.comments_count,
			(EXTRACT(EPOCH FROM p.created_at) * 1000)::bigint,
			(pl.user_id IS NOT NULL) AS liked_by_me
		FROM posts p
		LEFT JOIN post_likes pl ON pl.post_id = p.id AND pl.user_id = $1
		WHERE p.author_uid = ANY($2) AND p.seq < $3
		ORDER BY p.seq DESC
		LIMIT $4`,
		requesterUID, visibleUIDs, beforeSeq, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

func scanPosts(rows pgx.Rows) ([]models.Post, error) {
	posts := make([]models.Post, 0)
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Seq, &p.AuthorUID, &p.Body, &p.LikesCount, &p.CommentsCount, &p.CreatedAt, &p.LikedByMe); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// ToggleLike is a single atomic statement: post_likes' (post_id, user_id)
// primary key makes the delete-or-insert branch race-safe under
// READ COMMITTED, and the counter update happens in the same statement
// so likes_count never drifts from the post_likes row count. Note: a
// true concurrent double-click from the same user could still hit a
// unique-violation on the insert branch (both requests seeing "not yet
// liked" at once) — acceptable as a rare, retryable edge case for this
// project's scale rather than something worth a retry loop.
func (r *PgxPostRepo) ToggleLike(ctx context.Context, postID, userID uuid.UUID) (bool, int, error) {
	var liked bool
	var likesCount int
	err := r.Pool.QueryRow(ctx, `
		WITH del AS (
			DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2
			RETURNING 1
		), ins AS (
			INSERT INTO post_likes (post_id, user_id)
			SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM del)
			RETURNING 1
		), upd AS (
			UPDATE posts
			SET likes_count = likes_count + (SELECT COUNT(*) FROM ins) - (SELECT COUNT(*) FROM del)
			WHERE id = $1
			RETURNING likes_count
		)
		SELECT upd.likes_count, (SELECT COUNT(*) FROM ins) > 0 AS liked
		FROM upd`,
		postID, userID,
	).Scan(&likesCount, &liked)
	if err != nil {
		return false, 0, err
	}
	return liked, likesCount, nil
}

func (r *PgxPostRepo) AddComment(ctx context.Context, postID, authorUID uuid.UUID, body string) (*models.Comment, error) {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var c models.Comment
	err = tx.QueryRow(ctx, `
		INSERT INTO comments (post_id, author_uid, body)
		VALUES ($1, $2, $3)
		RETURNING id, post_id, author_uid, body, (EXTRACT(EPOCH FROM created_at) * 1000)::bigint`,
		postID, authorUID, body,
	).Scan(&c.ID, &c.PostID, &c.AuthorUID, &c.Body, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE posts SET comments_count = comments_count + 1 WHERE id = $1`, postID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PgxPostRepo) ListComments(ctx context.Context, postID uuid.UUID, limit int) ([]models.Comment, error) {
	limit = clampLimit(limit)
	rows, err := r.Pool.Query(ctx, `
		SELECT id, post_id, author_uid, body, (EXTRACT(EPOCH FROM created_at) * 1000)::bigint
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2`,
		postID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorUID, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
