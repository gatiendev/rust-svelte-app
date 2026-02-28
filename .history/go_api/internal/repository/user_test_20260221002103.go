package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections"),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	// Run migrations (simplified – you could also run the actual migration files)
	schema := `
    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );
    CREATE TABLE refresh_tokens (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        token_hash TEXT NOT NULL,
        expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        revoked_at TIMESTAMP WITH TIME ZONE
    );`
	_, err = pool.Exec(ctx, schema)
	require.NoError(t, err)

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}
	return pool, cleanup
}

func TestUserRepo_CreateAndGet(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	repo := NewUserRepo(pool)

	email := "test@example.com"
	password := "secret"
	user, err := repo.CreateUser(ctx, email, password)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, email, user.Email)

	// Fetch by email
	fetched, err := repo.GetUserByEmail(ctx, email)
	require.NoError(t, err)
	assert.Equal(t, user.ID, fetched.ID)
	assert.Equal(t, email, fetched.Email)

	// Check password
	assert.True(t, repo.CheckPassword(ctx, fetched, password))
	assert.False(t, repo.CheckPassword(ctx, fetched, "wrong"))
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	repo := NewUserRepo(pool)
	user, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, user)
}
