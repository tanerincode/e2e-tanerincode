# Infrastructure Documentation

## Visual Overview

For detailed visual representations of our infrastructure, please refer to the following diagrams in `docs/diagrams/architecture.md`:

1. **Project Structure Diagram**: Shows the complete directory structure and relationships
2. **Kubernetes vs Helm Workflow**: Illustrates the difference between direct Kubernetes deployment and Helm-managed deployment
3. **Deployment Process**: Sequence diagram showing the complete deployment flow
4. **Environment Management**: Shows how different environments are managed
5. **Release Management**: Visualizes the Helm release and rollback process

## Project Structure
```
.
├── deployment/
│   ├── charts/                 # Helm charts for deployment
│   │   └── e2e-app/           # Main application chart
│   │       ├── Chart.yaml     # Chart metadata
│   │       ├── values.yaml    # Default values
│   │       └── templates/     # Templated K8s manifests
│   ├── k8s/                   # Raw Kubernetes manifests
│   │   ├── base/             # Base configurations
│   │   └── overlays/         # Environment-specific overlays
│   ├── jenkins/              # Jenkins pipeline configurations
│   └── scripts/              # Deployment scripts
├── docs/                     # Documentation
└── examples/                 # Example configurations
```

## Kubernetes vs Helm

### Kubernetes (k8s/base)
Raw Kubernetes manifests used for:
- Direct kubectl applications
- Local development
- Simple deployments
- Testing configurations

Example:
```yaml
# k8s/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: e2e-app
  namespace: e2e-tanerincode
spec:
  replicas: 1
  # ... static configuration
```

### Helm Charts (charts/e2e-app)
Templated versions for:
- Environment management
- Configuration management
- Release management
- Easy rollbacks

Example:
```yaml
# charts/e2e-app/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "e2e-app.fullname" . }}
spec:
  replicas: {{ .Values.replicaCount }}
  # ... dynamic configuration
```

## Configuration Management

### Kubernetes Configuration
- Static values in YAML files
- Environment changes require file modifications
- Manual version control
- Direct kubectl commands:
  ```bash
  kubectl apply -f k8s/base/
  ```

### Helm Configuration
- Values separated from templates
- Environment-specific value files
- Built-in versioning
- Release management:
  ```bash
  # Install/Upgrade with default values
  helm upgrade --install e2e-app ./charts/e2e-app

  # Install/Upgrade with environment values
  helm upgrade --install e2e-app ./charts/e2e-app -f values-prod.yaml
  ```

## Environment Management

### Using Kubernetes
1. Create separate overlay directories:
   ```
   k8s/overlays/
   ├── dev/
   ├── staging/
   └── prod/
   ```

2. Apply environment-specific configurations:
   ```bash
   kubectl apply -f k8s/overlays/prod/
   ```

### Using Helm
1. Create environment-specific values:
   ```
   charts/e2e-app/
   ├── values.yaml          # Default values
   ├── values-dev.yaml     # Development values
   ├── values-staging.yaml # Staging values
   └── values-prod.yaml    # Production values
   ```

2. Deploy to specific environment:
   ```bash
   helm upgrade --install e2e-app ./charts/e2e-app -f values-prod.yaml
   ```

## Best Practices

### Kubernetes Best Practices
1. Keep base configurations simple
2. Use namespaces for isolation
3. Define resource limits
4. Implement health checks
5. Use labels consistently

### Helm Best Practices
1. Version your charts
2. Document all values
3. Use named templates
4. Keep release history
5. Implement rollback plans

## Common Operations

### Kubernetes Operations
```bash
# Apply configurations
kubectl apply -f k8s/base/

# Get resources
kubectl get pods -n e2e-tanerincode

# Describe resources
kubectl describe deployment e2e-app -n e2e-tanerincode
```

### Helm Operations
```bash
# Install/Upgrade release
helm upgrade --install e2e-app ./charts/e2e-app

# List releases
helm list

# Rollback to previous release
helm rollback e2e-app 1

# Uninstall release
helm uninstall e2e-app
```

## When to Use What

### Use Kubernetes (k8s/base) When:
- Developing locally
- Testing configurations
- Need simple, direct deployment
- Learning/understanding the infrastructure

### Use Helm (charts) When:
- Deploying to multiple environments
- Need configuration management
- Want release management
- Require easy rollbacks
- Sharing the application
- Production deployments

## Deployment Process

1. Development:
   ```bash
   # Local testing with Kubernetes
   kubectl apply -f k8s/base/

   # Development deployment with Helm
   helm upgrade --install e2e-app ./charts/e2e-app -f values-dev.yaml
   ```

2. Staging:
   ```bash
   helm upgrade --install e2e-app ./charts/e2e-app -f values-staging.yaml
   ```

3. Production:
   ```bash
   helm upgrade --install e2e-app ./charts/e2e-app -f values-prod.yaml
   ```

## Monitoring and Maintenance

1. Check application status:
   ```bash
   # Kubernetes
   kubectl get all -n e2e-tanerincode

   # Helm
   helm status e2e-app
   ```

2. View logs:
   ```bash
   kubectl logs -l app=e2e-app -n e2e-tanerincode
   ```

3. Debug issues:
   ```bash
   kubectl describe pod -l app=e2e-app -n e2e-tanerincode
   ```

## Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [Helm Best Practices](https://helm.sh/docs/chart_best_practices/)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)