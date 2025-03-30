# Understanding Kubernetes Overlays with Kustomize

## What are Overlays?
Think of overlays like layers of transparent sheets on a projector:
- Base layer: Your common, default configuration
- Overlay layers: Environment-specific changes that go on top

## Directory Structure
```
k8s/
├── base/                 # Your foundation
│   ├── deployment.yaml   # Common deployment settings
│   ├── service.yaml     # Common service settings
│   └── kustomization.yaml
└── overlays/            # Your environment-specific changes
    ├── dev/
    │   └── kustomization.yaml
    ├── staging/
    │   └── kustomization.yaml
    └── prod/
        └── kustomization.yaml
```

## How It Works

### 1. Base Configuration (k8s/base/deployment.yaml)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: e2e-app
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: e2e-app
        image: e2e-app:latest
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
```

### 2. Development Overlay (k8s/overlays/dev/kustomization.yaml)
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- ../../base
namespace: e2e-tanerincode-dev
commonLabels:
  environment: development
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: e2e-app
    spec:
      template:
        spec:
          containers:
          - name: e2e-app
            env:
            - name: DEBUG
              value: "true"
```

### 3. Production Overlay (k8s/overlays/prod/kustomization.yaml)
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- ../../base
namespace: e2e-tanerincode-prod
commonLabels:
  environment: production
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: e2e-app
    spec:
      replicas: 3
      template:
        spec:
          containers:
          - name: e2e-app
            resources:
              requests:
                cpu: "500m"
                memory: "512Mi"
```

## Real-World Examples

### 1. Changing Resource Limits
```yaml
# Dev environment (low resources)
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: e2e-app
    spec:
      template:
        spec:
          containers:
          - name: e2e-app
            resources:
              limits:
                cpu: "200m"
                memory: "256Mi"

# Prod environment (high resources)
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: e2e-app
    spec:
      template:
        spec:
          containers:
          - name: e2e-app
            resources:
              limits:
                cpu: "1000m"
                memory: "1Gi"
```

### 2. Different Environment Variables
```yaml
# Dev environment
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: e2e-app
    spec:
      template:
        spec:
          containers:
          - name: e2e-app
            env:
            - name: LOG_LEVEL
              value: "debug"
            - name: API_URL
              value: "http://dev-api.example.com"

# Prod environment
patches:
- patch: |-
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: e2e-app
    spec:
      template:
        spec:
          containers:
          - name: e2e-app
            env:
            - name: LOG_LEVEL
              value: "info"
            - name: API_URL
              value: "http://api.example.com"
```

## Common Use Cases

1. **Scaling**
   ```yaml
   # Dev: 1 replica
   replicas: 1

   # Staging: 2 replicas
   replicas: 2

   # Prod: 3 replicas with auto-scaling
   replicas: 3
   ```

2. **Different Images**
   ```yaml
   # Dev: Latest image
   image: e2e-app:latest

   # Prod: Specific version
   image: e2e-app:v1.2.3
   ```

3. **Resource Management**
   ```yaml
   # Dev: Minimal resources
   resources:
     requests:
       cpu: "100m"
       memory: "128Mi"

   # Prod: Production resources
   resources:
     requests:
       cpu: "500m"
       memory: "512Mi"
   ```

## How to Use

1. **View the final configuration:**
   ```bash
   # Development
   kubectl kustomize deployment/k8s/overlays/dev

   # Production
   kubectl kustomize deployment/k8s/overlays/prod
   ```

2. **Apply the configuration:**
   ```bash
   # Development
   kubectl apply -k deployment/k8s/overlays/dev

   # Production
   kubectl apply -k deployment/k8s/overlays/prod
   ```

## Best Practices

1. **Keep Base Minimal**
   - Base should contain only common configurations
   - Avoid environment-specific settings in base

2. **Use Clear Naming**
   - Name overlays clearly (dev, staging, prod)
   - Use consistent naming across resources

3. **Document Changes**
   - Comment significant changes in overlays
   - Explain why certain values are different

4. **Version Control**
   - Keep all overlays in version control
   - Review changes to overlays carefully

5. **Testing**
   - Test overlay configurations before applying
   - Use `kubectl kustomize` to preview changes