package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanerincode/e2e-profile/internal/model"
)

// ProfileServiceInterface defines the interface for profile-related operations
type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*model.UserProfile, error)
	CreateProfile(ctx context.Context, profile *model.ProfileData) error
	UpdateProfile(ctx context.Context, profile *model.ProfileData) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error
}