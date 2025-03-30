package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tanerincode/e2e-profile/internal/config"
	"github.com/tanerincode/e2e-profile/internal/grpc/client"
)

// AuthMiddleware validates tokens via gRPC with the auth service
func AuthMiddleware(cfg *config.Config, authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Check format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token via gRPC
		valid, userID, err := authClient.ValidateToken(context.Background(), token)
		if err != nil || !valid {
			errMsg := "invalid token"
			if err != nil {
				errMsg = err.Error()
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", userID)
		c.Next()
	}
}