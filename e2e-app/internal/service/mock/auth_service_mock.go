package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tanerincode/e2e-app/internal/model"
	"github.com/tanerincode/e2e-app/internal/service"
)

// AuthServiceMock is a mock implementation of AuthService
type AuthServiceMock struct {
	mock.Mock
}

// Ensure AuthServiceMock implements AuthService interface
var _ service.AuthService = (*AuthServiceMock)(nil)

// Register mocks the Register method
func (m *AuthServiceMock) Register(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Login mocks the Login method
func (m *AuthServiceMock) Login(ctx context.Context, email, password string) (*model.TokenResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenResponse), args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *AuthServiceMock) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenResponse), args.Error(1)
}