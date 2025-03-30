#!/bin/bash

# Exit on any error
set -e

# Change to the project root directory
cd "$(dirname "$0")/../.."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building Docker images...${NC}"

# Function to build and push Docker images
build_and_push() {
    local service=$1
    echo -e "${GREEN}Building $service image...${NC}"
    
    # Build the Docker image using the Dockerfile in the service directory
    docker build -t tanerincode/$service:latest ./$service
    
    # Tag the image with minikube registry
    echo -e "${GREEN}Tagging $service image for minikube...${NC}"
    docker tag tanerincode/$service:latest localhost:57049/tanerincode/$service:latest
    
    # Push to minikube registry
    echo -e "${GREEN}Pushing $service image to minikube registry...${NC}"
    docker push localhost:57049/tanerincode/$service:latest
}

# Build and push both services
build_and_push "e2e-app"
build_and_push "e2e-profile"

echo -e "${YELLOW}Deploying to Kubernetes...${NC}"

# Deploy PostgreSQL using Helm (if needed)
echo -e "${GREEN}Deploying PostgreSQL...${NC}"
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm upgrade --install postgres bitnami/postgresql \
  --set auth.username=app_user \
  --set auth.password=app_password \
  --set auth.database=app_db \
  --set persistence.size=1Gi \
  -n default

# Wait for PostgreSQL to be ready
echo -e "${GREEN}Waiting for PostgreSQL to be ready...${NC}"
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgresql --timeout=180s

# Deploy the e2e-app service
echo -e "${GREEN}Deploying e2e-app service...${NC}"
helm upgrade --install e2e-app ./deployment/charts/e2e-app \
  --set image.repository=localhost:57049/tanerincode/e2e-app \
  --set image.tag=latest \
  --set database.enabled=true \
  --set database.postgresql.auth.username=app_user \
  --set database.postgresql.auth.password=app_password \
  --set database.postgresql.auth.database=app_db \
  -n default

# Deploy the e2e-profile service
echo -e "${GREEN}Deploying e2e-profile service...${NC}"
helm upgrade --install e2e-profile ./deployment/charts/e2e-profile \
  --set image.repository=localhost:57049/tanerincode/e2e-profile \
  --set image.tag=latest \
  --set authService.url=http://e2e-app:8080 \
  --set authService.grpcAddr=e2e-app:9090 \
  --set database.enabled=true \
  --set database.postgresql.auth.username=app_user \
  --set database.postgresql.auth.password=app_password \
  --set database.postgresql.auth.database=app_db \
  -n default

# Wait for the services to be ready
echo -e "${GREEN}Waiting for services to be ready...${NC}"
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=e2e-app --timeout=180s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=e2e-profile --timeout=180s

# Get service URLs
echo -e "${YELLOW}Deployment complete!${NC}"
echo -e "${GREEN}Services:${NC}"
echo -e "e2e-app: $(minikube service e2e-app --url)"
echo -e "e2e-profile: $(minikube service e2e-profile --url)"

echo -e "${YELLOW}Run the following command to port-forward to the services:${NC}"
echo -e "kubectl port-forward svc/e2e-app 8080:8080"
echo -e "kubectl port-forward svc/e2e-profile 8081:8080"