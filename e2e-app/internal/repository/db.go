package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/tanerincode/e2e-app/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// IsMockDBEnabled returns whether the mock DB is enabled
func IsMockDBEnabled(cfg *config.Config) bool {
	return cfg.AppEnv == "development" && cfg.MockDB
}

// NewDB creates a new database connection
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	// Check if we're in dev mode with mock DB
	if IsMockDBEnabled(cfg) {
		log.Println("Using in-memory mock database for development")
		return nil, nil // Return nil to trigger mock implementations
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	var db *gorm.DB
	var err error
	maxRetries := 5

	// Attempt to connect with retries
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}

		retryAfter := time.Duration(i+1) * time.Second
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...",
			i+1, maxRetries, err, retryAfter)
		time.Sleep(retryAfter)
	}

	if err != nil {
		// If in development mode, we can fall back to mock DB even if mock wasn't explicitly enabled
		if cfg.AppEnv == "development" {
			log.Println("Failed to connect to database in development mode. Falling back to mock database")
			return nil, nil // Return nil to trigger mock implementations
		}
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	return db, nil
}
