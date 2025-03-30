package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tanerincode/e2e-app/internal/config"
	pb "github.com/tanerincode/e2e-app/internal/grpc/proto"
	"github.com/tanerincode/e2e-app/internal/grpc/server"
)

func TestValidateToken(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Create a new auth server
	authServer := server.NewAuthServer(cfg)

	// Generate a test user ID
	userID := uuid.New().String()

	// Create a valid JWT token for testing
	validTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   "test@example.com",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validTokenClaims)
	validTokenString, _ := validToken.SignedString([]byte(cfg.JWTSecret))

	// Create an expired token for testing
	expiredTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredTokenClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(cfg.JWTSecret))

	// Create a token with invalid signature
	invalidSignatureToken := validTokenString + "invalid"

	// Create a token without user_id
	missingUserIDClaims := jwt.MapClaims{
		"email": "test@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	missingUserIDToken := jwt.NewWithClaims(jwt.SigningMethodHS256, missingUserIDClaims)
	missingUserIDTokenString, _ := missingUserIDToken.SignedString([]byte(cfg.JWTSecret))

	// Test cases
	tests := []struct {
		name           string
		token          string
		expectedValid  bool
		expectedUserID string
		expectedEmail  string
		expectedError  bool
	}{
		{
			name:           "valid token",
			token:          validTokenString,
			expectedValid:  true,
			expectedUserID: userID,
			expectedEmail:  "test@example.com",
			expectedError:  false,
		},
		{
			name:           "expired token",
			token:          expiredTokenString,
			expectedValid:  false,
			expectedUserID: "",
			expectedEmail:  "",
			expectedError:  true,
		},
		{
			name:           "invalid signature",
			token:          invalidSignatureToken,
			expectedValid:  false,
			expectedUserID: "",
			expectedEmail:  "",
			expectedError:  true,
		},
		{
			name:           "empty token",
			token:          "",
			expectedValid:  false,
			expectedUserID: "",
			expectedEmail:  "",
			expectedError:  true,
		},
		{
			name:           "missing user_id",
			token:          missingUserIDTokenString,
			expectedValid:  false,
			expectedUserID: "",
			expectedEmail:  "",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := &pb.TokenRequest{
				Token: tt.token,
			}

			// Call ValidateToken
			resp, err := authServer.ValidateToken(context.Background(), req)

			// All calls should return without gRPC errors
			assert.NoError(t, err)
			assert.NotNil(t, resp)

			// Check response validity
			assert.Equal(t, tt.expectedValid, resp.Valid)

			if tt.expectedValid {
				// For valid tokens, check user details
				assert.Equal(t, tt.expectedUserID, resp.UserId)
				assert.Equal(t, tt.expectedEmail, resp.Email)
				assert.Nil(t, resp.Error)
			} else {
				// For invalid tokens, check error details
				assert.Empty(t, resp.UserId)
				if tt.expectedError {
					assert.NotNil(t, resp.Error)
					assert.NotEmpty(t, resp.Error.Code)
					assert.NotEmpty(t, resp.Error.Message)
				}
			}
		})
	}
}