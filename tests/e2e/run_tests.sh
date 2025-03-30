#!/bin/bash

# Script to run end-to-end tests for the authentication flow

# Set environment variables
export RUN_E2E_TESTS=true

# Check if services are running
echo "Checking if authentication service is running..."
curl -s http://localhost:8080/health > /dev/null
if [ $? -ne 0 ]; then
    echo "Error: Authentication service is not running on port 8080."
    echo "Please start the service before running end-to-end tests."
    exit 1
fi

echo "Checking if profile service is running..."
curl -s http://localhost:8081/health > /dev/null
if [ $? -ne 0 ]; then
    echo "Error: Profile service is not running on port 8081."
    echo "Please start the service before running end-to-end tests."
    exit 1
fi

# Run the tests
echo "Running end-to-end authentication flow tests..."
cd "$(dirname "$0")"
go test -v