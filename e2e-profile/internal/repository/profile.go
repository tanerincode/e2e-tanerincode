package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tanerincode/e2e-profile/internal/model"
	"gorm.io/gorm"
)

// ProfileRepository implements the ProfileRepository interface
type profileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new ProfileRepository instance
func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

// Create adds a new profile to the database
func (r *profileRepository) Create(ctx context.Context, profile *model.ProfileData) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetByID retrieves a profile by its ID
func (r *profileRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ProfileData, error) {
	var profile model.ProfileData
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

// GetByUserID retrieves a profile by user ID
func (r *profileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.ProfileData, error) {
	var profile model.ProfileData
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

// Update updates an existing profile
func (r *profileRepository) Update(ctx context.Context, profile *model.ProfileData) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// Delete removes a profile by its ID
func (r *profileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ProfileData{}, "id = ?", id).Error
}