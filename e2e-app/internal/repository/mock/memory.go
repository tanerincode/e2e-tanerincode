package mock

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tanerincode/e2e-app/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// MemoryUserRepository provides a memory-based mock implementation of UserRepository
type MemoryUserRepository struct {
	users      map[uuid.UUID]*model.User
	emailIndex map[string]uuid.UUID
	mu         sync.RWMutex
}

// NewMemoryUserRepository creates a new memory-based mock user repository
func NewMemoryUserRepository() UserRepository {
	return &MemoryUserRepository{
		users:      make(map[uuid.UUID]*model.User),
		emailIndex: make(map[string]uuid.UUID),
	}
}

// GetMockUserRepository returns a memory repository that can be used for testing/mocking
func GetMockUserRepository() UserRepository {
	return NewMemoryUserRepository()
}

// Create creates a new user in the mock repository
func (r *MemoryUserRepository) Create(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	if _, exists := r.emailIndex[user.Email]; exists {
		return errors.New("email already exists")
	}

	// If ID is not set, generate one
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Set timestamps if not set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	// Hash password if it's not already hashed
	if len(user.Password) > 0 && user.Password[0] != '$' {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	// Create a deep copy of the user to store
	userCopy := &model.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Store user in memory
	r.users[user.ID] = userCopy
	r.emailIndex[user.Email] = user.ID

	// Update the input user with the generated ID
	user.ID = userCopy.ID

	log.Printf("User created with ID: %s, Email: %s", user.ID, user.Email)
	return nil
}

// GetByID retrieves a user by their ID from the mock repository
func (r *MemoryUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Return a copy to prevent mutation of the stored user
	userCopy := &model.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userCopy, nil
}

// GetByEmail retrieves a user by their email from the mock repository
func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.emailIndex[email]
	if !exists {
		log.Printf("User with email %s not found", email)
		return nil, errors.New("user not found")
	}

	user := r.users[id]

	// Return a copy to prevent mutation of the stored user
	userCopy := &model.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	log.Printf("User retrieved by email: %s, ID: %s", email, user.ID)
	return userCopy, nil
}

// Update updates user information in the mock repository
func (r *MemoryUserRepository) Update(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Create a deep copy of the user
	userCopy := &model.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// If email has changed, update the index
	currentUser := r.users[user.ID]
	if currentUser.Email != user.Email {
		delete(r.emailIndex, currentUser.Email)
		r.emailIndex[user.Email] = user.ID
	}

	// Update user in memory
	r.users[user.ID] = userCopy

	return nil
}

// Delete removes a user from the mock repository
func (r *MemoryUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}

	// Remove from indices
	delete(r.emailIndex, user.Email)
	delete(r.users, id)

	return nil
}
