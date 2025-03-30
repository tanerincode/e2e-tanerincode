package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanerincode/e2e-profile/internal/model"
)

// ProfileRepository defines the interface for profile data operations
type ProfileRepository interface {
	Create(ctx context.Context, profile *model.ProfileData) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.ProfileData, error)
	Update(ctx context.Context, profile *model.ProfileData) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.ProfileData, error)
}