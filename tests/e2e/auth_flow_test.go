package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Constants for testing
const (
	AuthServiceURL    = "http://localhost:8080/api/v1"
	ProfileServiceURL = "http://localhost:8081/api/v1"
)

// Test data structures
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProfileRequest struct {
	Bio      string `json:"bio"`
	Location string `json:"location"`
	Website  string `json:"website"`
}

type ProfileResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Bio       string    `json:"bio"`
	Location  string    `json:"location"`
	Website   string    `json:"website"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TestAuthenticationFlow is an end-to-end test for the complete authentication flow
func TestAuthenticationFlow(t *testing.T) {
	// Skip if not running e2e tests
	if os.Getenv("RUN_E2E_TESTS") != "true" {
		t.Skip("Skipping end-to-end tests. Set RUN_E2E_TESTS=true to run.")
	}

	// Generate unique email for testing
	uniqueID := uuid.New().String()[:8]
	testEmail := fmt.Sprintf("test_%s@example.com", uniqueID)

	// Step 1: Register a new user
	t.Log("STEP 1: Registering a new user")
	registerReq := RegisterRequest{
		Email:     testEmail,
		Password:  "Password123!",
		FirstName: "Test",
		LastName:  "User",
	}

	var user UserResponse
	statusCode, err := sendRequest(http.MethodPost, AuthServiceURL+"/auth/register", registerReq, &user)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, statusCode)

	// Validate registration response
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)

	// Step 2: Login with the registered user
	t.Log("STEP 2: Logging in with the registered user")
	loginReq := LoginRequest{
		Email:    testEmail,
		Password: "Password123!",
	}

	var tokenResp TokenResponse
	statusCode, err = sendRequest(http.MethodPost, AuthServiceURL+"/auth/login", loginReq, &tokenResp)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	// Validate login response
	assert.NotEmpty(t, tokenResp.AccessToken)
	assert.NotEmpty(t, tokenResp.RefreshToken)
	assert.Equal(t, "Bearer", tokenResp.TokenType)
	assert.Greater(t, tokenResp.ExpiresIn, int64(0))

	// Step 3: Create a user profile using the auth token
	t.Log("STEP 3: Creating a user profile with authentication")
	createProfileReq := CreateProfileRequest{
		Bio:      "Test bio for end-to-end testing",
		Location: "Test Location",
		Website:  "https://test.example.com",
	}

	var profile ProfileResponse
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", tokenResp.AccessToken),
	}

	statusCode, err = sendRequestWithHeaders(http.MethodPost, ProfileServiceURL+"/profiles", createProfileReq, &profile, headers)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, statusCode)

	// Validate profile creation
	assert.NotEmpty(t, profile.ID)
	assert.Equal(t, user.ID, profile.UserID)
	assert.Equal(t, "Test bio for end-to-end testing", profile.Bio)
	assert.Equal(t, "Test Location", profile.Location)
	assert.Equal(t, "https://test.example.com", profile.Website)

	// Step 4: Get profile with auth token
	t.Log("STEP 4: Getting user profile with authentication")
	var retrievedProfile ProfileResponse
	statusCode, err = sendRequestWithHeaders(http.MethodGet, ProfileServiceURL+"/profiles/"+profile.ID, nil, &retrievedProfile, headers)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	// Validate retrieved profile
	assert.Equal(t, profile.ID, retrievedProfile.ID)
	assert.Equal(t, user.ID, retrievedProfile.UserID)
	assert.Equal(t, profile.Bio, retrievedProfile.Bio)

	// Step 5: Try accessing profile without token (should fail)
	t.Log("STEP 5: Attempting to access profile without authentication")
	var unauthorizedResp map[string]interface{}
	statusCode, err = sendRequest(http.MethodGet, ProfileServiceURL+"/profiles/"+profile.ID, nil, &unauthorizedResp)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, statusCode)

	// Step 6: Refresh the token
	t.Log("STEP 6: Refreshing the auth token")
	refreshHeaders := map[string]string{
		"X-Refresh-Token": tokenResp.RefreshToken,
	}

	var newTokenResp TokenResponse
	statusCode, err = sendRequestWithHeaders(http.MethodGet, AuthServiceURL+"/auth/refresh", nil, &newTokenResp, refreshHeaders)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	// Validate refreshed token
	assert.NotEmpty(t, newTokenResp.AccessToken)
	assert.NotEmpty(t, newTokenResp.RefreshToken)
	assert.NotEqual(t, tokenResp.AccessToken, newTokenResp.AccessToken)

	// Step 7: Use the new token to access the profile
	t.Log("STEP 7: Using refreshed token to access profile")
	newHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", newTokenResp.AccessToken),
	}

	var profileAfterRefresh ProfileResponse
	statusCode, err = sendRequestWithHeaders(http.MethodGet, ProfileServiceURL+"/profiles/"+profile.ID, nil, &profileAfterRefresh, newHeaders)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	// Validate profile access with new token
	assert.Equal(t, profile.ID, profileAfterRefresh.ID)
}

// Helper function to send HTTP requests
func sendRequest(method, url string, requestBody interface{}, responseBody interface{}) (int, error) {
	return sendRequestWithHeaders(method, url, requestBody, responseBody, nil)
}

// Helper function to send HTTP requests with custom headers
func sendRequestWithHeaders(method, url string, requestBody interface{}, responseBody interface{}, headers map[string]string) (int, error) {
	var reqBody io.Reader

	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return 0, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, err
	}

	// Set Content-Type header for requests with body
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// Handle empty response for non-success status codes
	if resp.StatusCode >= 400 && len(respBody) == 0 {
		return resp.StatusCode, nil
	}

	// Parse response body if provided
	if responseBody != nil && len(respBody) > 0 {
		err = json.Unmarshal(respBody, responseBody)
		if err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}
