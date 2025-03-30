#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== E2E Infrastructure API Testing Script ===${NC}"

# Function to make requests and display formatted output
make_request() {
  ENDPOINT=$1
  METHOD=${2:-GET}
  DATA=${3:-""}
  CONTENT_TYPE=${4:-"application/json"}
  AUTH_HEADER=${5:-""}
  
  echo -e "${YELLOW}Testing endpoint: ${BLUE}$METHOD $ENDPOINT${NC}"
  
  if [ -n "$DATA" ]; then
    echo -e "${YELLOW}Request body:${NC}"
    echo -e "${BLUE}$DATA${NC}"
  fi
  
  if [ -n "$AUTH_HEADER" ]; then
    HEADERS="-H \"Content-Type: $CONTENT_TYPE\" -H \"$AUTH_HEADER\""
  else
    HEADERS="-H \"Content-Type: $CONTENT_TYPE\""
  fi
  
  COMMAND="curl -s -X $METHOD $HEADERS"
  
  if [ -n "$DATA" ]; then
    COMMAND="$COMMAND -d '$DATA'"
  fi
  
  COMMAND="$COMMAND $ENDPOINT"
  
  echo -e "${YELLOW}Running command:${NC}"
  echo -e "${BLUE}$COMMAND${NC}"
  
  # Execute the command and capture response
  if [ -n "$DATA" ]; then
    # With data
    if [ -n "$AUTH_HEADER" ]; then
      # With auth
      RESPONSE=$(curl -s -X $METHOD -H "Content-Type: $CONTENT_TYPE" -H "$AUTH_HEADER" -d "$DATA" $ENDPOINT)
    else
      # Without auth
      RESPONSE=$(curl -s -X $METHOD -H "Content-Type: $CONTENT_TYPE" -d "$DATA" $ENDPOINT)
    fi
  else
    # Without data
    if [ -n "$AUTH_HEADER" ]; then
      # With auth
      RESPONSE=$(curl -s -X $METHOD -H "Content-Type: $CONTENT_TYPE" -H "$AUTH_HEADER" $ENDPOINT)
    else
      # Without auth
      RESPONSE=$(curl -s -X $METHOD -H "Content-Type: $CONTENT_TYPE" $ENDPOINT)
    fi
  fi
  
  # Format and print the response if it's JSON
  if [[ $RESPONSE == {* ]] || [[ $RESPONSE == [* ]]; then
    echo -e "${YELLOW}Response:${NC}"
    echo $RESPONSE | jq '.' 2>/dev/null || echo $RESPONSE
  else
    echo -e "${YELLOW}Response:${NC}"
    echo -e "${BLUE}$RESPONSE${NC}"
  fi
  
  echo -e "${GREEN}--------------------------------------------------------${NC}"
  echo ""
  
  # Return the response for further processing
  echo "$RESPONSE"
}

# Test the health endpoints
echo -e "${GREEN}=== Testing Health Endpoints ===${NC}"
make_request "http://localhost:8080/health" "GET"
make_request "http://localhost:8081/health" "GET"
make_request "http://dev.e2e-app.local/health" "GET"
make_request "http://dev.profile.e2e-app.local/health" "GET"

# Register a new user in the auth service
echo -e "${GREEN}=== Creating a test user ===${NC}"
REGISTER_DATA='{
  "email": "test@example.com",
  "password": "password123",
  "first_name": "Test",
  "last_name": "User"
}'
REGISTER_RESPONSE=$(make_request "http://localhost:8080/api/v1/auth/register" "POST" "$REGISTER_DATA")

# Log in to get an auth token
echo -e "${GREEN}=== Logging in to get auth token ===${NC}"
LOGIN_DATA='{
  "email": "test@example.com",
  "password": "password123"
}'
LOGIN_RESPONSE=$(make_request "http://localhost:8080/api/v1/auth/login" "POST" "$LOGIN_DATA")

# Extract the token from the login response
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token' 2>/dev/null)

if [ "$ACCESS_TOKEN" != "null" ] && [ -n "$ACCESS_TOKEN" ]; then
  echo -e "${GREEN}Successfully obtained access token!${NC}"
  
  # Get user profile
  echo -e "${GREEN}=== Testing protected user profile endpoint ===${NC}"
  make_request "http://localhost:8080/api/v1/user/profile" "GET" "" "application/json" "Authorization: Bearer $ACCESS_TOKEN"
  
  # Create a profile in the profile service
  echo -e "${GREEN}=== Creating a user profile ===${NC}"
  PROFILE_DATA='{
    "bio": "This is a test profile",
    "avatar": "https://example.com/avatar.jpg",
    "interests": ["testing", "automation"],
    "social_links": {
      "twitter": "https://twitter.com/testuser",
      "github": "https://github.com/testuser"
    }
  }'
  make_request "http://localhost:8081/api/v1/profiles/" "POST" "$PROFILE_DATA" "application/json" "Authorization: Bearer $ACCESS_TOKEN"
  
  # Get the profile from the profile service
  echo -e "${GREEN}=== Retrieving user profile from profile service ===${NC}"
  USER_ID=$(echo $LOGIN_RESPONSE | jq -r '.user_id' 2>/dev/null || echo "")
  if [ -n "$USER_ID" ]; then
    make_request "http://localhost:8081/api/v1/profiles/$USER_ID" "GET"
  else
    echo -e "${RED}Could not extract user ID from login response${NC}"
  fi
  
  # Test the same endpoints using the ingress hostnames
  echo -e "${GREEN}=== Testing the same endpoints via ingress ===${NC}"
  make_request "http://dev.e2e-app.local/api/v1/user/profile" "GET" "" "application/json" "Authorization: Bearer $ACCESS_TOKEN"
  make_request "http://dev.profile.e2e-app.local/api/v1/profiles/$USER_ID" "GET"
else
  echo -e "${RED}Failed to get access token. Check if the auth service is running correctly.${NC}"
fi

echo -e "${GREEN}=== API Testing Complete ===${NC}"