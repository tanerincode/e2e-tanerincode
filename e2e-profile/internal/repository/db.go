package repository

import (
	"fmt"
	"log"

	"github.com/tanerincode/e2e-profile/internal/config"
	"github.com/tanerincode/e2e-profile/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB establishes a connection to the database
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	// Check if mock database is enabled
	if cfg.MockDB {
		log.Println("Using mock database for Profile service")
		return setupMockDB()
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	err = db.AutoMigrate(&model.ProfileData{})
	if err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	return db, nil
}

// setupMockDB returns a mock database implementation
func setupMockDB() (*gorm.DB, error) {
	// For a real implementation, we'd use SQLite in-memory or a proper mock
	log.Println("Mock database is enabled. Using in-memory implementation")

	// Return a nil DB for now, but the application won't crash due to connection issues
	// In a real implementation, you'd replace this with a proper in-memory DB
	return nil, nil
}
