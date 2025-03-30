#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== E2E Infrastructure Hosts File Updater ===${NC}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Please run as root (sudo $0)${NC}"
  exit 1
fi

# Get minikube IP
if ! command -v minikube &> /dev/null; then
  echo -e "${RED}Error: minikube is not installed.${NC}"
  exit 1
fi

MINIKUBE_IP=$(minikube ip)
if [ -z "$MINIKUBE_IP" ]; then
  echo -e "${RED}Error: Could not get minikube IP. Is minikube running?${NC}"
  exit 1
fi

echo -e "${YELLOW}Minikube IP: ${MINIKUBE_IP}${NC}"

# Hosts to add
HOSTS=(
  "api.e2e-app.local"
  "profile.e2e-app.local"
  "dev.e2e-app.local"
  "dev.profile.e2e-app.local"
  "e2e-app.local"
  "e2e-profile.local"
)

# Backup hosts file
cp /etc/hosts /etc/hosts.bak
echo -e "${GREEN}Backed up hosts file to /etc/hosts.bak${NC}"

# Remove any previous entries for these hosts
for HOST in "${HOSTS[@]}"; do
  sed -i '' "/[[:space:]]$HOST$/d" /etc/hosts
done

# Add new entries
echo -e "${YELLOW}Adding hosts entries...${NC}"
for HOST in "${HOSTS[@]}"; do
  echo "$MINIKUBE_IP $HOST" >> /etc/hosts
  echo -e "${GREEN}Added: $MINIKUBE_IP $HOST${NC}"
done

echo -e "${GREEN}Hosts file updated successfully!${NC}"
echo -e "${YELLOW}You can now access your services at:${NC}"
echo -e "${GREEN}http://api.e2e-app.local${NC}"
echo -e "${GREEN}http://profile.e2e-app.local${NC}"
echo -e "${GREEN}http://dev.e2e-app.local${NC}"
echo -e "${GREEN}http://dev.profile.e2e-app.local${NC}"
echo -e "${GREEN}http://e2e-app.local${NC}"
echo -e "${GREEN}http://e2e-profile.local${NC}"