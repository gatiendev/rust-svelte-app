package handler

import (
	"myproject/internal/auth"
	"myproject/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(authHandler *AuthHandler) *gin.Engine {
	router := gin.New() // No default middleware; we'll add our own

	// Add custom request logger (or use gin's Logger middleware)
	router.Use(middleware.LoggerMiddleware(), gin.Recovery())

	// Public routes
	public := router.Group("/")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.Refresh)
	}

	// Protected routes
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware()) // Gin middleware
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/profile", authHandler.Profile) // Move profile handler to a method
	}

	return router
}
