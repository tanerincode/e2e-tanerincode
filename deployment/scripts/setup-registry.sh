#!/bin/bash

# Get the host and port of the registry
REGISTRY_PORT=$(kubectl -n kube-system get svc registry -o jsonpath='{.spec.ports[0].nodePort}')
REGISTRY_HOST=$(minikube ip)

# Create a port-forwarding from localhost to minikube registry
echo "Setting up port-forwarding to minikube registry..."
echo "Registry is available at $REGISTRY_HOST:$REGISTRY_PORT"
echo "Forwarding localhost:57049 -> $REGISTRY_HOST:$REGISTRY_PORT"

# Run this in background
kubectl port-forward -n kube-system service/registry 57049:80 &

# Save the PID to kill it later
echo $! > /tmp/registry-port-forward.pid

echo "Port forwarding is set up. To kill it later, run:"
echo "kill \$(cat /tmp/registry-port-forward.pid)"