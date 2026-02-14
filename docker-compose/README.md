# Docker Compose Configuration

This directory contains Docker Compose files for running the Portcall platform.

## Available Services

### Infrastructure (`docker-compose.db.yml`)

- **postgres**: PostgreSQL database
- **redis**: Redis cache/queue
- **minio**: S3-compatible object storage (port 9000 API, 9001 console)
- **minio-init**: MinIO bucket initialization (creates `quote-signatures` and `icon-logos` buckets)

### Authentication (`docker-compose.auth.yml`)

- **keycloak**: Keycloak identity and access management

### Development Tools (`docker-compose.tools.yml`)

- **asynqmon**: Redis queue monitoring UI (port 8090)
- **mailpit**: Email testing tool (SMTP 1025, UI 8025)
- **prisma_studio**: Database browser (port 5555)
- **redis_insight**: Redis GUI (port 5540)

### Observability (`docker-compose.observability.yml`)

- **loki**: Log aggregation
- **promtail**: Log collector
- Local observability config source:
  - `../observability/local/loki-config.yml`
  - `../observability/local/promtail-config.yml`
  - Full guide: `../observability/README.md`

### Applications (`docker-compose.apps.yml`)

- **api**: Public REST API (port 8080)
- **admin**: Admin API (port 8081)
- **file-api**: File service API (port 8085)
- **quote**: Quote service (port 8010)

### Frontend Apps

- **dashboard**: Dashboard UI + API (port 8082) - `docker-compose.dashboard.yml`
- **checkout**: Checkout UI + API (port 8700) - `docker-compose.checkout.yml`

### Workers (`docker-compose.workers.yml`)

- **email_worker**: Email processing worker
- **billing_worker**: Billing processing worker

## Quick Start

### Run Infrastructure + Tools

```bash
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.tools.yml up
```

### Run Infrastructure + All Apps

```bash
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.apps.yml -f docker-compose.dashboard.yml -f docker-compose.checkout.yml up
```

### Run Everything (Infrastructure + Apps + Workers + Tools)

```bash
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.tools.yml -f docker-compose.apps.yml -f docker-compose.dashboard.yml -f docker-compose.checkout.yml -f docker-compose.workers.yml up
```

### Run All Services (including observability)

```bash
docker compose $(ls docker-compose.*.yml | sort | xargs -I{} echo -f {}) up
```

## Individual Services

### Run specific services

```bash
# Just the database and Redis
docker compose -f docker-compose.db.yml up

# Just the API
docker compose -f docker-compose.db.yml -f docker-compose.apps.yml up api

# Just the dashboard
docker compose -f docker-compose.db.yml -f docker-compose.dashboard.yml up dashboard
```

## Rebuild Services

```bash
# Rebuild all services
docker compose -f docker-compose.apps.yml build

# Rebuild specific service
docker compose -f docker-compose.apps.yml build api
```

## Stop and Clean Up

```bash
# Stop all services
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.apps.yml down

# Stop and remove volumes (⚠️ deletes data)
docker compose -f docker-compose.db.yml down -v
```

## S3/Storage Configuration

### Local Development (MinIO)

The `docker-compose.db.yml` includes MinIO for S3-compatible object storage. **Buckets are automatically created** by the `minio-init` service.

**MinIO Configuration:**

- API Port: 9000
- Console Port: 9001 (web UI)
- Root User: `admin`
- Root Password: `admin123`
- Endpoint: `http://minio:9000` (internal) or `http://localhost:9000` (external)

**Auto-Created Buckets:**

- `quote-signatures`: For storing document signatures
- `icon-logos`: For storing company logos

**Access MinIO Console:**

```bash
# Start services
docker compose -f docker-compose.db.yml up

# Open http://localhost:9001
# Login: admin / admin123
```

### Environment Variables

Services using S3 storage (dashboard, file-api, quote) receive these environment variables:

```yaml
S3_REGION: "us-east-1"
S3_ENDPOINT: "http://minio:9000"
S3_ACCESS_KEY_ID: "admin"
S3_SECRET_ACCESS_KEY: "admin123"
```

These are defined in the respective docker-compose files (e.g., `docker-compose.dashboard.yml`).

### Manual Bucket Management

If you need to create additional buckets or manage existing ones:

```bash
# Install MinIO Client
brew install minio/stable/mc  # macOS
# or download from https://min.io/docs/minio/linux/reference/minio-mc.html

# Configure alias
mc alias set local http://localhost:9000 admin admin123

# List buckets
mc ls local

# Create a new bucket
mc mb local/my-new-bucket

# Upload file
mc cp myfile.png local/icon-logos/

# List files in bucket
mc ls local/icon-logos
```

### Troubleshooting

**Buckets not created:**

```bash
# Check minio-init logs
docker compose -f docker-compose.db.yml logs minio-init

# Manually create buckets
mc mb -p local/quote-signatures
mc mb -p local/icon-logos
```

**MinIO not accessible:**

```bash
# Check MinIO health
curl http://localhost:9000/minio/health/live

# Check logs
docker compose -f docker-compose.db.yml logs minio
```

**403 Forbidden on uploaded files:**

- Buckets have default access control
- Files uploaded via the dashboard/API include public-read ACL
- Check `libs/go/storex/storex.go` for ACL settings
