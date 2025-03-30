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