#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "[deprecated] k8s/deploy.sh now proxies to dev-cli deploy"
echo "[recommended] go run ./tools/dev-cli deploy --cluster digitalocean"

if [[ ! -x ./dev-cli ]]; then
  if [[ -f ./tools/dev-cli/main.go ]]; then
    go run ./tools/dev-cli deploy "$@"
    exit $?
  fi
  echo "dev-cli binary/source not found"
  exit 1
fi

./dev-cli deploy "$@"
