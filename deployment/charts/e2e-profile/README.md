# E2E Profile Service Helm Chart

This Helm chart deploys the E2E Profile Service and its PostgreSQL database to a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure
- The e2e-app service deployed in the same cluster

## Installing the Chart

To install the chart with the release name `my-profile`:

```bash
helm install my-profile /path/to/e2e-profile
```

## Configuration

The following table lists the configurable parameters of the e2e-profile chart and their default values.

| Parameter                         | Description                                                                                                        | Default                           |
|-----------------------------------|--------------------------------------------------------------------------------------------------------------------|-----------------------------------|
| `replicaCount`                    | Number of replicas                                                                                                 | `1`                               |
| `image.repository`                | Image repository                                                                                                   | `e2e-profile`                     |
| `image.tag`                       | Image tag                                                                                                          | `latest`                          |
| `image.pullPolicy`                | Image pull policy                                                                                                  | `IfNotPresent`                    |
| `service.type`                    | Kubernetes Service type                                                                                            | `ClusterIP`                       |
| `service.port`                    | Service HTTP port                                                                                                  | `8081`                            |
| `ingress.enabled`                 | Enable ingress controller resource                                                                                 | `true`                            |
| `ingress.className`               | IngressClass that will be be used                                                                                  | `nginx`                           |
| `ingress.hosts[0].host`           | Hostname to your installation                                                                                      | `profile.e2e-app.local`           |
| `ingress.hosts[0].paths[0].path`  | Path within the host                                                                                               | `/`                               |
| `resources.limits.cpu`            | CPU resource limits                                                                                                | `200m`                            |
| `resources.limits.memory`         | Memory resource limits                                                                                             | `256Mi`                           |
| `resources.requests.cpu`          | CPU resource requests                                                                                              | `100m`                            |
| `resources.requests.memory`       | Memory resource requests                                                                                           | `128Mi`                           |
| `env`                             | Environment variables                                                                                              | See `values.yaml`                 |
| `database.enabled`                | Enable database installation                                                                                       | `true`                            |
| `database.host`                   | Database host                                                                                                      | `postgres-profile`                |
| `database.port`                   | Database port                                                                                                      | `5432`                            |
| `database.name`                   | Database name                                                                                                      | `e2e_profile`                     |
| `database.user`                   | Database user                                                                                                      | `postgres`                        |
| `database.existingSecret`         | Secret containing the database password                                                                            | `profile-db-credentials`          |

## Dependencies

This chart depends on:

- PostgreSQL chart from Bitnami (required for database functionality)
- The e2e-app service should be deployed in the same cluster

## Environment-specific Values

This chart includes environment-specific values files:

- `values-dev.yaml`: Development environment configuration
- `values-prod.yaml`: Production environment configuration

To deploy with environment-specific values:

```bash
# Development deployment
helm install e2e-profile /path/to/charts/e2e-profile -f /path/to/charts/e2e-profile/values-dev.yaml

# Production deployment
helm install e2e-profile /path/to/charts/e2e-profile -f /path/to/charts/e2e-profile/values-prod.yaml
```

## Deployment Example

To deploy both the e2e-app and e2e-profile services in a development environment:

```bash
# First deploy e2e-app
helm install e2e-app /path/to/charts/e2e-app -f /path/to/charts/e2e-app/values-dev.yaml

# Then deploy e2e-profile
helm install e2e-profile /path/to/charts/e2e-profile -f /path/to/charts/e2e-profile/values-dev.yaml
```