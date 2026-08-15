package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leeryan2000/flashat/models"
)

type PgxUserRepo struct {
	Pool *pgxpool.Pool
}

func (r *PgxUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return r.Pool.QueryRow(ctx, `
		INSERT INTO users (
			uid,
			name,
			email,
			hashed_password,
			user_avatar_url
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uid, name, email, hashed_password, user_avatar_url, created_at
	`,
		user.UID,
		user.Name,
		user.Email,
		user.HashedPassword,
		user.UserAvatarURL,
	).Scan(
		&user.UID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.UserAvatarURL,
		&user.CreatedAt,
	)
}

func (r *PgxUserRepo) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT uid, name, email, hashed_password, user_avatar_url, created_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.UID,
			&user.Name,
			&user.Email,
			&user.HashedPassword,
			&user.UserAvatarURL,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PgxUserRepo) GetUserByUID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user := &models.User{}

	err := r.Pool.QueryRow(ctx, `
		SELECT uid, name, email, hashed_password, user_avatar_url, created_at
		FROM users
		WHERE uid = $1
	`, uid).Scan(
		&user.UID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.UserAvatarURL,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PgxUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := r.Pool.QueryRow(ctx, `
		SELECT uid, name, email, hashed_password, user_avatar_url, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.UID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.UserAvatarURL,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PgxUserRepo) UpdateName(ctx context.Context, uid uuid.UUID, name string) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE users
		SET name = $1
		WHERE uid = $2
	`, name, uid)

	return err
}

func (r *PgxUserRepo) UpdateAvatarURL(ctx context.Context, uid uuid.UUID, avatarURL string) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE users
		SET user_avatar_url = $1
		WHERE uid = $2
	`, avatarURL, uid)

	return err
}
