# Local Kubernetes Setup Guide

This document provides instructions for setting up the e2e-tanerincode project on a local Kubernetes environment.

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop) with Kubernetes enabled, or [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) configured to use your local cluster
- [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/) (version 4.0.0+)
- Git (to clone the repository)

## Setup Options

Choose one of the following options for your local Kubernetes environment:

### Option 1: Docker Desktop Kubernetes

1. Open Docker Desktop
2. Go to Settings > Kubernetes
3. Check "Enable Kubernetes"
4. Click "Apply & Restart"
5. Wait for the Kubernetes cluster to start

### Option 2: Minikube

```bash
# Install Minikube
# macOS with Homebrew
brew install minikube

# Start Minikube with adequate resources
minikube start --cpus 4 --memory 8192 --disk-size 20g

# Verify the cluster is running
minikube status
```

## Building and Loading the Images

Before deploying to your local Kubernetes cluster, you need to build and make the Docker images available:

### For Docker Desktop Kubernetes

```bash
# Navigate to the e2e-app directory
cd /Users/tombastaner/private-projects/e2e-app

# Build the Docker image
docker build -t e2e-app:latest .

# Navigate to the e2e-profile directory
cd /Users/tombastaner/private-projects/e2e-profile

# Build the Docker image
docker build -t e2e-profile:latest .
```

### For Minikube

```bash
# Navigate to the e2e-app directory
cd /Users/tombastaner/private-projects/e2e-app

# Set Docker to use Minikube's Docker daemon
eval $(minikube docker-env)

# Build the Docker image
docker build -t e2e-app:latest .

# Navigate to the e2e-profile directory
cd /Users/tombastaner/private-projects/e2e-profile

# Build the Docker image
docker build -t e2e-profile:latest .
```

## Deploying the Application

### Step 1: Create Required Namespaces

```bash
# Apply the base namespace configuration
kubectl apply -f /Users/tombastaner/private-projects/deployment/k8s/base/namespace.yaml

# Create the dev namespace manually since it's defined in kustomization
kubectl create namespace e2e-tanerincode-dev
```

### Step 2: Deploy Using Kustomize (Development Environment)

```bash
# Navigate to the deployment directory
cd /Users/tombastaner/private-projects/deployment/k8s

# Apply the development overlay
kubectl apply -k overlays/dev
```

## Verifying the Deployment

```bash
# Check all resources in the dev namespace
kubectl get all -n e2e-tanerincode-dev

# Check deployment status
kubectl get deployments -n e2e-tanerincode-dev

# Check pods status
kubectl get pods -n e2e-tanerincode-dev

# Check services
kubectl get services -n e2e-tanerincode-dev

# View ConfigMap details
kubectl get configmap e2e-app-config -n e2e-tanerincode-dev -o yaml

# View logs from a pod (replace pod-name with actual pod name)
kubectl logs -n e2e-tanerincode-dev pod-name
```

## Accessing the Application

### For Docker Desktop Kubernetes

The service is exposed as ClusterIP by default. To access it:

```bash
# Port forward the service to your local machine
kubectl port-forward -n e2e-tanerincode-dev service/e2e-app 8080:80
```

Now you can access the application at `http://localhost:8080`

### For Minikube

```bash
# Create a URL to access the service
minikube service e2e-app -n e2e-tanerincode-dev --url
```

This will print a URL that you can use to access the service.

## Troubleshooting

### Common Issues and Solutions

1. **Pods not starting**
   ```bash
   # Check the detailed status of the pod
   kubectl describe pod -n e2e-tanerincode-dev [pod-name]
   
   # Check the logs
   kubectl logs -n e2e-tanerincode-dev [pod-name]
   ```

2. **Image pull issues**
   - For Minikube, ensure you've built the images using Minikube's Docker daemon (`eval $(minikube docker-env)`)
   - For Docker Desktop, ensure the images are present in your local Docker (`docker images`)

3. **ConfigMap or Secret issues**
   ```bash
   # Check if ConfigMap exists
   kubectl get configmap -n e2e-tanerincode-dev
   
   # Check if Secret exists
   kubectl get secret -n e2e-tanerincode-dev
   ```

## Cleaning Up

When you're done with the deployment, you can clean up the resources:

```bash
# Delete all resources created with kustomize
kubectl delete -k /Users/tombastaner/private-projects/deployment/k8s/overlays/dev

# Delete the namespaces
kubectl delete namespace e2e-tanerincode-dev
kubectl delete namespace e2e-tanerincode
```

## Next Steps

- Set up a local PostgreSQL instance if required by your services
- Configure port forwarding for multiple services
- Add Ingress controller for more sophisticated routing