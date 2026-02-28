package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"myproject/internal/db"
	"myproject/internal/handler"
	"myproject/internal/repository"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // added
	"github.com/rs/zerolog"
)

func waitForDB(dsn string, timeout time.Duration) error {
	start := time.Now()
	for {
		dbConn, err := sql.Open("postgres", dsn)
		if err == nil {
			err = dbConn.Ping()
			dbConn.Close()
		}
		if err == nil {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for database: %w", err)
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	// Read database configuration from environment (with defaults)
	dbUser := getEnv("DB_USER", "auth_user")
	dbPass := getEnv("DB_PASSWORD", "auth_pass")
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "auth_db")

	// Build DSN for migrations and waiting
	migrationDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Wait for database to be ready
	if err := waitForDB(migrationDSN, 30*time.Second); err != nil {
		logger.Fatal().Err(err).Msg("Database not ready")
	}

	// Get absolute path to migrations directory
	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get working directory")
	}
	migrationsPath := "file://" + filepath.Join(wd, "migrations")

	// Run migrations
	m, err := migrate.New(migrationsPath, migrationDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create migrate instance")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}
	logger.Info().Msg("Migrations applied successfully")

	// Connect to database for application
	pool, err := db.ConnectPool(ctx) // new function
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	// Initialize repositories with pool
	userRepo := repository.NewUserRepo(pool)
	tokenRepo := repository.NewTokenRepo(pool)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userRepo, tokenRepo)

	// Setup routes (returns *gin.Engine)
	router := handler.SetupRoutes(authHandler)

	// Get server port from environment or default
	port := getEnv("SERVER_PORT", "8000")

	logger.Info().Msgf("Server starting on :%s", port)

	// Start the Gin server
	if err := router.Run(":" + port); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}

// Helper to read environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
