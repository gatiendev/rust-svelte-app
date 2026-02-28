package handler

import (
	"net/http"
	"time"

	"myproject/internal/auth"
	"myproject/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userRepo  repository.UserRepository  // changed from *repository.UserRepo
	tokenRepo repository.TokenRepository // changed from *repository.TokenRepo
}

// Update the constructor accordingly
func NewAuthHandler(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type authResponse struct {
	Message string `json:"message"`
}

// Register handler
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context() // Get context from request

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if user exists
	existing, _ := h.userRepo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user, err := h.userRepo.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Store refresh token
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	if err := h.tokenRepo.StoreRefreshToken(ctx, user.ID, refreshToken, refreshExp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	// Set cookies
	c.SetCookie(
		"access_token",
		accessToken,
		int(15*time.Minute.Seconds()), // MaxAge in seconds
		"/",
		"",
		true, // secure (set to false in dev if needed)
		true, // httpOnly
	)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(7*24*time.Hour.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusCreated, authResponse{Message: "Registration successful"})
}

// Login handler
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !h.userRepo.CheckPassword(ctx, user, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Store refresh token
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	if err := h.tokenRepo.StoreRefreshToken(ctx, user.ID, refreshToken, refreshExp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	// Set cookies
	c.SetCookie(
		"access_token",
		accessToken,
		int(15*time.Minute.Seconds()),
		"/",
		"",
		true,
		true,
	)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(7*24*time.Hour.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, authResponse{Message: "Login successful"})
}

// Refresh handler
func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var refreshToken string
	// Try cookie first
	if cookie, err := c.Cookie("refresh_token"); err == nil {
		refreshToken = cookie
	} else {
		var req refreshRequest
		if err := c.ShouldBindJSON(&req); err == nil {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	// Validate JWT
	if _, err := auth.ValidateRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check DB
	userID, err := h.tokenRepo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token revoked or expired"})
		return
	}

	// Issue new access token
	accessToken, err := auth.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Set new access token cookie
	c.SetCookie(
		"access_token",
		accessToken,
		int(15*time.Minute.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, authResponse{Message: "Token refreshed"})
}

// Logout handler
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	// Revoke all refresh tokens
	if err := h.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// Clear cookies
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, authResponse{Message: "Logout successful"})
}

// Profile handler (example protected endpoint)
func (h *AuthHandler) Profile(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)
	c.String(http.StatusOK, "User ID: %s", userID.String())
}
