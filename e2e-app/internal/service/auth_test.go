package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tanerincode/e2e-app/internal/config"
	"github.com/tanerincode/e2e-app/internal/model"
	repoMock "github.com/tanerincode/e2e-app/internal/repository/mock"
	"github.com/tanerincode/e2e-app/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name     string
		user     *model.User
		mockFunc func(*repoMock.UserRepositoryMock)
		wantErr  bool
	}{
		{
			name: "successful registration",
			user: &model.User{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			user: &model.User{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(errors.New("repository error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the mock repository
			mockRepo := new(repoMock.UserRepositoryMock)
			tt.mockFunc(mockRepo)

			// Create test configuration
			cfg := &config.Config{
				JWTSecret:         "test-secret",
				JWTExpiration:     "1h",
				RefreshExpiration: "24h",
			}

			// Initialize the auth service with mock repo
			authService := service.NewAuthService(mockRepo, cfg)

			// Call Register method
			err := authService.Register(context.Background(), tt.user)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Assert expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	// Hash a test password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	testUser := &model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name     string
		email    string
		password string
		mockFunc func(*repoMock.UserRepositoryMock)
		wantErr  bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				m.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, errors.New("user not found"))
			},
			wantErr: true,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the mock repository
			mockRepo := new(repoMock.UserRepositoryMock)
			tt.mockFunc(mockRepo)

			// Create test configuration
			cfg := &config.Config{
				JWTSecret:         "test-secret",
				JWTExpiration:     "1h",
				RefreshExpiration: "24h",
			}

			// Initialize the auth service with mock repo
			authService := service.NewAuthService(mockRepo, cfg)

			// Call Login method
			response, err := authService.Login(context.Background(), tt.email, tt.password)

			// Check error and response
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, "Bearer", response.TokenType)
				assert.Greater(t, response.ExpiresIn, int64(0))
			}

			// Assert expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	// Create a test user
	testUser := &model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create a valid JWT token for testing
	tokenClaims := jwt.MapClaims{
		"user_id": testUser.ID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	cfg := &config.Config{
		JWTSecret:         "test-secret",
		JWTExpiration:     "1h",
		RefreshExpiration: "24h",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	validToken, _ := token.SignedString([]byte(cfg.JWTSecret))

	// Create an expired token for testing
	expiredClaims := jwt.MapClaims{
		"user_id": testUser.ID.String(),
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(cfg.JWTSecret))

	tests := []struct {
		name        string
		refreshToken string
		mockFunc    func(*repoMock.UserRepositoryMock)
		wantErr     bool
	}{
		{
			name:        "successful refresh",
			refreshToken: validToken,
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				m.On("GetByID", mock.Anything, testUser.ID).Return(testUser, nil)
			},
			wantErr: false,
		},
		{
			name:        "expired token",
			refreshToken: expiredTokenString,
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				// No mock calls since validation will fail before repository call
			},
			wantErr: true,
		},
		{
			name:        "user not found",
			refreshToken: validToken,
			mockFunc: func(m *repoMock.UserRepositoryMock) {
				m.On("GetByID", mock.Anything, testUser.ID).Return(nil, errors.New("user not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the mock repository
			mockRepo := new(repoMock.UserRepositoryMock)
			tt.mockFunc(mockRepo)

			// Initialize the auth service with mock repo
			authService := service.NewAuthService(mockRepo, cfg)

			// Call RefreshToken method
			response, err := authService.RefreshToken(context.Background(), tt.refreshToken)

			// Check error and response
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, "Bearer", response.TokenType)
				assert.Greater(t, response.ExpiresIn, int64(0))
			}

			// Assert expectations
			mockRepo.AssertExpectations(t)
		})
	}
}