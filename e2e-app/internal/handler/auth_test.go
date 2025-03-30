package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tanerincode/e2e-app/internal/handler"
	"github.com/tanerincode/e2e-app/internal/model"
	serviceMock "github.com/tanerincode/e2e-app/internal/service/mock"
)

func TestRegisterHandler(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Test cases
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(*serviceMock.AuthServiceMock)
		expectedStatus int
		expectedJSON   bool
	}{
		{
			name: "successful registration",
			requestBody: model.RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// Configure mock to return success for Register
				m.On("Register", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedJSON:   true,
		},
		{
			name: "service error",
			requestBody: model.RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// Configure mock to return an error for Register
				m.On("Register", mock.Anything, mock.AnythingOfType("*model.User")).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   true,
		},
		{
			name: "invalid request",
			requestBody: map[string]interface{}{
				"email":      "invalid-email",
				"password":   "123", // Too short
				"first_name": "",    // Empty
				"last_name":  "",    // Empty
			},
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// No mock calls expected for invalid request
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock auth service
			mockService := new(serviceMock.AuthServiceMock)
			tt.mockFunc(mockService)

			// Create handler with the mock service
			authHandler := handler.NewAuthHandler(mockService)

			// Create a new Gin router
			router := gin.New()

			// Register the handler
			router.POST("/register", authHandler.Register)

			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Record the response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check JSON response if expected
			if tt.expectedJSON {
				assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
			}

			// Assert all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a sample token response
	tokenResponse := &model.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	// Test cases
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(*serviceMock.AuthServiceMock)
		expectedStatus int
		expectedJSON   bool
	}{
		{
			name: "successful login",
			requestBody: model.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// Configure mock to return token response for successful login
				m.On("Login", mock.Anything, "test@example.com", "password123").Return(tokenResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   true,
		},
		{
			name: "service error",
			requestBody: model.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// Configure mock to return an error for Login
				m.On("Login", mock.Anything, "test@example.com", "password123").Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedJSON:   true,
		},
		{
			name: "invalid request",
			requestBody: map[string]interface{}{
				"email":    "invalid-email",
				"password": "123", // Too short
			},
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// No mock calls expected for invalid request
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock auth service
			mockService := new(serviceMock.AuthServiceMock)
			tt.mockFunc(mockService)

			// Create handler with the mock service
			authHandler := handler.NewAuthHandler(mockService)

			// Create a new Gin router
			router := gin.New()

			// Register the handler
			router.POST("/login", authHandler.Login)

			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Record the response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check JSON response if expected
			if tt.expectedJSON {
				assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
			}

			// For successful login, verify token in response
			if tt.expectedStatus == http.StatusOK {
				var response model.TokenResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tokenResponse.AccessToken, response.AccessToken)
				assert.Equal(t, tokenResponse.RefreshToken, response.RefreshToken)
				assert.Equal(t, tokenResponse.TokenType, response.TokenType)
				assert.Equal(t, tokenResponse.ExpiresIn, response.ExpiresIn)
			}

			// Assert all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestRefreshTokenHandler(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a sample token response for refresh
	tokenResponse := &model.TokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	// Test cases
	tests := []struct {
		name           string
		refreshToken   string
		mockFunc       func(*serviceMock.AuthServiceMock)
		expectedStatus int
		expectedJSON   bool
	}{
		{
			name:         "successful refresh",
			refreshToken: "valid-refresh-token",
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// Configure mock to return token response for successful refresh
				m.On("RefreshToken", mock.Anything, "valid-refresh-token").Return(tokenResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   true,
		},
		{
			name:         "invalid token",
			refreshToken: "invalid-token",
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// Configure mock to return error for invalid token
				m.On("RefreshToken", mock.Anything, "invalid-token").Return(nil, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedJSON:   true,
		},
		{
			name:         "expired token",
			refreshToken: "expired-token",
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// Configure mock to return error for expired token
				m.On("RefreshToken", mock.Anything, "expired-token").Return(nil, errors.New("token expired"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedJSON:   true,
		},
		{
			name:         "missing token",
			refreshToken: "",
			mockFunc: func(m *serviceMock.AuthServiceMock) {
				// No mock calls expected for missing token
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock auth service
			mockService := new(serviceMock.AuthServiceMock)
			tt.mockFunc(mockService)

			// Create handler with the mock service
			authHandler := handler.NewAuthHandler(mockService)

			// Create a new Gin router
			router := gin.New()

			// Register the handler
			router.POST("/refresh", authHandler.RefreshToken)

			// Create request with refresh token in header
			req, _ := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			if tt.refreshToken != "" {
				req.Header.Set("X-Refresh-Token", tt.refreshToken)
			}

			// Record the response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check JSON response if expected
			if tt.expectedJSON {
				assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
			}

			// For successful refresh, verify token in response
			if tt.expectedStatus == http.StatusOK {
				var response model.TokenResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tokenResponse.AccessToken, response.AccessToken)
				assert.Equal(t, tokenResponse.RefreshToken, response.RefreshToken)
				assert.Equal(t, tokenResponse.TokenType, response.TokenType)
				assert.Equal(t, tokenResponse.ExpiresIn, response.ExpiresIn)
			}

			// Assert all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}