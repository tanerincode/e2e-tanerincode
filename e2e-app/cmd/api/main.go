package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tanerincode/e2e-app/internal/config"
	auth "github.com/tanerincode/e2e-app/internal/grpc/proto"
	"github.com/tanerincode/e2e-app/internal/grpc/server"
	"github.com/tanerincode/e2e-app/internal/handler"
	"github.com/tanerincode/e2e-app/internal/repository"
	"github.com/tanerincode/e2e-app/internal/service"
	"google.golang.org/grpc"
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
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Create a wait group to wait for both servers
	var wg sync.WaitGroup
	wg.Add(2)

	// Start gRPC server in a separate goroutine
	go startGRPCServer(&wg, cfg)

	// Start HTTP server in a separate goroutine
	go startHTTPServer(&wg, cfg, authHandler, userHandler)

	// Wait for both servers to finish (which should never happen)
	wg.Wait()
}

// startGRPCServer starts the gRPC server
func startGRPCServer(wg *sync.WaitGroup, cfg *config.Config) {
	defer wg.Done()

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register auth service
	authServer := server.NewAuthServer(cfg)
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	log.Printf("gRPC server starting on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

// startHTTPServer starts the HTTP server
func startHTTPServer(wg *sync.WaitGroup, cfg *config.Config, authHandler *handler.AuthHandler, userHandler *handler.UserHandler) {
	defer wg.Done()

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
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("/user")
		protected.Use(handler.AuthMiddleware(cfg))
		{
			protected.GET("/profile", userHandler.GetProfile)
		}
	}

	// Start server
	log.Printf("HTTP server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
