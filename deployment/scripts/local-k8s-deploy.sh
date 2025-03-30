#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== E2E Infrastructure Local Kubernetes Deployment ===${NC}"

# Check if required tools are installed
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed. Please install it first.${NC}"
        exit 1
    fi
}

check_command kubectl
check_command docker
check_command minikube

# Check if minikube is running
if ! minikube status | grep -q "Running"; then
    echo -e "${YELLOW}Starting minikube...${NC}"
    minikube start --driver=docker --memory=4096 --cpus=2
else
    echo -e "${GREEN}Minikube is already running.${NC}"
fi

# Enable ingress addon if not already enabled
if ! minikube addons list | grep -q "ingress.*enabled"; then
    echo -e "${YELLOW}Enabling ingress addon...${NC}"
    minikube addons enable ingress
fi

# Create namespaces if they don't exist
echo -e "${YELLOW}Creating namespaces...${NC}"
kubectl create namespace dev --dry-run=client -o yaml | kubectl apply -f -

# Set the current context to dev namespace
kubectl config set-context --current --namespace=dev

# Build and push Docker images to local minikube
echo -e "${YELLOW}Building Docker images and pushing to minikube...${NC}"

# Point shell to minikube's Docker daemon
eval $(minikube docker-env)

# Build e2e-app image
echo -e "${YELLOW}Building e2e-app image...${NC}"
cd services/e2e-app
docker build -t e2e-app:latest .

# Build e2e-profile image
echo -e "${YELLOW}Building e2e-profile image...${NC}"
cd ../e2e-profile
docker build -t e2e-profile:latest .
cd ../..

# Create PostgreSQL secrets
echo -e "${YELLOW}Creating PostgreSQL secrets...${NC}"
kubectl create secret generic postgres-secret \
    --from-literal=postgres-password=postgres \
    --namespace=dev \
    --dry-run=client -o yaml | kubectl apply -f -

# Deploy using Helm
echo -e "${YELLOW}Deploying applications using Helm...${NC}"

# Add bitnami repo if not already added
if ! helm repo list | grep -q "bitnami"; then
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
fi

# Deploy e2e-app
echo -e "${YELLOW}Deploying e2e-app...${NC}"
helm upgrade --install e2e-app ./deployment/charts/e2e-app \
    --set image.repository=e2e-app \
    --set image.tag=latest \
    --set image.pullPolicy=Never \
    --set database.postgresPassword=postgres \
    --namespace=dev \
    -f ./deployment/charts/e2e-app/values-dev.yaml

# Deploy e2e-profile
echo -e "${YELLOW}Deploying e2e-profile...${NC}"
helm upgrade --install e2e-profile ./deployment/charts/e2e-profile \
    --set image.repository=e2e-profile \
    --set image.tag=latest \
    --set image.pullPolicy=Never \
    --set database.postgresPassword=postgres \
    --namespace=dev \
    -f ./deployment/charts/e2e-profile/values-dev.yaml

# Wait for pods to be ready
echo -e "${YELLOW}Waiting for pods to be ready...${NC}"
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=e2e-app --timeout=180s -n dev
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=e2e-profile --timeout=180s -n dev

# Initialize databases if needed using the script
echo -e "${YELLOW}Initializing databases...${NC}"
./deployment/scripts/init-database.sh

# Display ingress configurations
echo -e "${YELLOW}Displaying ingress configurations...${NC}"
kubectl get ingress -A

# Get minikube IP
MINIKUBE_IP=$(minikube ip)
echo -e "${GREEN}Minikube IP: ${MINIKUBE_IP}${NC}"

# Hosts file update instructions
echo -e "${YELLOW}====== IMPORTANT: UPDATE YOUR HOSTS FILE ======${NC}"
echo -e "${YELLOW}Add the following entries to your /etc/hosts file:${NC}"
echo -e "${GREEN}${MINIKUBE_IP} api.e2e-app.local${NC}"
echo -e "${GREEN}${MINIKUBE_IP} profile.e2e-app.local${NC}"
echo -e "${GREEN}${MINIKUBE_IP} dev.e2e-app.local${NC}"
echo -e "${GREEN}${MINIKUBE_IP} dev.profile.e2e-app.local${NC}"
echo -e "${GREEN}${MINIKUBE_IP} e2e-app.local${NC}"
echo -e "${GREEN}${MINIKUBE_IP} e2e-profile.local${NC}"
echo -e "${YELLOW}Run: sudo nano /etc/hosts${NC}"
echo -e "${YELLOW}====== END HOSTS FILE INSTRUCTIONS ======${NC}"

# Set up port forwarding
echo -e "${GREEN}Setting up port forwarding...${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop port forwarding when you're done.${NC}"
echo -e "${GREEN}e2e-app will be available at http://localhost:8080${NC}"
echo -e "${GREEN}e2e-profile will be available at http://localhost:8081${NC}"

# Start port forwarding in the background
kubectl port-forward svc/e2e-app 8080:8080 -n dev &
PF1_PID=$!
kubectl port-forward svc/e2e-profile 8081:8081 -n dev &
PF2_PID=$!

# Trap to kill port forwarding when the script is terminated
trap "kill $PF1_PID $PF2_PID; echo -e '${YELLOW}Port forwarding stopped.${NC}'" INT TERM EXIT

# Keep the script running until manually terminated
echo -e "${GREEN}Deployment complete!${NC}"
echo -e "${YELLOW}Services are now accessible at:${NC}"
echo -e "${GREEN}e2e-app: http://localhost:8080${NC}"
echo -e "${GREEN}e2e-profile: http://localhost:8081${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop port forwarding.${NC}"

# Wait indefinitely
while true; do
    sleep 1
done