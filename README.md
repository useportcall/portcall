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

Portcall is a modern, open-source platform for building, launching, and scaling SaaS products with usage-based billing, entitlement management, and feature flagging. Built for developers, Portcall provides robust APIs, a beautiful dashboard, and ready-to-use example apps.

**Key features:**

- 🔌 **Metered & subscription billing** (usage-based, seat-based, and more)
- 🛡️ **Entitlement management** (feature flags, quotas, limits)
- ⚡ **Modern APIs** (REST, webhooks, event-driven)
- 🖥️ **Beautiful dashboard** (Vite + React)
- 🧩 **Monorepo**: Go backend, TypeScript/React frontends, Dockerized services
- 🧪 **Ready-to-use example apps** (Next.js, with more coming soon...)

---

## 🗂️ Monorepo Structure

```
portcall/
├── apps/           # Main backend and frontend apps
│   ├── api/        # Go REST API for billing, entitlements, subscriptions
│   ├── dashboard/  # Go backend with Vite+React frontend dashboard
│   ├── checkout/   # Go backend with Next.js frontend checkout
│   ├── billing/    # Go Billing worker microservice
│   ├── ...         # Other Go worker microservices (email, file, webhook, etc)
├── libs/           # Shared Go libraries (dbx, apix, authx, etc)
├── docker/         # Docker Compose, infra, and local dev tools
├── example/        # Example Next.js app for integration
├── CONTRIBUTING.md # Contribution guidelines
├── LICENSE         # Apache 2.0
└── README.md       # This file
```

---

## 🚀 Quick Start

### 1. Prerequisites

- [Go 1.20+](https://golang.org/doc/install)
- [Node.js (LTS)](https://nodejs.org/)
- [Docker](https://www.docker.com/)

### 2. Clone & Bootstrap

```bash
git clone https://github.com/useportcall/portcall.git
cd portcall
```

### 3. Run Everything (Local Dev)

```bash
# Start all services (API, dashboard, DB, etc)
cd docker
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.tools.yml -f docker-compose.workers.yml up
```

### 4. Frontend Apps

```bash
# Dashboard (Vite+React)
cd apps/dashboard/frontend
npm install && npm run dev

# Example Next.js app
cd example/example-next-app
npm install && npm run dev
```

### 5. API (Go)

```bash
cd apps/api
go run main.go
```

---

## 🏛️ Architecture

- **Go microservices**: Modular, scalable, and event-driven
- **Frontend**: Vite+React dashboard, Next.js checkout & example apps
- **Database**: Postgres (Dockerized for local dev)
- **Auth**: Keycloak (Dockerized), JWT, API keys
- **Queue**: Background jobs via Redis
- **Observability**: Loki, Promtail, Prisma Studio
- **CI/CD**: GitHub Actions (coming soon)

---

## 📦 Example Apps

- [`example/example-next-app`](./example/example-next-app): Next.js demo for integrating Portcall billing & entitlements
- [`apps/checkout/frontend`](./apps/checkout/frontend): Checkout UI
- [`apps/dashboard/frontend`](./apps/dashboard/frontend): Dashboard

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for:

- Code of Conduct
- How to contribute
- Development setup (Go, Node.js, Docker)
- Coding standards
- Issue reporting & PR process

---

## 📚 Documentation & Community

- [Documentation](https://useportcall.com/docs)
- [Website](https://useportcall.com)

---

## 🛡️ License

Portcall is licensed under the [Apache 2.0 License](./LICENSE).

---

<p align="center">
	<a href="https://github.com/useportcall/portcall/actions"><img src="https://github.com/useportcall/portcall/workflows/CI/badge.svg" alt="CI Status"></a>
	<a href="https://github.com/useportcall/portcall/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
</p>
