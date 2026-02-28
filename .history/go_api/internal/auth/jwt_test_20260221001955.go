package auth

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	// Set environment (or use mock secrets)
	os.Setenv("JWT_ACCESS_SECRET", "testsecret")
	os.Setenv("ACCESS_TOKEN_DURATION", "5m")
	defer os.Unsetenv("JWT_ACCESS_SECRET")
	defer os.Unsetenv("ACCESS_TOKEN_DURATION")

	userID := uuid.New()
	token, err := GenerateAccessToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), claims.ExpiresAt.Time, time.Second)
}

func TestGenerateAndValidateRefreshToken(t *testing.T) {
	os.Setenv("JWT_REFRESH_SECRET", "testrefresh")
	os.Setenv("REFRESH_TOKEN_DURATION", "72h")
	defer os.Unsetenv("JWT_REFRESH_SECRET")
	defer os.Unsetenv("REFRESH_TOKEN_DURATION")

	userID := uuid.New()
	token, err := GenerateRefreshToken(userID)
	require.NoError(t, err)

	claims, err := ValidateRefreshToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateAccessToken_Expired(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "testsecret")
	os.Setenv("ACCESS_TOKEN_DURATION", "-1m") // expired
	defer os.Unsetenv("JWT_ACCESS_SECRET")
	defer os.Unsetenv("ACCESS_TOKEN_DURATION")

	userID := uuid.New()
	token, err := GenerateAccessToken(userID)
	require.NoError(t, err)

	_, err = ValidateAccessToken(token)
	assert.Error(t, err)
}
