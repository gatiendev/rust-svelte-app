package handler

import (
	"net/http"

	"myproject/internal/auth"
	"myproject/internal/middleware" // <-- new import

	"github.com/google/uuid"
)

func SetupRoutes(authHandler *AuthHandler) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /refresh", authHandler.Refresh)

	// Protected routes
	mux.Handle("POST /logout", auth.AuthMiddleware(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /profile", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(uuid.UUID)
		w.Write([]byte("User ID: " + userID.String()))
	})))

	// Wrap the entire mux with the logger middleware
	return middleware.LoggerMiddleware(mux)
}
