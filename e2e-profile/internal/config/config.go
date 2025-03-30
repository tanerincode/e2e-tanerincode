package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	AuthServiceURL string
	AuthGRPCAddr   string
	Port           string
	MockDB         bool

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// New creates a new Config with values from environment or defaults
func New() *Config {
	return &Config{
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
		AuthGRPCAddr:   getEnv("AUTH_GRPC_ADDR", "localhost:50051"),
		Port:           getEnv("PORT", "8081"),
		MockDB:         getBoolEnv("MOCK_DB", false),

		// Database defaults - typically overridden by environment in production
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "e2e_profile"),
	}
}

// Helper to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper to get boolean environment variable with fallback
func getBoolEnv(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return fallback
}
