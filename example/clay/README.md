# Portcall Next.js Example App (Clay)

A complete example demonstrating how to integrate Portcall for billing, subscriptions, and entitlement management in a Next.js application.

## 🚀 Quick Start

> **🆕 First time? Pulling fresh changes?** See [QUICKSTART.md](../../QUICKSTART.md) in the repo root for complete setup instructions.

### Prerequisites

- [Node.js](https://nodejs.org/) (LTS version recommended)
- [pnpm](https://pnpm.io/) - Install with `npm install -g pnpm`
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) - Required for running Portcall services

### Automated Setup (Recommended)

From the **monorepo root** (`portcall/`):

```bash
# First time setup (installs deps, builds SDK, starts services)
bash scripts/first-time-setup.sh
```

### Manual Setup

If you prefer to run steps manually or have already done the initial setup:

#### 1. Start Portcall Services

From the **monorepo root** (`portcall/`):

```bash
cd docker-compose

# IMPORTANT: Rebuild images if you pulled new changes
docker compose -f docker-compose.db.yml \
  -f docker-compose.auth.yml \
  -f docker-compose.apps.yml \
  build

# Start services
docker compose -f docker-compose.db.yml \
  -f docker-compose.auth.yml \
  -f docker-compose.apps.yml \
  up -d

# Wait for services to start
sleep 30
```

This starts:

- PostgreSQL database
- Redis
- Keycloak auth
- Portcall API (port 9080)
- Portcall Dashboard (port 8082)

> **⚠️ Important**: If you get 404 errors later, the API is running old code. Run `bash scripts/rebuild-api.sh` to rebuild it.

#### 2. Install Dependencies & Build SDK

From the **monorepo root** (only needed once or after pulling SDK changes):

```bash
pnpm install
cd sdks/node-typescript
pnpm build
cd ../..
```

#### 3. Get Your API Key

1. Open the Portcall dashboard: http://localhost:8082
2. Create a new app or use an existing one
3. Copy the API secret key (starts with `sk_`)

#### 4. Setup Clay Plans & Features

From the **clay example directory** (`example/clay/`):

```bash
cd example/clay

# Run the setup script to create plans and features
bash scripts/setup-plans.sh
```

If you get 404 errors, the API needs to be rebuilt:

```bash
# From the monorepo root
bash scripts/rebuild-api.sh
```

#### 5. Generate Portcall Types

From the **clay example directory** (`example/clay/`):

```bash
npx portcall generate
```

The CLI will prompt you for:

1. **Select environment**: Choose "Local (localhost)" for development
2. **Enter API port**: Press Enter to use default (9080)
3. **Enter Portcall API secret key**: Paste the key you copied from the dashboard

This will:

- Generate typed SDK client in `src/portcall/client.ts`
- Generate YAML config reference in `src/portcall/portcall.yaml`
- Create/update `.env` with `PC_API_SECRET`

#### 6. Run the App

Still in the clay example directory:

```bash
pnpm dev
```

The app will be available at http://localhost:3000

## 📁 Project Structure

```
example/clay/
├── src/
│   ├── app/               # Next.js app routes
│   │   ├── actions/       # Server actions (uses Portcall SDK)
│   │   ├── login/         # Login page
│   │   └── pricing/       # Pricing page
│   ├── components/        # React components
│   ├── lib/               # Utilities
│   │   └── portcall.ts    # Server-side Portcall client
│   └── portcall/          # Auto-generated Portcall SDK
│       ├── client.ts      # Typed client with plan/entitlement helpers
│       ├── index.ts       # Re-exports
│       └── portcall.yaml  # Human-readable config reference
└── .env                   # Environment variables
```

## 🔑 Environment Setup

The `.env` file only needs one variable:

```bash
PC_API_SECRET="sk_your_api_key_here"
```

This is automatically created/updated when you run `npx portcall generate`.

## 🎯 Key Features Demonstrated

### 1. Type-Safe Plan Subscriptions

```typescript
import { portcall } from "@/lib/portcall";

// Subscribe user to a plan (type-safe!)
await portcall.plans.free.subscribe("user@email.com");
```

### 2. Entitlement Checks

```typescript
// Check if user has access to a feature
const hasAccess =
  await portcall.entitlements.ai_claygent.isEnabled("user@email.com");

// Get usage and quota
const usage = await portcall.entitlements.credits.getUsage("user@email.com");
```

### 3. Subscription Management

```typescript
// Get active subscription
const subscription = await portcall.subscriptions.getActive(userId);

// Update subscription plan
await portcall.subscriptions.update(subscriptionId, { plan_id: newPlanId });

// Cancel subscription
await portcall.subscriptions.cancel(subscriptionId);
```

### 4. Metered Billing

```typescript
// Record usage events
await portcall.meterEvents.record({
  user_id: userId,
  feature_id: "tokens",
  quantity: 100,
});
```

## 🔄 Regenerating Types

Whenever you add/update plans or features in the Portcall dashboard, regenerate types.

From the clay example directory (`example/clay/`):

```bash
npx portcall generate
```

The CLI will:

- Read `PC_API_SECRET` from your `.env` file (or prompt if missing)
- Generate types in `src/portcall/` (default output directory)
- Create or update `.env` with your API key

## 📖 API Reference

### Core Endpoints Used

```bash
# Create checkout session
POST /v1/checkout-sessions

# Get user subscriptions
GET /v1/subscriptions?user_id={user_id}

# Update subscription
POST /v1/subscriptions/{subscription_id}

# Cancel subscription
POST /v1/subscriptions/{subscription_id}/cancel

# Get user entitlement
GET /v1/entitlements/{user_id}/{entitlement_id}

# Record meter event
POST /v1/meter-events
```

## 🛠️ Development

### Running the Local API

Make sure the Portcall API is running:

````bStopping Services

To stop the Docker services:

```bash
cd ../../docker-compose
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.apps.yml down
````

### Viewing Logs

To see API logs:

```bash
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.apps.yml logs -f api
```

### Troubleshooting

**"Cannot find module 'portcall'" or CLI not found:**

- Make sure you ran `pnpm install` from the monorepo root
- Make sure you built the SDK: `cd sdks/node-typescript && pnpm build`

**Docker services won't start:**

- Make sure Docker Desktop is running
- Try `docker compose down -v` then start again (⚠️ this deletes data)

**API connection errors:**

- Make sure Docker services are running: `docker ps` should show containers
- Check API is accessible: `curl http://localhost:9080/ping`

### Running the Local API

The Portcall API runs in Docker at http://localhost:9080

If you need to restart just the API:

````bash
cd ../../docker-compose
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.apps.yml restart api
```8082 to:

- Create/manage plans
- View subscriptions
- Configure features and entitlements

## 📝 Notes

- This example uses Next.js 16 with App Router
- Server Actions are used for all Portcall API calls
- The generated SDK provides full TypeScript type safety
- Plans and features are synced from your Portcall dashboard

## 🔗 Learn More

- [Portcall Documentation](https://useportcall.com/docs)
- [Next.js Documentation](https://nextjs.org/docs)
- [Portcall GitHub](https://github.com/useportcall/portcall)
````
