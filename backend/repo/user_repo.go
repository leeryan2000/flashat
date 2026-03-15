package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByUID(ctx context.Context, uid uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateName(ctx context.Context, uid uuid.UUID, name string) error
}
