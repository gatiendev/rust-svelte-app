package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myproject/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister_Success(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Mock repositories
	mockUser := &MockUserRepo{}
	mockToken := &MockTokenRepo{}

	userID := uuid.New()
	expectedUser := &models.User{ID: userID, Email: "test@example.com"}

	mockUser.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	mockUser.On("CreateUser", mock.Anything, "test@example.com", "secret").Return(expectedUser, nil)
	mockToken.On("StoreRefreshToken", mock.Anything, userID, mock.Anything, mock.Anything).Return(nil)

	handler := NewAuthHandler(mockUser, mockToken)

	// Create request
	reqBody := registerRequest{Email: "test@example.com", Password: "secret"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp authResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Registration successful", resp.Message)

	mockUser.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUser := &MockUserRepo{}
	mockToken := &MockTokenRepo{}

	existingUser := &models.User{Email: "test@example.com"}
	mockUser.On("GetUserByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

	handler := NewAuthHandler(mockUser, mockToken)

	reqBody := registerRequest{Email: "test@example.com", Password: "secret"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Register(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockUser.AssertExpectations(t)
	mockToken.AssertNotCalled(t, "StoreRefreshToken")
}

// Add tests for Login, Refresh, Logout, Profile similarly.
