# Observability Guide

This directory is the observability home for Portcall.

## Where Config Lives

- Local Promtail config: `observability/local/promtail-config.yml`
- Local Loki config: `observability/local/loki-config.yml`
- Local Compose stack: `docker-compose/docker-compose.observability.yml`
- Kubernetes observability chart templates:
  `k8s/portcall-chart/templates/observability/`
- Kubernetes observability values (defaults):
  `k8s/portcall-chart/values.yaml`
- Kubernetes observability values (DigitalOcean example):
  `k8s/deploy/digitalocean/values.example.yaml`

## Production Access (DigitalOcean)

1. Open Grafana:
  - `https://admin.useportcall.com/grafana`

2. Log in:
   - Username: `admin`
   - Password comes from secret key `GRAFANA_ADMIN_PASSWORD` in `portcall-secrets`.

```bash
kubectl -n portcall get secret portcall-secrets -o jsonpath='{.data.GRAFANA_ADMIN_PASSWORD}' | base64 -d; echo
```

3. Open prebuilt dashboards:
   - `Portcall Logs Overview`
   - `Portcall Error Drilldown`

If ingress is unavailable, use port-forward temporarily:

```bash
kubectl -n portcall port-forward svc/grafana 33000:3000
```

Then open `http://localhost:33000` and use the same credentials.

## How To View Logs In Grafana

1. Go to `Explore`.
2. Choose data source `Loki`.
3. Run LogQL queries:

```logql
{namespace="portcall", app="dashboard"}
{namespace="portcall", app="checkout"}
{namespace="portcall", app="api"}
{namespace="portcall"} |~ "(?i)error|fatal|panic"
```

Useful metric-style query:

```logql
sum by (app) (count_over_time({namespace="portcall"}[5m]))
```

## Local Access (Compose)

From repo root:

```bash
docker compose -f docker-compose/docker-compose.observability.yml up -d
```

- Grafana: `http://localhost:3030`
- Loki: `http://localhost:3100`
- Grafana login: `admin` / `admin`

## Deploy/Update In Kubernetes

Use dev-cli:

```bash
go run ./tools/dev-cli deploy --apps observability --cluster digitalocean --values k8s/deploy/digitalocean/values.yaml --version skip --skip-build --yes
```

Quick health checks:

```bash
kubectl get deployment -n portcall grafana
kubectl get statefulset -n portcall loki
kubectl get daemonset -n portcall promtail
kubectl get pods -n portcall | rg "grafana|loki|promtail"
```
