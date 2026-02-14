# API Build & Deployment Guide

## Overview

The API service includes the public REST API and serves the API Reference documentation statically at `/docs`. The API Reference is a Vite/React frontend embedded into the API service.

## Environment-Specific Builds

### 1. Local Development (http://localhost:8080)

```bash
# Backend only
cd apps/api
go run main.go

# Or build api-reference frontend separately for development
cd apps/api-reference/frontend
npm install && npm run dev  # Runs on separate port for hot reload
```

### 2. Local Kubernetes (http://localhost:30080)

```bash
# Build Docker image for local k8s
docker build -f apps/api/Dockerfile \
  -t portcall/api:local \
  .

# Deploy to local k8s
kubectl apply -f k8s/deploy/local/api-deployment.yaml
```

### 3. Staging Environment

```bash
# Build Docker image for staging
docker build -f apps/api/Dockerfile \
  -t your-registry/portcall-api:staging \
  .

# Push to registry
docker push your-registry/portcall-api:staging

# Deploy to k8s
kubectl set image deployment/api \
  api=your-registry/portcall-api:staging \
  -n portcall-staging
```

### 4. Production

```bash
# Build Docker image for production
docker build -f apps/api/Dockerfile \
  -t your-registry/portcall-api:v1.0.x \
  .

# Push to registry
docker push your-registry/portcall-api:v1.0.x

# Deploy via Helm or deployment script
cd k8s
./deploy-api.sh production v1.0.x
```

## API Reference Integration

The API Reference frontend is built during the Docker build process and served at `/docs`:

- **Development**: API runs on port 8080, docs available at `http://localhost:8080/docs`
- **Production**: Docs available at `https://your-domain.com/docs`

The integration happens in two places:

1. **Dockerfile**: Builds the frontend and copies it to `/app/api-reference/dist`
2. **main.go**: Serves static files from `/docs` path with SPA fallback

## Docker Build Process

The Dockerfile uses multi-stage builds:

1. **Stage 1 (builder)**: Builds the Go API binary
2. **Stage 2 (frontend-builder)**: Builds the API Reference frontend (Vite)
3. **Stage 3 (final)**: Combines both into a minimal Alpine image

## Verifying the Build

Check that the API Reference is included in the image:

```bash
# For Docker image
docker run --rm portcall/api:latest ls -la /app/api-reference/dist

# For deployed pod
kubectl exec deployment/api -- ls -la /app/api-reference/dist
```

## Troubleshooting

### API Reference not loading?

1. Check if the dist directory exists in the container:

   ```bash
   docker run --rm portcall/api:latest ls -la /app/api-reference/
   ```

2. Verify the frontend build succeeded:

   ```bash
   cd apps/api-reference/frontend
   npm run build
   ls -la dist/
   ```

3. Check API logs for file serving errors:
   ```bash
   kubectl logs -n portcall deployment/api --tail=100
   ```

### Docker Build Failing?

Force a fresh build:

```bash
docker build --no-cache -f apps/api/Dockerfile \
  -t portcall/api:latest \
  .
```

Or clean the frontend build artifacts first:

```bash
cd apps/api-reference/frontend
rm -rf dist node_modules/.vite
cd ../../..
docker build -f apps/api/Dockerfile \
  -t portcall/api:latest \
  .
```

## Development Workflow

### Running Locally

```bash
# Terminal 1: Run backend
cd apps/api
go run main.go

# Terminal 2: Build frontend once
cd apps/api-reference/frontend
npm run build

# Now access http://localhost:8080/docs
```

### Testing Docker Build Locally

```bash
# Build the image
docker build -f apps/api/Dockerfile -t portcall/api:test .

# Run it
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_URL=your_db_url \
  portcall/api:test

# Access http://localhost:8080/docs
```

## Environment Variables

Required environment variables for the API service:

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string (optional, for queue)

No build-time environment variables are needed for the API Reference (unlike the dashboard which needs Keycloak URL).

## Deployment Scripts

See [deploy-api.sh](../../k8s/deploy-api.sh) for automated deployment to different environments.

## Related Files

- [apps/api/main.go](./main.go) - API server with /docs routing
- [apps/api/Dockerfile](./Dockerfile) - Multi-stage build definition
- [apps/api-reference/frontend/](../api-reference/frontend/) - API Reference source
- [k8s/deploy-api.sh](../../k8s/deploy-api.sh) - Deployment automation script
