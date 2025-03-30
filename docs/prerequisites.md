# Project Prerequisites

This document outlines the required tools, commands, and configurations needed before running the e2e-tanerincode project on your machine.

## System Requirements Check

Before setting up the environment, run these commands to verify your system capabilities:

```bash
# Check available RAM
sysctl -n hw.memsize | awk '{print $0/1024/1024/1024 " GB"}'

# Check CPU cores
sysctl -n hw.ncpu

# Check available disk space
df -h / | awk 'NR==2 {print $2, $4}'
```

## Required Software

### 1. Docker & Kubernetes

```bash
# Check if Docker is installed
docker version

# Check if Docker has Kubernetes enabled
docker info | grep -E 'Kubernetes|Context'

# Alternative: Check for minikube
minikube version
```

### 2. Kubernetes CLI Tools

```bash
# Verify kubectl installation
kubectl version --client

# Verify kustomize installation
kustomize version
```

## Project Setup Commands

Run these commands in sequence before deploying the application:

### 1. Set up Go environment

The project requires Go 1.23+ (latest version recommended), which may not be available in standard Docker images. For development:

```bash
# Check your Go version
go version

# Install/upgrade Go if needed
# For macOS:
brew install go
# or update:
brew upgrade go
```

### 2. Update Docker Images (Optional)

If you need to modify the Docker build process to use a compatible Go version:

```bash
# To update the Dockerfile to use a newer Go version
sed -i '' 's/FROM golang:1.21-alpine/FROM golang:1.23-alpine/g' e2e-app/Dockerfile
sed -i '' 's/FROM golang:1.21-alpine/FROM golang:1.23-alpine/g' e2e-profile/Dockerfile
```

### 3. Start Minikube with Adequate Resources

```bash
# Start minikube with appropriate resources
minikube start --cpus 4 --memory 6144 --disk-size 20g

# Enable required addons
minikube addons enable ingress
minikube addons enable metrics-server
```

### 4. Configure Docker Environment

```bash
# Set Docker to use Minikube's Docker daemon
eval $(minikube docker-env)
```

### 5. Prepare Kubernetes Namespaces

```bash
# Create required namespaces
kubectl apply -f deployment/k8s/base/namespace.yaml
kubectl create namespace e2e-tanerincode-dev
```

## Verification Steps

After running the prerequisite commands, verify your setup is ready:

```bash
# Verify minikube is running
minikube status

# Check Kubernetes cluster is accessible
kubectl get nodes

# Ensure Docker is configured to use minikube's daemon
docker info | grep -E 'Context'
```

## Common Issues and Solutions

1. **Go Version Compatibility**
   - Issue: Docker images may use older Go versions than required by the project
   - Solution: Modify Dockerfiles to use compatible Go versions or update go.mod

2. **Resource Limitations**
   - Issue: Minikube fails to start due to insufficient resources
   - Solution: Reduce the allocated resources in the minikube start command

3. **Docker Context**
   - Issue: Images built locally not found in Kubernetes
   - Solution: Ensure you're using minikube's Docker daemon with `eval $(minikube docker-env)`

4. **Network Issues**
   - Issue: Services not accessible from host
   - Solution: Use `minikube tunnel` to expose services with LoadBalancer type

## Next Steps

After completing these prerequisites, proceed to the [Local Kubernetes Setup Guide](local-k8s-setup.md) for detailed deployment instructions.