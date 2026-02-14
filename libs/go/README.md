# Go Libraries

This folder contains reusable Go modules for Portcall.

## Constructor Pattern

Constructors now return errors instead of terminating the process on bad runtime config:

- `authx.NewFromEnv()`
- `cryptox.NewFromEnv()`
- `cryptox.NewFromBase64Key(...)`
- `dbx.NewFromEnv()`
- `dbx.NewFromDSN(...)`
- `emailx.NewFromEnv()`
- `qx.NewFromEnv()`
- `storex.NewFromEnv()`

The default `New()` constructors in these modules also return `(value, error)`.

## Modules

- `apix`: API-facing DTOs/helpers
- `authx`: auth providers and JWT validation
- `cronx`: scheduled queue orchestration
- `cryptox`: encryption/token helpers
- `dbx`: database abstraction and models
- `discordx`: Discord webhook client
- `dogfoodx`: Portcall API client for dogfooding
- `emailx`: email providers (local, Resend, Postmark)
- `envx`: env loading and mode helpers
- `logx`: logging middleware/helpers
- `paymentx`: Braintree/Stripe helpers
- `qx`: queue client/server helpers
- `ratelimitx`: Redis-backed rate limiting
- `routerx`: Gin router wrappers/middleware
- `storex`: S3-compatible storage helper
- `webhookx`: webhook verification/handlers

`services/` is intentionally separate and not documented here.
