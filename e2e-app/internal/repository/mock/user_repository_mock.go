package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/tanerincode/e2e-app/internal/model"
	"github.com/tanerincode/e2e-app/internal/repository"
)

// UserRepositoryMock is a mock implementation of UserRepository
type UserRepositoryMock struct {
	mock.Mock
}

// Ensure UserRepositoryMock implements UserRepository interface
var _ repository.UserRepository = (*UserRepositoryMock)(nil)

// Create mocks the Create method
func (m *UserRepositoryMock) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *UserRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// GetByEmail mocks the GetByEmail method
func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Update mocks the Update method
func (m *UserRepositoryMock) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *UserRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}