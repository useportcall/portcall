# Portcall Helm Chart

This chart defines all Kubernetes resources for Portcall services.

## Use Through Dev CLI

Production deployments should go through `tools/dev-cli deploy`, which wraps Helm and rollout checks.

```bash
go run ./tools/dev-cli deploy --cluster digitalocean
```

## Manual Helm Usage (When Needed)

Render for review:

```bash
helm template portcall ./k8s/portcall-chart -f ./k8s/deploy/digitalocean/values.yaml
```

Install/upgrade:

```bash
helm upgrade --install portcall ./k8s/portcall-chart \
  --namespace portcall \
  --create-namespace \
  -f ./k8s/deploy/digitalocean/values.yaml
```

## Key Value Files

- Base defaults: `k8s/portcall-chart/values.yaml`
- Production overrides: `k8s/deploy/digitalocean/values.yaml`

Production values disable in-cluster Postgres/Redis/MinIO and use managed services plus externally managed secrets.

## Post-Deploy Verification

```bash
kubectl get deployment -n portcall
kubectl get pods -n portcall
kubectl get svc -n portcall
kubectl rollout status deployment/api -n portcall
kubectl rollout status deployment/dashboard -n portcall
```

## Notes

- Keep secrets in Kubernetes secrets, not in values files.
- Prefer pinned image tags over `latest`.
- If template behavior changes, validate with `helm template` before deploy.
