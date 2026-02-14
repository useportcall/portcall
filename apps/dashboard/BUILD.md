# Dashboard Build & Deployment Guide

## Environment-Specific Builds

The dashboard frontend requires the Keycloak URL to be set at **build time** (not runtime). Choose the appropriate build method for your target environment.

### 1. Local Development (http://localhost:8080)

```bash
# Frontend only (hot reload)
cd apps/dashboard/frontend
npm run dev

# Or build locally
npm run build:local
```

### 2. Local Kubernetes (http://localhost:30800)

```bash
# Build Docker image for local k8s
docker build -f apps/dashboard/Dockerfile \
  --build-arg VITE_KEYCLOAK_URL=http://localhost:30800 \
  -t portcall/dashboard:local \
  .

# Or use the frontend script
cd apps/dashboard/frontend
npm run build:k8s-local
```

### 3. Staging Environment (https://auth.staging.useportcall.com)

```bash
# Build Docker image for staging
docker build -f apps/dashboard/Dockerfile \
  --build-arg VITE_KEYCLOAK_URL=https://auth.staging.useportcall.com \
  -t registry.digitalocean.com/portcall-registry/portcall-dashboard:staging \
  .

# Push to registry
docker push registry.digitalocean.com/portcall-registry/portcall-dashboard:staging

# Deploy to k8s
kubectl set image deployment/dashboard \
  dashboard=registry.digitalocean.com/portcall-registry/portcall-dashboard:staging \
  -n portcall-staging
```

### 4. Production (https://auth.useportcall.com)

```bash
# Build Docker image for production
docker build -f apps/dashboard/Dockerfile \
  --build-arg VITE_KEYCLOAK_URL=https://auth.useportcall.com \
  -t registry.digitalocean.com/portcall-registry/portcall-dashboard:v1.0.x \
  .

# Push to registry
docker push registry.digitalocean.com/portcall-registry/portcall-dashboard:v1.0.x

# Deploy via Helm
cd k8s
helm upgrade portcall ./portcall-chart \
  -n portcall \
  -f ./deploy/digitalocean/values.yaml \
  --set dashboard.image.tag=v1.0.x
```

## Important Notes

### ⚠️ VITE_KEYCLOAK_URL is a Build-Time Variable

- The Keycloak URL is **embedded** into the frontend JavaScript during build
- Changing environment variables after build will NOT update the URL
- Always rebuild with the correct `--build-arg` for your target environment

### Verifying the Build

Check which Keycloak URL is in a built image:

```bash
# For Docker image
docker run --rm portcall/dashboard:latest \
  find /app/frontend/dist -name "*.js" -exec grep -o "https://auth[^\"]*" {} \; | head -1

# For deployed pod
kubectl exec -n portcall deployment/dashboard -- \
  grep -r "auth\." /app/frontend/dist/ | grep -o "https://auth[^\"]*" | head -1
```

## Package.json Scripts

- `npm run dev` - Local dev server (localhost:8080 Keycloak)
- `npm run build` - Base build command (requires VITE_KEYCLOAK_URL env var)
- `npm run build:local` - Build for localhost development
- `npm run build:k8s-local` - Build for local Kubernetes
- `npm run build:staging` - Build for staging environment
- `npm run build:production` - Build for production environment

## Keycloak Realm Configuration

Ensure your Keycloak realm has the correct redirect URIs configured:

```json
{
  "redirectUris": [
    "http://localhost:5400/*",
    "http://localhost:8082/*",
    "http://localhost:30082/*",
    "https://dashboard.staging.useportcall.com/*",
    "https://dashboard.useportcall.com/*"
  ]
}
```

## Troubleshooting

### "Invalid parameter: redirect_uri" Error

This means the frontend is using a different Keycloak URL than expected:

1. Check what URL is in the deployed frontend:

   ```bash
   kubectl exec -n portcall deployment/dashboard -- \
     grep -r "keycloak" /app/frontend/dist/ | grep "url:"
   ```

2. Rebuild with the correct `--build-arg VITE_KEYCLOAK_URL`

3. Verify the Keycloak realm has the redirect URI whitelisted

### Docker Build Using Cache

Force a fresh build:

```bash
docker build --no-cache -f apps/dashboard/Dockerfile \
  --build-arg VITE_KEYCLOAK_URL=https://auth.useportcall.com \
  -t portcall/dashboard:latest \
  .
```

Or clean the frontend build artifacts first:

```bash
cd apps/dashboard/frontend
rm -rf dist node_modules/.vite
cd ../../..
docker build -f apps/dashboard/Dockerfile \
  --build-arg VITE_KEYCLOAK_URL=https://auth.useportcall.com \
  -t portcall/dashboard:latest \
  .
```
