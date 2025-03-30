#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Viewing Ingress Controller Logs ===${NC}"

# Check if namespace argument is provided
NAMESPACE=${1:-"ingress-nginx"}

# Find ingress controller pod
echo -e "${YELLOW}Looking for ingress controller pod in namespace: $NAMESPACE${NC}"

INGRESS_POD=$(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=ingress-nginx,app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

# If not found in the provided namespace, try kube-system
if [ -z "$INGRESS_POD" ]; then
  echo -e "${YELLOW}Ingress controller not found in $NAMESPACE, trying kube-system namespace...${NC}"
  NAMESPACE="kube-system"
  INGRESS_POD=$(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=ingress-nginx,app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
  
  # If still not found, try searching for the nginx pod by common labels
  if [ -z "$INGRESS_POD" ]; then
    echo -e "${YELLOW}Trying to find ingress by common labels...${NC}"
    INGRESS_POD=$(kubectl get pods -n $NAMESPACE -l app=nginx-ingress -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
  fi
  
  # For minikube specifically
  if [ -z "$INGRESS_POD" ]; then
    echo -e "${YELLOW}Trying to find minikube's ingress controller...${NC}"
    NAMESPACE="ingress-nginx"
    INGRESS_POD=$(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
  fi
fi

# Final attempt with a broader search
if [ -z "$INGRESS_POD" ]; then
  echo -e "${YELLOW}Searching across all namespaces for ingress controller...${NC}"
  INGRESS_NS_POD=$(kubectl get pods -A | grep -i ingress | grep -i controller | head -1)
  NAMESPACE=$(echo "$INGRESS_NS_POD" | awk '{print $1}')
  INGRESS_POD=$(echo "$INGRESS_NS_POD" | awk '{print $2}')
fi

if [ -z "$INGRESS_POD" ]; then
  echo -e "${RED}Error: Could not find the ingress controller pod.${NC}"
  echo -e "${YELLOW}Available pods across all namespaces:${NC}"
  kubectl get pods -A
  exit 1
fi

echo -e "${GREEN}Found ingress controller pod: $INGRESS_POD in namespace: $NAMESPACE${NC}"

# View pod details
echo -e "${YELLOW}Ingress controller pod details:${NC}"
kubectl describe pod $INGRESS_POD -n $NAMESPACE

# View logs
echo -e "${YELLOW}Fetching logs for ingress controller...${NC}"
echo -e "${GREEN}Press Ctrl+C to exit log view${NC}"
kubectl logs -f $INGRESS_POD -n $NAMESPACE