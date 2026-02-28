package main

import (
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
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
)

func waitForDB(dsn string, timeout time.Duration) error {
	start := time.Now()
	for {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			db.Close()
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

	// Database connection string for migrations and waiting
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" {
		dbUser = "auth_user"
	}
	if dbPass == "" {
		dbPass = "auth_pass"
	}
	if dbHost == "" {
		dbHost = "postgres"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbName == "" {
		dbName = "auth_db"
	}

	migrationDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Wait for database
	if err := waitForDB(migrationDSN, 30*time.Second); err != nil {
		logger.Fatal().Err(err).Msg("Database not ready")
	}

	// Get absolute migrations path
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
	database, err := db.Connect()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepo(database)
	tokenRepo := repository.NewTokenRepo(database)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userRepo, tokenRepo)

	// Setup routes (returns *gin.Engine)
	router := handler.SetupRoutes(authHandler)

	// Get server port from environment or default
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	logger.Info().Msgf("Server starting on :%s", port)

	// Start the Gin server
	if err := router.Run(":" + port); err != nil {
		logger.Fatal().Err(err).Msg("Server failed")
	}
}
