# Portcall Admin Panel

A full-stack admin panel for managing and monitoring all Portcall apps, users, subscriptions, plans, quotes, and triggering queue jobs.

## Features

- **Dashboard**: Platform-wide statistics overview
- **Apps Management**: View all apps with user/subscription/plan/quote counts
- **App Details**: Deep dive into each app with tabs for users, subscriptions, plans, and quotes
- **Queue Management**: Trigger background queue jobs (billing and email queues)
- **Authentication**: Separate Keycloak realm with a single superadmin user

## Tech Stack

### Backend

- Go 1.25 with Gin framework
- PostgreSQL database (shared with main app)
- Redis for queue management

### Frontend

- Vite 6 + React 19
- TypeScript
- TailwindCSS 4
- React Query 5 (TanStack Query)
- React Router 7
- Keycloak 26 for authentication

## Running Locally

### Prerequisites

- Go 1.20+
- Node.js LTS
- Docker (for infrastructure)

### Start Infrastructure

```bash
cd docker-compose
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.tools.yml up
```

### Run Backend (Serving Built Frontend)

```bash
cd apps/admin
go run main.go
```

The backend serves both the API and the built frontend at `http://localhost:8081`.

> **Note**: The backend serves the static frontend build from `frontend/dist/`. Make sure to build the frontend first (see below).

### Build Frontend

```bash
cd apps/admin/frontend
npm install
npm run build:local
```

The built files are placed in `frontend/dist/` and automatically served by the Go backend.

### Run Frontend (Development with Hot Reload)

For frontend development with hot reload:

```bash
cd apps/admin/frontend
npm install
npm run dev
```

The frontend dev server runs at `http://localhost:5401` and proxies API requests to the backend at `http://localhost:8081`.

### Using the Dev Script

The easiest way to run the admin app locally with hot reload is using the dev script:

```bash
# From repo root
./scripts/dev-dashboard.sh --apps admin

# Or select admin interactively when prompted
./scripts/dev-dashboard.sh
```

This will:

- Start all required infrastructure (database, auth, workers) in Docker
- Start the admin backend in a Terminal window
- Start the admin frontend with hot reload in another Terminal window

## Authentication

The admin panel uses a separate Keycloak realm called `admin` with:

- **Realm**: `admin`
- **Client ID**: `portcall-admin`
- **Superadmin User**: `superadmin`
- **Password**: `X9$mK@2pLqR#7vNw!8ZbCyDe`

> ⚠️ **Important**: Change the password in production by updating the Keycloak realm configuration.

## API Endpoints

### Health

- `GET /ping` - Health check
- `GET /healthz` - Health status

### Stats

- `GET /api/stats` - Platform-wide statistics

### Apps

- `GET /api/apps` - List all apps
- `GET /api/apps/:app_id` - Get app details

### Users (per app)

- `GET /api/apps/:app_id/users` - List users for an app
- `GET /api/apps/:app_id/users/:user_id` - Get user details

### Subscriptions (per app)

- `GET /api/apps/:app_id/subscriptions` - List subscriptions
- `GET /api/apps/:app_id/subscriptions/:subscription_id` - Get subscription details

### Plans (per app)

- `GET /api/apps/:app_id/plans` - List plans
- `GET /api/apps/:app_id/plans/:plan_id` - Get plan details

### Quotes (per app)

- `GET /api/apps/:app_id/quotes` - List quotes
- `GET /api/apps/:app_id/quotes/:quote_id` - Get quote details

### Queues

- `GET /api/queues` - List available queues and tasks
- `POST /api/queues/enqueue` - Enqueue a task

## Docker

Build and run with Docker:

```bash
# From repo root
docker compose -f docker-compose/docker-compose.db.yml \
               -f docker-compose/docker-compose.auth.yml \
               -f docker-compose/docker-compose.admin.yml up
```

## Environment Variables

| Variable             | Description                  | Default                 |
| -------------------- | ---------------------------- | ----------------------- |
| `PORT`               | Server port                  | `8081`                  |
| `DATABASE_URL`       | PostgreSQL connection string | Required                |
| `REDIS_ADDR`         | Redis address                | Required                |
| `AES_ENCRYPTION_KEY` | Encryption key for secrets   | Required                |
| `KEYCLOAK_URL`       | Keycloak server URL          | `http://localhost:8080` |
| `KEYCLOAK_REALM`     | Keycloak realm               | `admin`                 |
| `KEYCLOAK_CLIENT_ID` | Keycloak client ID           | `portcall-admin`        |

    "date_paid": "15th November, 2025",
    "company_name": "The Prancing Pony",
    "recipient_name": "Frodo Baggins"

}'

```

```
