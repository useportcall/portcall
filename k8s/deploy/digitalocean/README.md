# DigitalOcean Production Values

This directory holds production overrides for the DigitalOcean cluster.

## Canonical Deploy Path

Use `tools/dev-cli` from repo root:

```bash
go run ./tools/dev-cli deploy --cluster digitalocean --version patch
```

`--cluster digitalocean` resolves to:

- Kubernetes context: `your-k8s-context` (or an alias-specific context from `.dev-cli.infra.json`)
- Values template: `k8s/deploy/digitalocean/values.example.yaml` (copy to `values.yaml` locally)
- Helm chart: `k8s/portcall-chart`

## Bootstrap A New Cluster

For open-source/self-hosted setup, prefer:

```bash
go run ./tools/dev-cli infra init --cluster digitalocean --mode micro
go run ./tools/dev-cli infra doctor --cluster digitalocean
```

This keeps production defaults intact while generating a micro-mode values file for the new cluster.

## Secrets Management

Use `dev-cli` for secret updates instead of ad-hoc scripts:

```bash
go run ./tools/dev-cli secrets --cluster digitalocean --namespace portcall list
go run ./tools/dev-cli secrets --cluster digitalocean --namespace portcall get portcall-secrets
go run ./tools/dev-cli secrets --cluster digitalocean --namespace portcall set portcall-secrets KEY=value
```

Required production keys are expected in `portcall-secrets` (database, Redis, Keycloak, S3, email, Grafana alerting, etc.).

## Safe Deploy Checklist

1. Confirm context: `kubectl config current-context`
2. Confirm tests: unit/integration enabled (E2E optional)
3. Run deploy: `go run ./tools/dev-cli deploy --cluster digitalocean`
4. Verify rollout:
   - `kubectl get deployment -n portcall`
   - `kubectl get pods -n portcall`
   - `kubectl rollout status deployment/api -n portcall`
   - `kubectl rollout status deployment/dashboard -n portcall`

## Notes

- Keep credentials out of `values.yaml`.
- Keep image tags explicit (avoid `latest`) for reproducible rollouts.
- Update this file only for environment-level config, not transient hotfixes.
