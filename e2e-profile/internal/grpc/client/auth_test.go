package client_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tanerincode/e2e-profile/internal/grpc/client"
	pb "github.com/tanerincode/e2e-profile/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// mockAuthServer implements the AuthService gRPC server for testing
type mockAuthServer struct {
	pb.UnimplementedAuthServiceServer
	jwtSecret string
}

// ValidateToken implementation for the mock server
func (s *mockAuthServer) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	if req.Token == "" {
		return &pb.TokenResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    "invalid_token",
				Message: "Token is empty",
			},
		}, nil
	}

	// Parse and validate the token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return &pb.TokenResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    "invalid_token",
				Message: "Token is invalid or expired",
			},
		}, nil
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &pb.TokenResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    "invalid_claims",
				Message: "Failed to parse token claims",
			},
		}, nil
	}

	// Get user ID
	userID, ok := claims["user_id"].(string)
	if !ok {
		return &pb.TokenResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    "missing_user_id",
				Message: "User ID not found in token",
			},
		}, nil
	}

	return &pb.TokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}

// setupGRPCServer creates an in-memory gRPC server for testing
func setupGRPCServer(jwtSecret string) (*bufconn.Listener, *grpc.Server) {
	bufListener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	// Register the mock auth server
	pb.RegisterAuthServiceServer(grpcServer, &mockAuthServer{jwtSecret: jwtSecret})

	// Start the server in background
	go func() {
		if err := grpcServer.Serve(bufListener); err != nil {
			panic(err)
		}
	}()

	return bufListener, grpcServer
}

// getBufDialer returns a dialer for the bufconn listener
func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// TestAuthClient tests the gRPC client for auth validation
func TestAuthClient(t *testing.T) {
	// Setup
	jwtSecret := "test-secret-key"
	userID := uuid.New().String()

	// Create a valid JWT token
	validTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validTokenClaims)
	validTokenString, err := validToken.SignedString([]byte(jwtSecret))
	require.NoError(t, err)

	// Create an expired token
	expiredTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredTokenClaims)
	expiredTokenString, err := expiredToken.SignedString([]byte(jwtSecret))
	require.NoError(t, err)

	// Setup in-memory gRPC server
	listener, server := setupGRPCServer(jwtSecret)
	defer server.Stop()

	// Setup gRPC connection
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", 
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	// Create a gRPC service client
	grpcClient := pb.NewAuthServiceClient(conn)

	tests := []struct {
		name          string
		token         string
		expectedValid bool
		expectedUserID string
		expectedError bool
	}{
		{
			name:          "valid token",
			token:         validTokenString,
			expectedValid: true,
			expectedUserID: userID,
			expectedError: false,
		},
		{
			name:          "expired token",
			token:         expiredTokenString,
			expectedValid: false,
			expectedUserID: "",
			expectedError: true,
		},
		{
			name:          "empty token",
			token:         "",
			expectedValid: false,
			expectedUserID: "",
			expectedError: true,
		},
		{
			name:          "invalid token format",
			token:         "invalid-token",
			expectedValid: false,
			expectedUserID: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test auth client with the exported fields
			authClient := &client.AuthClient{}
			
			// Use reflection to set private fields for testing
			// This is a workaround for testing private fields
			// In real-world applications, consider creating a factory function for tests
			authClient.SetClientForTesting(grpcClient, conn)

			// Call ValidateToken
			valid, userID, err := authClient.ValidateToken(ctx, tt.token)

			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, valid)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.True(t, valid)
				assert.Equal(t, tt.expectedUserID, userID)
			}
		})
	}
}