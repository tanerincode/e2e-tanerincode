#!/bin/bash

# Change to the project root directory
cd "$(dirname "$0")/../.."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Cleaning up Kubernetes resources...${NC}"

# Uninstall Helm releases
echo -e "${GREEN}Uninstalling e2e-profile service...${NC}"
helm uninstall e2e-profile

echo -e "${GREEN}Uninstalling e2e-app service...${NC}"
helm uninstall e2e-app

echo -e "${GREEN}Uninstalling PostgreSQL...${NC}"
helm uninstall postgres

# Delete any remaining resources
echo -e "${GREEN}Checking for any remaining resources...${NC}"
kubectl delete pod,svc,deploy,sts,pvc,pv -l app.kubernetes.io/part-of=e2e-tanerincode --ignore-not-found

echo -e "${YELLOW}Cleanup complete!${NC}"