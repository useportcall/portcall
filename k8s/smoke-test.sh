#!/usr/bin/env bash
set -euo pipefail
NAMESPACE=portcall-dev

echo "Checking services in namespace: $NAMESPACE"

check_nodeport() {
  svc=$1
  path=$2
  port=$(kubectl get svc "$svc" -n "$NAMESPACE" -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null || true)
  if [ -z "$port" ]; then
    echo "  - $svc: no nodePort (skipping)"
    return
  fi
  echo -n "  - $svc (nodePort=$port) -> "
  if curl -sS --max-time 5 "http://localhost:$port$path" >/dev/null; then
    echo "OK"
  else
    echo "FAILED"
  fi
}

# Services to check (service, path)
services=(
  "api:/ping"
  "dashboard:/ping"
  "keycloak:/realms/dev"
  "checkout:/"
  "quote:/"
  "file-api:/ping"
  "minio-console:/"
)

for s in "${services[@]}"; do
  IFS=':' read -r svc path <<< "$s"
  check_nodeport "$svc" "$path"
done

# Check Postgres tables via exec
echo "Checking Postgres tables (via kubectl exec postgres-0)..."
if kubectl exec -n "$NAMESPACE" postgres-0 -- psql -U admin -d main_portcall_db -c "\dt" >/dev/null 2>&1; then
  echo "  - Postgres: OK (tables present)"
else
  echo "  - Postgres: FAILED"
fi

# Show job statuses
echo "Jobs:"
kubectl get jobs -n "$NAMESPACE" -o wide

echo "Done."
