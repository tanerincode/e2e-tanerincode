.PHONY: test test-unit test-integration test-e2e

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	@cd e2e-app && go test -v ./internal/service/...
	@cd e2e-app && go test -v ./internal/grpc/server/...
	@cd e2e-profile && go test -v ./internal/grpc/client/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@cd e2e-app && go test -v ./internal/handler/...

# Run end-to-end tests (services must be running)
test-e2e:
	@echo "Running end-to-end tests..."
	@cd tests/e2e && ./run_tests.sh