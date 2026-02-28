package handler

import (
	"context"
	"time"

	"myproject/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) CheckPassword(user *models.User, password string) bool {
	args := m.Called(user, password)
	return args.Bool(0)
}

type MockTokenRepo struct {
	mock.Mock
}

func (m *MockTokenRepo) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockTokenRepo) ValidateRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockTokenRepo) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
