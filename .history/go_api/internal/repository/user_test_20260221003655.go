package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(60 * time.Second),
		NetworkMode: "bridge",
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)
	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	connStr := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, mappedPort.Port())

	// Give it a moment after readiness
	time.Sleep(1 * time.Second)

	config, err := pgxpool.ParseConfig(connStr)
	require.NoError(t, err)
	config.MaxConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)

	// Verify connection with ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = pool.Ping(pingCtx)
	require.NoError(t, err)

	// Run schema (or migrations)
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
	assert.True(t, repo.CheckPassword(fetched, password))
	assert.False(t, repo.CheckPassword(fetched, "wrong"))
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
