# CI/CD Pipeline - e2e-app

This document describes the CI/CD pipeline setup for the e2e-app service.

## Pipeline Overview

The CI/CD pipeline for this service uses Jenkins to automate the build, test, and deployment process. The pipeline is defined in the `Jenkinsfile` in the root of this repository.

### Pipeline Stages

1. **Checkout** - Retrieves the source code from the repository
2. **Unit Tests** - Runs unit tests to verify individual components
3. **Code Quality** - Performs static code analysis using Go tools
4. **Build Docker Image** - Creates a Docker image for the application
5. **Security Scan** - Scans the Docker image for security vulnerabilities
6. **Push Docker Image** - Pushes the Docker image to the registry
7. **Deploy to Dev** - Deploys the application to the development environment
8. **Integration Tests** - Runs integration tests against the development environment
9. **Approve Production Deployment** - Manual approval step for production deployment
10. **Deploy to Production** - Deploys the application to the production environment

## Environment Configuration

The pipeline is configured to deploy to different environments based on the branch:

- **develop branch**: Deploys automatically to the development environment
- **main branch**: Deploys to production after manual approval

## Configuration

Pipeline configuration is stored in `jenkins-config.yaml` in the root of this repository. This file defines:

- Environment configurations
- Test settings
- Security scan thresholds
- Notification settings

## Jenkins Setup Requirements

To use this pipeline, the following must be configured in Jenkins:

1. **Kubernetes Plugin** - For running the pipeline in Kubernetes pods
2. **Credentials**:
   - `docker-credentials` - Docker registry credentials
   - `kubeconfig` - Kubernetes configuration file
3. **Jenkins Pipeline Plugin** - For pipeline execution

## Running the Pipeline Manually

To trigger a manual pipeline run:

1. Navigate to the Jenkins UI
2. Select the pipeline job for this service
3. Click "Build with Parameters"
4. Provide any required parameters
5. Click "Build"

## Troubleshooting

If the pipeline fails, check:

1. Jenkins console output for specific error messages
2. Kubernetes pod logs if failures occurred in the build/test containers
3. Helm deployment logs for issues related to deployment steps