# dev-cli

`dev-cli` is the operational CLI for local development, infrastructure lifecycle, and Kubernetes deploys.

It is designed to be:
- Idempotent: reruns reconcile state instead of blindly recreating resources.
- Interactive by default: `infra create` / `infra update` ask for confirmation and stage selection.
- Safe: guardrails block risky defaults and require stronger confirmation when existing infra is detected.
- Guided in-CLI: domain, stage, safety, and DNS setup prompts are built into the flow (README is optional).

Start the interactive mode:

```bash
portcall
```

`dev-cli` remains a supported alias.
Repository-local wrapper is available at `./portcall`.

## Quick Start (DigitalOcean)

1. Start infra flow (it can guide `doctl auth init/switch` interactively if needed):

```bash
portcall
```

2. Preview infra changes and preflight checks (no resources created):

```bash
go run ./tools/dev-cli infra create --cluster digitalocean --name portcall-micro --domain yourdomain.com --dry-run
```

3. Apply core infra (cluster/network/registry):

```bash
go run ./tools/dev-cli infra create --cluster digitalocean --name portcall-micro --domain yourdomain.com --step core
```

4. Apply managed services (Postgres -> Redis -> Spaces):

```bash
go run ./tools/dev-cli infra update --cluster digitalocean --name portcall-micro --domain yourdomain.com --step services
```

5. Validate cluster wiring before deploy:

```bash
go run ./tools/dev-cli infra doctor --cluster digitalocean
```

6. Deploy applications:

```bash
go run ./tools/dev-cli deploy --cluster digitalocean --apps all --version patch
```

7. Verify rollout in namespace `portcall`:

```bash
kubectl get deployment -n portcall
kubectl get pods -n portcall
kubectl rollout status deployment/portcall-api -n portcall
```

## Safety Defaults

- `infra create` and `infra update` are interactive unless `--yes` is passed.
- If an existing cluster is detected in `create` mode, the CLI requires explicit typed confirmation.
- Placeholder domain guardrail: non-dry-run apply with `example.com` and `--yes` is blocked.
- DigitalOcean preflight validates:
  - required CLIs,
  - resolved token/auth context,
  - required read endpoints,
  - write-permission probes (non-destructive failure probes).

## DNS Readiness

After infra apply, the CLI prints a DNS checklist with required hostnames:
- `admin.<domain>`
- `dashboard.<domain>`
- `api.<domain>`
- `auth.<domain>`
- `quote.<domain>`
- `checkout.<domain>`
- `webhook.<domain>`
- `file.<domain>`

These must point to the ingress load balancer before endpoints are reachable.

Cloudflare assist is available directly in the init flow:
- Interactive mode asks whether to check Cloudflare CLI auth and DNS records.
- `--dns-provider cloudflare` enables Cloudflare checks automatically.
- `--dns-auto` creates/updates missing records idempotently via Cloudflare API.
- `--cloudflare-zone-id` can pin zone resolution when needed.

Cloudflare auth options:
- preferred: interactive prompt in `infra create/update` (can save token in `.dev-cli.auth.json`)
- optional env var:

```bash
export CLOUDFLARE_API_TOKEN=<token>
```

Optional Cloudflare CLI auth check:

```bash
wrangler login
```

## Idempotent Workflow Notes

- Use `infra create` for new stacks.
- Use `infra update` for existing stacks.
- Re-running either command is expected and supported.
- Use `infra status` to inspect saved local alias state:

```bash
go run ./tools/dev-cli infra status --cluster digitalocean
```

- If a cluster already exists but local state is missing, pull it first:

```bash
go run ./tools/dev-cli infra pull --cluster digitalocean
```

## Additional Commands

Cleanup deprecated resources:

```bash
go run ./tools/dev-cli infra cleanup legacy --cluster digitalocean --dry-run
go run ./tools/dev-cli infra cleanup legacy --cluster digitalocean --yes
```

Downscale playground/test stacks:

```bash
go run ./tools/dev-cli infra downscale --cluster digitalocean
```

## Validation

Run CLI tests from module root:

```bash
cd tools/dev-cli
GOCACHE=/tmp/go-build-cache go test ./...
```
