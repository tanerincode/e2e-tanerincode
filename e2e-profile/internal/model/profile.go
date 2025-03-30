package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProfileData represents the profile information stored in our database
type ProfileData struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	Bio         string         `gorm:"type:text" json:"bio"`
	Avatar      string         `gorm:"type:varchar(255)" json:"avatar"`
	Interests   []string       `gorm:"type:text[]" json:"interests"`
	SocialLinks map[string]string `gorm:"-" json:"social_links"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserProfile combines user data from auth service with profile data
type UserProfile struct {
	ID        uuid.UUID         `json:"id"`
	Email     string            `json:"email"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Profile   *ProfileData      `json:"profile"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}