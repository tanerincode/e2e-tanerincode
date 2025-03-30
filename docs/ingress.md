# Kubernetes Ingress Configuration Guide

## What is Ingress?
Ingress is a Kubernetes resource that manages external access to services in your cluster. Think of it as a smart router or entry point to your applications.

## Our Ingress Setup

We've configured three layers of Ingress settings:

### 1. Default Configuration (values.yaml)
```yaml
ingress:
  enabled: false              # Disabled by default
  className: "nginx"          # Using nginx ingress controller
  annotations:
    kubernetes.io/ingress.class: nginx
  hosts:
    - host: e2e-app.local    # Default local domain
      paths:
        - path: /            # Root path
          pathType: Prefix   # Match all paths starting with /
  tls: []                    # No TLS by default
```

### 2. Development Environment (values-dev.yaml)
```yaml
ingress:
  enabled: true              # Enabled for development
  className: "nginx"
  annotations:
    kubernetes.io/ingress.class: nginx
  hosts:
    - host: dev.e2e-app.local
      paths:
        - path: /
          pathType: Prefix
  tls: []                    # No SSL in development
```

### 3. Production Environment (values-prod.yaml)
```yaml
ingress:
  enabled: true
  className: "nginx"
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: "letsencrypt-prod"    # Automatic SSL certificates
    nginx.ingress.kubernetes.io/ssl-redirect: "true"      # Force HTTPS
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
  hosts:
    - host: e2e-app.example.com    # Production domain
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: e2e-app-tls      # SSL certificate storage
      hosts:
        - e2e-app.example.com
```

## Key Features Implemented

1. **Environment Separation**
   - Development: Simple HTTP setup with local domain
   - Production: HTTPS with automatic SSL certificates

2. **SSL/TLS Configuration**
   - Development: No SSL (faster local development)
   - Production: Automatic SSL with Let's Encrypt
   - Force SSL redirect in production

3. **Path Management**
   - Using Prefix path type for flexibility
   - All paths (/*) are routed to the application

4. **Security Features**
   - SSL enforcement in production
   - Automatic certificate management
   - Proper ingress class isolation

## How It Works

1. **Local Development**
   ```
   Request → dev.e2e-app.local → Ingress → Service → Pods
   ```

2. **Production**
   ```
   Request → e2e-app.example.com → SSL Termination → Ingress → Service → Pods
   ```

## Usage

1. **Apply in Development**
   ```bash
   helm upgrade --install e2e-app ./charts/e2e-app -f values-dev.yaml
   ```

2. **Apply in Production**
   ```bash
   helm upgrade --install e2e-app ./charts/e2e-app -f values-prod.yaml
   ```

## Prerequisites

1. **Nginx Ingress Controller**
   - Must be installed in the cluster
   - Configured as the default ingress class

2. **For Production**
   - cert-manager installed for SSL
   - DNS configured for your domain
   - Let's Encrypt cluster issuer configured

## Common Operations

1. **Check Ingress Status**
   ```bash
   kubectl get ingress -n your-namespace
   ```

2. **Check SSL Certificate**
   ```bash
   kubectl get certificates -n your-namespace
   ```

3. **View Ingress Details**
   ```bash
   kubectl describe ingress e2e-app -n your-namespace
   ```

## Troubleshooting

1. **SSL Issues**
   - Check cert-manager logs
   - Verify DNS settings
   - Check certificate status

2. **Routing Issues**
   - Verify ingress controller logs
   - Check service and pod health
   - Confirm correct port configurations

3. **Domain Issues**
   - Verify DNS resolution
   - Check host configuration in values
   - Confirm ingress controller is receiving traffic