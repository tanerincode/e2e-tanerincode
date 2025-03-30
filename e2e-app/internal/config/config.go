package config

import (
	"os"
	"time"
)

// Config holds the application configuration
type Config struct {
	// HTTP Server
	Port string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret       string
	JWTExpiration   string
	RefreshSecret   string
	RefreshExpiration string
	
	// gRPC
	GRPCPort string
}

// New creates a new Config with values from environment or defaults
func New() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),

		// Database settings
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "e2e_app"),

		// JWT settings
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiration:   getEnv("JWT_EXPIRATION", "24h"),
		RefreshSecret:   getEnv("REFRESH_SECRET", "your-refresh-secret-key"),
		RefreshExpiration: getEnv("REFRESH_EXPIRATION", "168h"),
		
		// gRPC settings
		GRPCPort: getEnv("GRPC_PORT", "50051"),
	}
}

// GetJWTExpiration returns the parsed JWT expiration duration
func (c *Config) GetJWTExpiration() time.Duration {
	duration, err := time.ParseDuration(c.JWTExpiration)
	if err != nil {
		return 24 * time.Hour // Default to 24 hours
	}
	return duration
}

// GetRefreshExpiration returns the parsed refresh token expiration duration
func (c *Config) GetRefreshExpiration() time.Duration {
	duration, err := time.ParseDuration(c.RefreshExpiration)
	if err != nil {
		return 168 * time.Hour // Default to 1 week
	}
	return duration
}

// Helper to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}