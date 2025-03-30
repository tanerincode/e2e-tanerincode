#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgresql --timeout=180s || true

# Initialize e2e-app database
echo "Creating database for e2e-app..."
kubectl exec -it e2e-app-postgresql-0 -- psql -U postgres -c "CREATE DATABASE e2e_app;" || true

# Initialize e2e-profile database
echo "Creating database for e2e-profile..."
kubectl exec -it e2e-profile-postgresql-0 -- psql -U postgres -c "CREATE DATABASE e2e_profile;" || true

echo "Database initialization completed successfully!"