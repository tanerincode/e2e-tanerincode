package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanerincode/e2e-profile/internal/model"
	"github.com/tanerincode/e2e-profile/internal/service"
)

// ProfileHandler handles HTTP requests for profile operations
type ProfileHandler struct {
	profileService service.ProfileServiceInterface
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(profileService service.ProfileServiceInterface) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile retrieves a user profile by ID
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	idParam := c.Param("id")
	
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}
	
	profile, err := h.profileService.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, profile)
}

// CreateProfile creates a new profile
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var profile model.ProfileData
	
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Ensure valid UUID
	if profile.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	
	if err := h.profileService.CreateProfile(c.Request.Context(), &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, profile)
}

// UpdateProfile updates an existing profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	idParam := c.Param("id")
	
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile ID format"})
		return
	}
	
	var profile model.ProfileData
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Ensure ID in path matches body
	profile.ID = id
	
	if err := h.profileService.UpdateProfile(c.Request.Context(), &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, profile)
}

// DeleteProfile deletes a profile
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	idParam := c.Param("id")
	
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile ID format"})
		return
	}
	
	if err := h.profileService.DeleteProfile(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "profile deleted successfully"})
}