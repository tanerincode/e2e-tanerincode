#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Initializing Database Tables ===${NC}"

# Initialize e2e-app database tables
echo -e "${YELLOW}Creating tables for e2e-app...${NC}"
kubectl exec -it e2e-app-postgresql-0 -n default -- psql -U postgres -d e2e_app -c "
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
"

echo -e "${GREEN}Successfully created tables for e2e-app.${NC}"

# Initialize e2e-profile database tables
echo -e "${YELLOW}Creating tables for e2e-profile...${NC}"
kubectl exec -it e2e-profile-postgresql-0 -n default -- psql -U postgres -d e2e_profile -c "
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";

CREATE TABLE IF NOT EXISTS profile_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    bio TEXT,
    avatar VARCHAR(255),
    interests TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_profile_data_user_id ON profile_data(user_id);

CREATE TABLE IF NOT EXISTS social_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profile_data(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_social_links_profile_id ON social_links(profile_id);
"

echo -e "${GREEN}Successfully created tables for e2e-profile.${NC}"

# Restart deployments to ensure they connect with the new tables
echo -e "${YELLOW}Restarting deployments...${NC}"
kubectl rollout restart deployment e2e-app -n default
kubectl rollout restart deployment e2e-profile -n default

echo -e "${GREEN}Tables initialized successfully!${NC}"