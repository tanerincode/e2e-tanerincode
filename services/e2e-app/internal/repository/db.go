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

// NewDB creates a new database connection
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	// Check if we're in dev mode with mock DB
	if cfg.AppEnv == "development" && cfg.MockDB {
		log.Println("Using in-memory mock database for development")
		return setupMockDB()
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
			return setupMockDB()
		}
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	return db, nil
}

// setupMockDB creates an in-memory mock database for development
func setupMockDB() (*gorm.DB, error) {
	// For a real implementation, you would use an SQLite in-memory database or a mock implementation
	// For now, we'll just return a stub with a message
	log.Println("Mock database not fully implemented. This is a placeholder!")

	// Placeholder for the actual implementation
	// In a real implementation, you'd return something like:
	// return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	return nil, fmt.Errorf("mock database implementation pending")
}
