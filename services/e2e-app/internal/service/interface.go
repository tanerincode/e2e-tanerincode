package service

import (
	"context"

	"github.com/tanerincode/e2e-app/internal/model"
)

// UserServiceInterface defines the interface for user-related operations
type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
}
