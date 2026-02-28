package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/auth-system/internal/models"
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
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(query, userID, string(hashed), expiresAt)
	return err
}

// ValidateRefreshToken checks if the token exists and is not revoked, and returns the userID
func (r *TokenRepo) ValidateRefreshToken(token string) (uuid.UUID, error) {
	// Since tokens are hashed, we need to scan all active tokens and compare
	var tokens []models.RefreshToken
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at FROM refresh_tokens 
	          WHERE expires_at > NOW() AND revoked_at IS NULL`
	err := r.db.Select(&tokens, query)
	if err != nil {
		return uuid.Nil, err
	}

	for _, t := range tokens {
		err = bcrypt.CompareHashAndPassword([]byte(t.TokenHash), []byte(token))
		if err == nil {
			// Found matching token
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
