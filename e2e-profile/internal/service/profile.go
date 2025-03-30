package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/tanerincode/e2e-profile/internal/config"
	"github.com/tanerincode/e2e-profile/internal/model"
	"github.com/tanerincode/e2e-profile/internal/repository"
)

// ProfileService handles profile data operations
type ProfileService struct {
	client     *http.Client
	config     *config.Config
	profileRepo repository.ProfileRepository
}

// NewProfileService creates a new instance of ProfileService
func NewProfileService(cfg *config.Config, profileRepo repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		client:     &http.Client{},
		config:     cfg,
		profileRepo: profileRepo,
	}
}

// GetProfile retrieves a user profile by ID
// It combines data from our database with data from the auth service
func (s *ProfileService) GetProfile(ctx context.Context, id uuid.UUID) (*model.UserProfile, error) {
	// First, get user data from auth service
	authUserData, err := s.getUserFromAuthService(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}
	
	// Second, get profile data from our database
	profileData, err := s.profileRepo.GetByUserID(ctx, id)
	if err != nil {
		// If profile doesn't exist, that's ok - we'll just return user data
		// with a nil profile
		return &model.UserProfile{
			ID:        id,
			Email:     authUserData.Email,
			FirstName: authUserData.FirstName,
			LastName:  authUserData.LastName,
			CreatedAt: authUserData.CreatedAt,
			UpdatedAt: authUserData.UpdatedAt,
			Profile:   nil,
		}, nil
	}
	
	// Combine data and return
	return &model.UserProfile{
		ID:        id,
		Email:     authUserData.Email,
		FirstName: authUserData.FirstName,
		LastName:  authUserData.LastName,
		CreatedAt: authUserData.CreatedAt,
		UpdatedAt: authUserData.UpdatedAt,
		Profile:   profileData,
	}, nil
}

// CreateProfile creates a new profile
func (s *ProfileService) CreateProfile(ctx context.Context, profile *model.ProfileData) error {
	return s.profileRepo.Create(ctx, profile)
}

// UpdateProfile updates an existing profile
func (s *ProfileService) UpdateProfile(ctx context.Context, profile *model.ProfileData) error {
	return s.profileRepo.Update(ctx, profile)
}

// DeleteProfile deletes a profile
func (s *ProfileService) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	return s.profileRepo.Delete(ctx, id)
}

// Helper method to get user data from auth service
func (s *ProfileService) getUserFromAuthService(ctx context.Context, id uuid.UUID) (*model.UserProfile, error) {
	url := fmt.Sprintf("%s/api/v1/user/profile/%s", s.config.AuthServiceURL, id)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch profile (status %d): %s", resp.StatusCode, string(body))
	}
	
	var profile model.UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &profile, nil
}