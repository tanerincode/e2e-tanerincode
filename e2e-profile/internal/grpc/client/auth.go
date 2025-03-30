package client

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/tanerincode/e2e-profile/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient is a gRPC client for the auth service
type AuthClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

// NewAuthClient creates a new gRPC client for the auth service
func NewAuthClient(address string) (*AuthClient, error) {
	// Use the recommended NewClient method
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		client: pb.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

// ValidateToken validates a JWT token with the auth service
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (bool, string, error) {
	resp, err := c.client.ValidateToken(ctx, &pb.TokenRequest{
		Token: token,
	})
	if err != nil {
		return false, "", fmt.Errorf("failed to validate token: %w", err)
	}

	if !resp.Valid {
		errorMsg := "token invalid"
		if resp.Error != nil {
			errorMsg = resp.Error.Message
		}
		return false, "", errors.New(errorMsg)
	}

	return true, resp.UserId, nil
}

// Close closes the gRPC connection
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
