package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"myproject/internal/auth"
	"myproject/internal/repository"

	"github.com/google/uuid"
)

type AuthHandler struct {
	userRepo  *repository.UserRepo
	tokenRepo *repository.TokenRepo
}

func NewAuthHandler(userRepo *repository.UserRepo, tokenRepo *repository.TokenRepo) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, tokenRepo: tokenRepo}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"` // optional if sent in cookie
}

type authResponse struct {
	Message string `json:"message"`
}

// setTokenCookies sets http-only cookies for access and refresh tokens
func (h *AuthHandler) setTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	accessExp := time.Now().Add(15 * time.Minute) // match env
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  accessExp,
		HttpOnly: true,
		Secure:   true, // set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshExp,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if user exists
	existing, _ := h.userRepo.GetUserByEmail(req.Email)
	if existing != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	user, err := h.userRepo.CreateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Store refresh token hash in DB
	refreshExp := time.Now().Add(7 * 24 * time.Hour) // from env
	if err := h.tokenRepo.StoreRefreshToken(user.ID, refreshToken, refreshExp); err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	h.setTokenCookies(w, accessToken, refreshToken)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse{Message: "Registration successful"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !h.userRepo.CheckPassword(user, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	if err := h.tokenRepo.StoreRefreshToken(user.ID, refreshToken, refreshExp); err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	h.setTokenCookies(w, accessToken, refreshToken)

	json.NewEncoder(w).Encode(authResponse{Message: "Login successful"})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Try to get refresh token from cookie, fallback to request body
	var refreshToken string
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		refreshToken = cookie.Value
	} else {
		var req refreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken == "" {
		http.Error(w, "Refresh token required", http.StatusBadRequest)
		return
	}

	// Validate the refresh token JWT
	claims, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Check if token exists in DB and not revoked
	userID, err := h.tokenRepo.ValidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Refresh token revoked or expired", http.StatusUnauthorized)
		return
	}

	// Optionally, rotate refresh token: generate new one and revoke old
	// For simplicity, we just issue a new access token
	accessToken, err := auth.GenerateAccessToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Set new access token cookie (refresh token remains the same)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse{Message: "Token refreshed"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Revoke all refresh tokens for this user
	if err := h.tokenRepo.RevokeAllUserTokens(userID); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse{Message: "Logout successful"})
}
