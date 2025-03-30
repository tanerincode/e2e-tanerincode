package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tanerincode/e2e-app/internal/model"
	"github.com/tanerincode/e2e-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for user-related operations
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser handles user creation with password hashing
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	if user.Password == "" {
		return errors.New("password is required")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Create(ctx, user)
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail retrieves a user by their email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// AuthenticateUser verifies user credentials
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
