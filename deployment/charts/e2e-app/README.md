# E2E App Helm Chart

This Helm chart deploys the E2E App authentication and user service to a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure

## Installing the Chart

To install the chart with the release name `my-auth`:

```bash
helm install my-auth /path/to/e2e-app
```

## Configuration

The following table lists the configurable parameters of the e2e-app chart and their default values.

| Parameter                         | Description                                                                                                        | Default                           |
|-----------------------------------|--------------------------------------------------------------------------------------------------------------------|-----------------------------------|
| `replicaCount`                    | Number of replicas                                                                                                 | `1`                               |
| `image.repository`                | Image repository                                                                                                   | `e2e-app`                         |
| `image.tag`                       | Image tag                                                                                                          | `latest`                          |
| `image.pullPolicy`                | Image pull policy                                                                                                  | `Always`                          |
| `service.type`                    | Kubernetes Service type                                                                                            | `ClusterIP`                       |
| `service.port`                    | Service HTTP port                                                                                                  | `8080`                            |
| `service.grpcPort`                | Service gRPC port                                                                                                  | `50051`                           |
| `ingress.enabled`                 | Enable ingress controller resource                                                                                 | `true`                            |
| `ingress.className`               | IngressClass that will be be used                                                                                  | `nginx`                           |
| `ingress.hosts[0].host`           | Hostname to your installation                                                                                      | `api.e2e-app.local`               |
| `ingress.hosts[0].paths[0].path`  | Path within the host                                                                                               | `/`                               |
| `resources.limits.cpu`            | CPU resource limits                                                                                                | `500m`                            |
| `resources.limits.memory`         | Memory resource limits                                                                                             | `256Mi`                           |
| `resources.requests.cpu`          | CPU resource requests                                                                                              | `100m`                            |
| `resources.requests.memory`       | Memory resource requests                                                                                           | `128Mi`                           |
| `config.appEnv`                   | Application environment                                                                                            | `production`                      |
| `config.logLevel`                 | Logging level                                                                                                      | `info`                            |
| `database.enabled`                | Enable database installation                                                                                       | `true`                            |
| `database.host`                   | Database host                                                                                                      | `postgres`                        |
| `database.port`                   | Database port                                                                                                      | `5432`                            |
| `database.name`                   | Database name                                                                                                      | `e2e_app`                         |
| `database.user`                   | Database user                                                                                                      | `postgres`                        |
| `database.existingSecret`         | Secret containing the database password                                                                            | `auth-db-credentials`             |
| `jwt.secret`                      | JWT signing secret                                                                                                 | `your-secret-key`                 |
| `jwt.expiration`                  | JWT token expiration                                                                                               | `24h`                             |
| `jwt.refreshExpiration`           | JWT refresh token expiration                                                                                       | `168h`                            |

## Dependencies

This chart depends on:

- PostgreSQL chart from Bitnami (required for database functionality)

## Environment-specific Values

This chart includes environment-specific values files:

- `values-dev.yaml`: Development environment configuration
- `values-prod.yaml`: Production environment configuration

To deploy with environment-specific values:

```bash
# Development deployment
helm install e2e-app /path/to/charts/e2e-app -f /path/to/charts/e2e-app/values-dev.yaml

# Production deployment
helm install e2e-app /path/to/charts/e2e-app -f /path/to/charts/e2e-app/values-prod.yaml
```