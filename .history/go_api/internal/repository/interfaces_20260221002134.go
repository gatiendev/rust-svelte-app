package repository

import (
	"context"
	"time"

	"myproject/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CheckPassword(user *models.User, password string) bool
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (uuid.UUID, error)
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
}
