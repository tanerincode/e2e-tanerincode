package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tanerincode/e2e-profile/internal/config"
	"github.com/tanerincode/e2e-profile/internal/grpc/client"
	"github.com/tanerincode/e2e-profile/internal/handler"
	"github.com/tanerincode/e2e-profile/internal/repository"
	"github.com/tanerincode/e2e-profile/internal/service"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	profileRepo := repository.NewProfileRepository(db)

	// Initialize gRPC client
	authClient, err := client.NewAuthClient(cfg.AuthGRPCAddr)
	if err != nil {
		log.Fatalf("Failed to connect to auth gRPC service: %v", err)
	}
	defer authClient.Close()

	// Initialize services
	profileService := service.NewProfileService(cfg, profileRepo)

	// Initialize handlers
	profileHandler := handler.NewProfileHandler(profileService)

	// Setup router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		profiles := api.Group("/profiles")
		{
			profiles.GET("/:id", profileHandler.GetProfile)
		}

		// Protected routes
		protected := api.Group("/profiles")
		protected.Use(handler.AuthMiddleware(cfg, authClient))
		{
			protected.POST("/", profileHandler.CreateProfile)
			protected.PUT("/:id", profileHandler.UpdateProfile)
			protected.DELETE("/:id", profileHandler.DeleteProfile)
		}
	}

	// Start server
	log.Printf("HTTP server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}