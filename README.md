# E2E Infrastructure Project

A comprehensive example of modern infrastructure setup using Kubernetes, Helm, and Jenkins.

## Project Overview

This project demonstrates a production-ready infrastructure setup with:
- Kubernetes for container orchestration
- Helm for package management
- Jenkins for CI/CD
- Infrastructure as Code (IaC) principles

## Directory Structure

```
.
├── deployment/
│   ├── charts/              # Helm charts
│   │   ├── e2e-app/         # Authentication service chart
│   │   └── e2e-profile/     # Profile service chart
│   ├── k8s/                # Kubernetes manifests
│   │   ├── base/          # Base configurations
│   │   └── overlays/      # Environment-specific overlays
│   ├── jenkins/           # Jenkins pipeline configurations
│   └── scripts/           # Deployment scripts
├── docs/                  # Documentation
│   ├── infrastructure.md  # Infrastructure overview
│   ├── overlays.md       # Kubernetes overlays guide
│   ├── ingress.md        # Ingress configuration guide
│   ├── jenkins.md        # Jenkins pipeline guide
│   ├── prerequisites.md  # System requirements & setup
│   ├── local-k8s-setup.md # Local Kubernetes deployment
│   └── diagrams/         # Architecture diagrams
├── services/             # Microservices
│   ├── e2e-app/          # Authentication & User Service
│   │   ├── cmd/          # Command line applications
│   │   └── internal/     # Internal packages
│   └── e2e-profile/      # Profile Service
│       ├── cmd/          # Command line applications
│       └── internal/     # Internal packages
└── examples/             # Example configurations
```

## Services

The project consists of two main microservices:

1. **E2E App Service (e2e-app)**
   - Authentication and user management
   - JWT token generation and validation
   - gRPC service for internal communication
   - REST API for external clients

2. **E2E Profile Service (e2e-profile)**
   - User profile management
   - Communicates with e2e-app via gRPC
   - Stores extended user profile information
   - REST API for profile operations

## Features

### Kubernetes
- Base configurations with best practices
- Environment-specific overlays (dev, staging, prod)
- Resource management and limits
- Health checks and probes

### Helm
- Templated Kubernetes resources
- Environment-specific values
- Configurable deployments
- Ingress configurations with SSL

### Infrastructure
- Multi-environment support
- Scalable architecture
- Security best practices
- Monitoring readiness

## Prerequisites

- Kubernetes cluster (1.19+)
- Helm (3.0+)
- kubectl configured with cluster access
- nginx-ingress-controller installed
- cert-manager (for SSL in production)

Detailed setup instructions:
- [System and Environment Prerequisites](docs/prerequisites.md)
- [Local Kubernetes Setup Guide](docs/local-k8s-setup.md)

## Helm Deployment

### Helm Charts Structure

This project provides Helm charts for both services:

1. **e2e-app Chart**: Authentication and user service with PostgreSQL database
2. **e2e-profile Chart**: Profile service with PostgreSQL database

Each chart includes:
- Environment-specific values files (dev, prod)
- PostgreSQL database dependency
- Service, deployment, ingress resources
- Configurable environment variables

### Deploying with Helm

#### Development Environment

```bash
# Step 1: Add the Bitnami repository for PostgreSQL dependency
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Step 2: Deploy the authentication service (e2e-app)
cd deployment/charts
helm install e2e-app ./e2e-app -f ./e2e-app/values-dev.yaml

# Step 3: Deploy the profile service (e2e-profile)
helm install e2e-profile ./e2e-profile -f ./e2e-profile/values-dev.yaml
```

#### Production Environment

```bash
# Step 1: Add the Bitnami repository for PostgreSQL dependency
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Step 2: Deploy the authentication service (e2e-app)
cd deployment/charts
helm install e2e-app ./e2e-app -f ./e2e-app/values-prod.yaml

# Step 3: Deploy the profile service (e2e-profile)
helm install e2e-profile ./e2e-profile -f ./e2e-profile/values-prod.yaml
```

### Configuration Options

Both services can be customized through their respective `values.yaml` files:

#### Authentication Service (e2e-app)

| Parameter                | Description                                | Default               |
|--------------------------|--------------------------------------------|----------------------|
| `replicaCount`          | Number of replicas                         | `1` (dev), `3` (prod) |
| `image.repository`      | Docker image repository                    | `e2e-app`             |
| `image.tag`             | Docker image tag                           | `latest` (dev), `stable` (prod) |
| `service.port`          | HTTP service port                          | `8080`                |
| `service.grpcPort`      | gRPC service port                          | `50051`               |
| `ingress.enabled`       | Enable ingress                             | `true`                |
| `database.enabled`      | Enable PostgreSQL database                 | `true`                |
| `jwt.secret`            | JWT signing secret                         | `your-secret-key`     |
| `jwt.expiration`        | JWT token expiration                       | `24h`                 |

#### Profile Service (e2e-profile)

| Parameter                | Description                                | Default               |
|--------------------------|--------------------------------------------|----------------------|
| `replicaCount`          | Number of replicas                         | `1` (dev), `3` (prod) |
| `image.repository`      | Docker image repository                    | `e2e-profile`         |
| `image.tag`             | Docker image tag                           | `latest` (dev), `stable` (prod) |
| `service.port`          | HTTP service port                          | `8081`                |
| `ingress.enabled`       | Enable ingress                             | `true`                |
| `database.enabled`      | Enable PostgreSQL database                 | `true`                |
| `env.AUTH_SERVICE_URL`  | URL to auth service                        | `http://e2e-app:8080` |
| `env.AUTH_GRPC_ADDR`    | gRPC address for auth service             | `e2e-app:50051`       |

## Documentation

- [System and Environment Prerequisites](docs/prerequisites.md)
- [Local Kubernetes Setup Guide](docs/local-k8s-setup.md)
- [Infrastructure Overview](docs/infrastructure.md)
- [Kubernetes Overlays Guide](docs/overlays.md)
- [Ingress Configuration](docs/ingress.md)

## Environment Configuration

### Development
- Minimal resource requirements
- Local domain setup
- Debug logging enabled
- No SSL requirement

### Production
- High availability setup
- SSL enabled
- Resource limits enforced
- Production-grade monitoring

## Security Features

- RBAC configuration
- Network policies
- SSL/TLS encryption
- Secret management
- Security context settings

## Monitoring & Health Checks

- Liveness probes
- Readiness probes
- Resource monitoring
- Health check endpoints

## Best Practices Implemented

1. **Infrastructure as Code**
   - Version-controlled configurations
   - Environment-specific settings
   - Documented changes

2. **Security**
   - No hardcoded secrets
   - Proper RBAC setup
   - SSL/TLS configuration

3. **Scalability**
   - Resource management
   - Horizontal scaling
   - Load balancing

4. **Maintainability**
   - Clear documentation
   - Modular structure
   - Consistent naming

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Kubernetes best practices
- Helm chart structure guidelines
- Infrastructure security guidelines

## Contact

- Author: tanerincode
- GitHub: [Your GitHub Profile]
- Website: [Your Website]