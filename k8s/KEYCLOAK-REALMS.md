# Keycloak Realms Configuration for Production

## Overview

Portcall uses **two separate Keycloak realms**:

1. **`admin` realm** - For the Admin Portal (admin app)
2. **`dev` realm** - For the Dashboard and Checkout apps

Both realms are baked into the custom Keycloak Docker image during build time.

## Realm Configuration

### Admin Realm (`admin`)

**Purpose:** Superadmin access to the Admin Portal

- **Client ID:** `portcall-admin`
- **Theme:** `portcall-admin` (purple branding)
- **User:** `superadmin` with fixed password (change in production!)
- **Allowed Redirects:**
  - `http://localhost:5401/*` (local dev server)
  - `http://localhost:8081/*` (local production build)
  - `http://localhost:30081/*` (Kubernetes NodePort)
  - `https://admin.useportcall.com/*` (production)

**Configuration File:** `docker-compose/keycloak/realms/admin-realm.json`

### Dev Realm (`dev`)

**Purpose:** User authentication for Dashboard and Checkout

- **Client ID:** `portcall`
- **Theme:** `portcall` (greyscale branding)
- **Allowed Redirects:**
  - `http://localhost:5400/*` (local dev server)
  - `http://localhost:8082/*` (dashboard local)
  - `http://localhost:30082/*` (dashboard K8s)
  - `http://localhost:30700/*` (checkout K8s)
  - `https://dashboard.useportcall.com/*` (production)
  - `https://checkout.useportcall.com/*` (production)

**Configuration File:** `docker-compose/keycloak/realms/dev-realm.json`

## Docker Image Build

The custom Keycloak image includes:

1. **Base image:** `quay.io/keycloak/keycloak:26.3.2`
2. **Custom themes:**
   - `portcall` (greyscale)
   - `portcall-admin` (purple)
3. **Realm configurations:**
   - `admin-realm.json`
   - `dev-realm.json`

### Dockerfile Location

`docker-compose/keycloak/Dockerfile`

```dockerfile
FROM quay.io/keycloak/keycloak:26.3.2

# Copy custom Portcall themes
COPY themes/portcall /opt/keycloak/themes/portcall
COPY themes/portcall-admin /opt/keycloak/themes/portcall-admin

# Copy realm configurations
COPY realms/dev-realm.json /opt/keycloak/data/import/dev-realm.json
COPY realms/admin-realm.json /opt/keycloak/data/import/admin-realm.json

# Build optimized Keycloak image
RUN /opt/keycloak/bin/kc.sh build
```

## Building for Production

### 1. Build And Deploy Keycloak With Dev CLI

```bash
cd /path/to/portcall

# Deploy Keycloak image/chart changes to DigitalOcean production
go run ./tools/dev-cli deploy --cluster digitalocean --apps keycloak --version patch
```

### 2. Tag and Push to Registry

```bash
# Tag images for your registry
docker tag portcall/keycloak:latest your-registry.example.com/portcall-keycloak:latest
docker tag portcall/admin:latest your-registry.example.com/portcall-admin:latest
docker tag portcall/dashboard:latest your-registry.example.com/portcall-dashboard:latest

# Push to registry
docker push your-registry.example.com/portcall-keycloak:latest
docker push your-registry.example.com/portcall-admin:latest
docker push your-registry.example.com/portcall-dashboard:latest
```

### 3. Update Helm Values

Ensure your Helm values file uses the custom Keycloak image:

```yaml
keycloak:
  enabled: true
  image: portcall/keycloak:latest # Custom image with both realms
  # ... other settings
```

## Application Configuration

### Admin App

**Backend Environment Variables:**

```bash
PORT=8081
KEYCLOAK_API_URL=http://keycloak:8080  # Internal K8s service
KEYCLOAK_REALM=admin
```

**Frontend Build Arg:**

```bash
VITE_KEYCLOAK_URL=https://auth.useportcall.com  # External URL
```

**Helm Configuration:**

```yaml
admin:
  enabled: true
  env:
    keycloakApiUrl: "https://auth.useportcall.com"
    keycloakInternalUrl: "http://keycloak:8080"
    keycloakRealm: "admin"
```

### Dashboard App

**Backend Environment Variables:**

```bash
PORT=8082
KEYCLOAK_API_URL=http://keycloak:8080  # Internal K8s service
KEYCLOAK_REALM=dev  # Defaults to "dev" if not set
```

**Frontend Build Arg:**

```bash
VITE_KEYCLOAK_URL=https://auth.useportcall.com  # External URL
```

**Helm Configuration:**

```yaml
dashboard:
  enabled: true
  env:
    keycloakApiUrl: "https://auth.useportcall.com"
    keycloakInternalUrl: "http://keycloak:8080"
    keycloakRealm: "dev"
```

## Important Notes

### ✅ No Hardcoded Values in Images

- Realm JSON files contain **both localhost and production URLs**
- This is correct! Keycloak uses pattern matching (`/*` wildcards)
- Production URLs like `https://admin.useportcall.com/*` are already in the realm files
- No need to parameterize redirect URIs

### ⚠️ Environment-Specific Values

These are passed via Helm values, **NOT baked into images**:

- `KEYCLOAK_API_URL` (external URL for frontend)
- `KEYCLOAK_INTERNAL_URL` (internal K8s service URL)
- `KEYCLOAK_REALM` (which realm to use)
- Database connection strings
- Admin credentials

### 🔒 Security Considerations

1. **Change default superadmin password** in production
2. **Use strong KC_BOOTSTRAP_ADMIN_PASSWORD** for Keycloak master realm
3. **Enable HTTPS** in production (handled by ingress)
4. **Set proper CORS/Web Origins** in realm configs

## Testing in Kubernetes

```bash
# Deploy with Helm (default chart values for local/dev)
helm install portcall ./k8s/portcall-chart \
  -n portcall-dev \
  --create-namespace \
  -f ./k8s/portcall-chart/values.yaml

# Wait for pods
kubectl get pods -n portcall-dev -w

# Test realm endpoints
kubectl port-forward -n portcall-dev svc/keycloak 8080:8080
curl http://localhost:8080/realms/admin
curl http://localhost:8080/realms/dev

# Access admin portal
kubectl port-forward -n portcall-dev svc/admin 8081:8081
# Visit http://localhost:8081

# Access dashboard
kubectl port-forward -n portcall-dev svc/dashboard 8082:8082
# Visit http://localhost:8082
```

## Troubleshooting

### Realm not found

If you get "Realm not found" errors:

1. Check Keycloak logs: `kubectl logs -n portcall-dev deployment/keycloak`
2. Verify image includes realms: `docker run --rm portcall/keycloak:latest ls /opt/keycloak/data/import/`
3. Ensure `--import-realm` flag is set in deployment

### Redirect URI mismatch

If you get redirect errors:

1. Check the realm JSON includes your production URL
2. Verify frontend is using correct `VITE_KEYCLOAK_URL` build arg
3. Check backend is using correct `KEYCLOAK_REALM` env var

### Wrong theme loading

1. Verify theme is copied in Dockerfile
2. Check realm JSON has correct `loginTheme` value
3. Rebuild Keycloak image if themes were updated
