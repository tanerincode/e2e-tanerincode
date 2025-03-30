package server

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tanerincode/e2e-app/internal/config"
	pb "github.com/tanerincode/e2e-app/internal/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer implements the gRPC auth service for token validation
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	config *config.Config
}

// NewAuthServer creates a new auth gRPC server
func NewAuthServer(cfg *config.Config) *AuthServer {
	return &AuthServer{
		config: cfg,
	}
}

// ValidateToken validates a JWT token and returns user information
func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	if req.Token == "" {
		return &pb.TokenResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    "invalid_token",
				Message: "Token is empty",
			},
		}, nil
	}

	// Parse token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return &pb.TokenResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    "token_parsing_failed",
				Message: err.Error(),
			},
		}, nil
	}

	// Validate token
	if !token.Valid {
		return &pb.TokenResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    "invalid_token",
				Message: "Token is invalid",
			},
		}, nil
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to parse token claims")
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

	// Get email (if available)
	email, _ := claims["email"].(string)

	return &pb.TokenResponse{
		Valid:  true,
		UserId: userID,
		Email:  email,
	}, nil
}
