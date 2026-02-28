package repository

import (
	"crypto/sha256"
	"database/sql"
	"time"

	"myproject/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type TokenRepo struct {
	db *sqlx.DB
}

func NewTokenRepo(db *sqlx.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

// StoreRefreshToken hashes the token and stores it
func (r *TokenRepo) StoreRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	// 1. SHA‑256 the token
	hash := sha256.Sum256([]byte(token))
	// 2. bcrypt the SHA‑256 hash
	hashed, err := bcrypt.GenerateFromPassword(hash[:], bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(query, userID, string(hashed), expiresAt)
	return err
}

// ValidateRefreshToken checks if the token exists and is not revoked, and returns the userID
func (r *TokenRepo) ValidateRefreshToken(token string) (uuid.UUID, error) {
	// Pre‑hash the incoming token the same way
	hash := sha256.Sum256([]byte(token))

	var tokens []models.RefreshToken
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at FROM refresh_tokens 
              WHERE expires_at > NOW() AND revoked_at IS NULL`
	err := r.db.Select(&tokens, query)
	if err != nil {
		return uuid.Nil, err
	}

	for _, t := range tokens {
		// Compare the pre‑hashed token with the stored bcrypt hash
		err = bcrypt.CompareHashAndPassword([]byte(t.TokenHash), hash[:])
		if err == nil {
			return t.UserID, nil
		}
	}
	return uuid.Nil, sql.ErrNoRows
}

// RevokeAllUserTokens revokes all refresh tokens for a user (logout from all devices)
func (r *TokenRepo) RevokeAllUserTokens(userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(query, userID)
	return err
}
