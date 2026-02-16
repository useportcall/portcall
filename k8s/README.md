# Portcall Kubernetes Assets

This directory contains the Helm chart and deployment values used by Portcall.

## Current Source Of Truth

- Deploy orchestration: `tools/dev-cli deploy`
- Secret management: `tools/dev-cli secrets`
- Production values file: `k8s/deploy/digitalocean/values.yaml`
- Helm chart: `k8s/portcall-chart`

## Production Deploy (DigitalOcean)

```bash
go run ./tools/dev-cli deploy --cluster digitalocean --version patch
```

## New Cluster Init (DigitalOcean Micro)

```bash
go run ./tools/dev-cli infra init --cluster digitalocean --mode micro
go run ./tools/dev-cli infra doctor --cluster digitalocean
go run ./tools/dev-cli deploy --cluster digitalocean --apps all --version patch
```

`infra init` provisions VPC + cluster + registry + managed Postgres/Redis/Spaces and generates `.infra/<alias>/values.micro.yaml`.

For existing clusters, align local state first:

```bash
go run ./tools/dev-cli infra pull --cluster digitalocean
```

Default preflight behavior:

- Unit tests: enabled
- Integration tests: enabled
- E2E tests: disabled

Common flags:

- `--apps api,dashboard` deploy only selected apps
- `--skip-build` skip image build/push
- `--skip-preflight` skip pre-deploy tests
- `--yes` skip interactive confirmation

## Secrets Workflow

```bash
go run ./tools/dev-cli secrets --cluster digitalocean --namespace portcall list
go run ./tools/dev-cli secrets --cluster digitalocean --namespace portcall get portcall-secrets
go run ./tools/dev-cli secrets --cluster digitalocean --namespace portcall set portcall-secrets KEY=value
```

Do not commit plaintext credentials into values files.

## Verification Commands

```bash
kubectl get deployment -n portcall
kubectl get pods -n portcall
kubectl rollout status deployment/api -n portcall
kubectl rollout status deployment/dashboard -n portcall
```

## Directory Layout

```text
k8s/
├── deploy/
│   └── digitalocean/
│       ├── README.md
│       └── values.yaml
├── portcall-chart/
│   ├── Chart.yaml
│   ├── README.md
│   ├── values.yaml
│   └── templates/
├── deploy.sh
├── smoke-test.sh
└── test-deployment.sh
```

`k8s/deploy.sh`, `k8s/smoke-test.sh`, and `k8s/test-deployment.sh` are legacy helpers. Prefer `tools/dev-cli` for all production work.
