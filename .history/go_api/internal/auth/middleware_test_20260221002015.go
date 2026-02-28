package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "testsecret")
	os.Setenv("ACCESS_TOKEN_DURATION", "5m")
	defer func() {
		os.Unsetenv("JWT_ACCESS_SECRET")
		os.Unsetenv("ACCESS_TOKEN_DURATION")
	}()

	userID := uuid.New()
	token, _ := GenerateAccessToken(userID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: token})

	// Create a dummy next handler
	nextCalled := false
	next := gin.HandlerFunc(func(c *gin.Context) {
		nextCalled = true
		val, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, userID, val)
	})

	// Run middleware
	AuthMiddleware()(c)
	if !c.IsAborted() {
		next(c)
	}

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	AuthMiddleware()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
