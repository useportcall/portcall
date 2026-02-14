# Portcall

<p align="center">
	<img height="166" alt="portcall-read-me-banner" src="https://github.com/user-attachments/assets/15e25206-ff6a-4753-bf5c-66b72b7cf33f" />
</p>

<p align="center">
	<b>Open-source, developer-first platform for metered billing, entitlements, and feature management.</b><br>
	<a href="https://useportcall.com">Website</a> · <a href="https://useportcall.com/docs">Docs</a>
</p>

---

## 🚢 What is Portcall?

Portcall is a modern, open-source platform for building, launching, and scaling SaaS products with usage-based billing, entitlement management, and feature flagging.

**Key features:**

- 🔌 **Metered & subscription billing** (usage-based, seat-based, and more)
- 🛡️ **Entitlement management** (feature flags, quotas, limits)
- ⚡ **Modern APIs** (REST, webhooks, event-driven)
- 🖥️ **Beautiful dashboard** (Vite + React)
- 🧩 **Monorepo**: Go backend, TypeScript/React frontends, Dockerized services

---

## 🚀 Quick Start

### Prerequisites

- [Docker Desktop](https://www.docker.com/) (running)
- [Go 1.20+](https://golang.org/doc/install)
- [Node.js LTS](https://nodejs.org/) + [pnpm](https://pnpm.io/) (`npm install -g pnpm`)

### First-Time Setup

```bash
git clone https://github.com/useportcall/portcall.git
cd portcall

# Build and run the dev CLI
cd tools/dev-cli && go build -o ../../dev-cli . && cd ../..

# One-command setup (installs deps, builds SDK, starts services)
./dev-cli setup
```

### Daily Development

```bash
# Start dashboard development (dashboard in terminal, infra in Docker)
./dev-cli run --preset=dashboard

# Quick mode (minimal - just api + dashboard)
./dev-cli run --preset=quick

# See all options
./dev-cli run --list

# Stop everything
./dev-cli stop
```

### Email E2E Checks (Local + Live)

```bash
# Deterministic local e2e (mock Resend API)
make e2e-email-local

# Live Resend e2e (sends real emails)
RESEND_API_KEY=... \
E2E_EMAIL_FROM=relay-test@mail.useportcall.com \
E2E_EMAIL_TO=hello@useportcall.com \
make e2e-email-live
```

`make e2e-email-local` covers:
- Email worker transactional flow (invoice + status email tasks)
- SMTP relay flow (password-reset style SMTP message)

### Discord Notification E2E Checks (Local + Live)

```bash
# Local deterministic run (uses in-process webhook capture)
make e2e-discord
make e2e-browser-discord

# Live run (sends real Discord messages)
cp apps/api/.envs.example apps/api/.envs
cp apps/dashboard/.envs.example apps/dashboard/.envs
cp apps/billing/.envs.example apps/billing/.envs
# Fill DISCORD_WEBHOOK_URL_SIGNUP and DISCORD_WEBHOOK_URL_BILLING in each .envs
make e2e-discord-live
make e2e-browser-discord-live
```

### Available Presets

| Preset | Description |
|--------|-------------|
| `dashboard` | Dashboard + checkout in terminal, others in Docker |
| `quick` | Minimal - just API in Docker, dashboard in terminal |
| `billing` | Billing worker development |
| `all-docker` | All apps in Docker containers |
| `minimal` | Infrastructure only (no apps) |

---

## 🗂️ Project Structure

```
portcall/
├── apps/              # Main backend and frontend apps
│   ├── api/           # Public REST API (port 8080)
│   ├── dashboard/     # Go backend + Vite/React frontend (port 8082)
│   ├── checkout/      # Go backend + Next.js frontend (port 8700)
│   ├── admin/         # Admin API (port 8081)
│   └── ...            # billing, email, cron workers
├── libs/              # Shared Go libraries
├── docker-compose/    # Docker Compose files
├── observability/     # Grafana/Loki/Promtail access + config guide
├── example/           # Example Next.js apps
└── tools/dev-cli/     # Development CLI
```

---

## 🏛️ Architecture

- **Go microservices**: Modular, scalable, event-driven
- **Frontend**: Vite+React dashboard, Next.js checkout
- **Database**: Postgres · **Auth**: Keycloak · **Queue**: Redis

---

## 📦 Example Apps

- [`example/example-next-app`](./example/example-next-app): Next.js demo for Portcall integration

---

## 🤝 Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on code style, testing, and PR process.

---

## 📚 Documentation

- [Documentation](https://useportcall.com/docs)
- [Website](https://useportcall.com)

---

## 🛡️ License

[Apache 2.0 License](./LICENSE)

<p align="center">
	<a href="https://github.com/useportcall/portcall/actions"><img src="https://github.com/useportcall/portcall/workflows/CI/badge.svg" alt="CI Status"></a>
	<a href="https://github.com/useportcall/portcall/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
</p>
