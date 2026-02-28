package repository

import (
	"context"
	"crypto/sha256"
	"time"

	"myproject/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type TokenRepo struct {
	pool *pgxpool.Pool
}

func NewTokenRepo(pool *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{pool: pool}
}

func (r *TokenRepo) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	// SHA‑256 the token
	hash := sha256.Sum256([]byte(token))
	// bcrypt the SHA‑256 hash
	hashed, err := bcrypt.GenerateFromPassword(hash[:], bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err = r.pool.Exec(ctx, query, userID, string(hashed), expiresAt)
	return err
}

func (r *TokenRepo) ValidateRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	hash := sha256.Sum256([]byte(token))

	var tokens []models.RefreshToken
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at FROM refresh_tokens 
              WHERE expires_at > NOW() AND revoked_at IS NULL`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return uuid.Nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.RefreshToken
		err = rows.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt)
		if err != nil {
			return uuid.Nil, err
		}
		err = bcrypt.CompareHashAndPassword([]byte(t.TokenHash), hash[:])
		if err == nil {
			return t.UserID, nil
		}
	}
	return uuid.Nil, pgx.ErrNoRows
}

func (r *TokenRepo) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
