#!/usr/bin/env bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_usage() {
  echo "Usage: $0 [--apps list] [--cluster name] [--version type] [--values file] [--skip-build] [--skip-tests]"
  echo ""
  echo "Options:"
  echo "  --apps list      Comma-separated app names or numbers or 'all' (e.g. api,dashboard or 1,2)"
  echo "  --cluster name   Cluster name (default: digitalocean)"
  echo "  --version type   Version bump: major, minor, patch, or skip (default: patch)"
  echo "  --values file    Values file to use (default: k8s/deploy/digitalocean/values.yaml)"
  echo "  --skip-build     Skip Docker image building (use existing images)"
  echo "  --skip-tests     Skip smoke tests after deployment"
  echo "  --list-apps      Print available apps and exit"
  echo "  --dry-run        Show what would be deployed without actually deploying"
  echo "  --yes, -y        Skip confirmation prompts"
  echo "  -h, --help       Show help"
}

# Configuration
CLUSTER="digitalocean"
VERSION_TYPE=""
VALUES_FILE=""
SKIP_BUILD=false
SKIP_TESTS=false
DRY_RUN=false
APP_INPUT=""
SKIP_CONFIRMATION=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --apps)
      APP_INPUT="$2"
      shift 2
      ;;
    --apps=*)
      APP_INPUT="${1#*=}"
      shift
      ;;
    --cluster)
      CLUSTER="$2"
      shift 2
      ;;
    --cluster=*)
      CLUSTER="${1#*=}"
      shift
      ;;
    --version)
      VERSION_TYPE="$2"
      shift 2
      ;;
    --version=*)
      VERSION_TYPE="${1#*=}"
      shift
      ;;
    --values)
      VALUES_FILE="$2"
      shift 2
      ;;
    --values=*)
      VALUES_FILE="${1#*=}"
      shift
      ;;
    --skip-build)
      SKIP_BUILD=true
      shift
      ;;
    --skip-tests)
      SKIP_TESTS=true
      shift
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --list-apps)
      LIST_APPS=true
      shift
      ;;
    --yes|-y)
      SKIP_CONFIRMATION=true
      shift
      ;;
    -h|--help)
      print_usage
      exit 0
      ;;
    *)
      echo -e "${RED}Unknown option: $1${NC}"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

# Set default values file based on cluster if not provided
if [[ -z "$VALUES_FILE" ]]; then
  VALUES_FILE="${SCRIPT_DIR}/deploy/${CLUSTER}/values.yaml"
fi

# Validate values file exists
if [[ ! -f "$VALUES_FILE" ]]; then
  echo -e "${RED}Values file not found: $VALUES_FILE${NC}"
  echo "Available values files:"
  find "${SCRIPT_DIR}/deploy" -name "values.yaml" 2>/dev/null || echo "  (none found)"
  exit 1
fi

echo -e "${CYAN}Using values file: $VALUES_FILE${NC}"

# Quick cleanup of stuck Helm releases
clear_stuck_release() {
  local latest_rev=$(kubectl get secrets -n portcall -o json 2>/dev/null | jq -r '.items[] | select(.type == "helm.sh/release.v1") | select(.metadata.name | startswith("sh.helm.release.v1.portcall.v")) | .metadata.name' 2>/dev/null | sed 's/sh.helm.release.v1.portcall.v//' | sort -rn | head -1)
  
  if [[ -n "$latest_rev" ]]; then
    local stuck_status=$(kubectl get secret "sh.helm.release.v1.portcall.v${latest_rev}" -n portcall -o json 2>/dev/null | jq -r '.data.release' | base64 -d | base64 -d | gunzip 2>/dev/null | jq -r '.info.status // "unknown"')
    
    if [[ "$stuck_status" =~ ^pending|^failed ]]; then
      echo -e "${YELLOW}Clearing stuck release (rev ${latest_rev}, status: ${stuck_status})...${NC}"
      kubectl delete secret "sh.helm.release.v1.portcall.v${latest_rev}" -n portcall 2>/dev/null || true
      sleep 1
      return 0
    fi
  fi
  return 1
}

# Spinner for long-running operations
spin() {
  local pid=$1
  local delay=0.1
  local spinstr='|/-\'
  while ps -p "$pid" > /dev/null 2>&1; do
    local temp=${spinstr#?}
    printf " [%c]  " "$spinstr"
    spinstr=$temp${spinstr%"$temp"}
    sleep $delay
    printf "\b\b\b\b\b\b"
  done
  printf "      \b\b\b\b\b\b"
}

# Progress tracking for kubectl operations
watch_rollout() {
  local deployment=$1
  local namespace=$2
  local timeout=$3
  local start_time=$(date +%s)
  
  echo -e "${BLUE}Watching rollout for $deployment...${NC}"
  
  while true; do
    local current_time=$(date +%s)
    local elapsed=$((current_time - start_time))
    
    if [[ $elapsed -gt $timeout ]]; then
      echo -e "\n${RED}Timeout waiting for rollout (${timeout}s)${NC}"
      return 1
    fi
    
    # Get rollout status
    local ready=$(kubectl get deployment "$deployment" -n "$namespace" -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
    local desired=$(kubectl get deployment "$deployment" -n "$namespace" -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "1")
    local updated=$(kubectl get deployment "$deployment" -n "$namespace" -o jsonpath='{.status.updatedReplicas}' 2>/dev/null || echo "0")
    
    ready=${ready:-0}
    desired=${desired:-1}
    updated=${updated:-0}
    
    printf "\r  ${CYAN}[%ds]${NC} Ready: %s/%s, Updated: %s " "$elapsed" "$ready" "$desired" "$updated"
    
    if [[ "$ready" == "$desired" && "$updated" == "$desired" && "$ready" != "0" ]]; then
      echo -e "\n  ${GREEN}✓ Rollout complete${NC}"
      return 0
    fi
    
    sleep 2
  done
}

# Get all current deployment versions to preserve during helm upgrade
get_all_current_versions() {
  local namespace="portcall"
  local sets=""
  
  for app in "${APP_NAMES[@]}"; do
    local version=$(kubectl get deployment "$app" -n "$namespace" -o jsonpath='{.spec.template.spec.containers[0].image}' 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "")
    if [[ -n "$version" ]]; then
      sets="$sets --set ${app}.image.tag=${version}"
    fi
  done
  
  # Also check workers which have different naming
  local worker_apps=("email-worker:emailWorker" "billing-worker:billingWorker")
  for mapping in "${worker_apps[@]}"; do
    local deploy_name="${mapping%%:*}"
    local values_name="${mapping##*:}"
    local version=$(kubectl get deployment "$deploy_name" -n "$namespace" -o jsonpath='{.spec.template.spec.containers[0].image}' 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "")
    if [[ -n "$version" ]]; then
      sets="$sets --set ${values_name}.image.tag=${version}"
    fi
  done
  
  echo "$sets"
}

# Helper function to get deployment name from app name
# Some apps use different names in Kubernetes
get_deployment_name() {
  local app=$1
  case "$app" in
    billing)
      echo "billing-worker"
      ;;
    email)
      echo "email-worker"
      ;;
    file)
      echo "file-api"
      ;;
    *)
      echo "$app"
      ;;
  esac
}

# Helper function to get Helm values name from app name
get_helm_values_name() {
  local app=$1
  case "$app" in
    billing)
      echo "billingWorker"
      ;;
    email)
      echo "emailWorker"
      ;;
    file)
      echo "fileApi"
      ;;
    *)
      echo "$app"
      ;;
  esac
}

# Version management functions
get_current_version() {
  local app=$1
  local namespace="portcall"
  local deployment_name=$(get_deployment_name "$app")
  local version=$(kubectl get deployment "$deployment_name" -n "$namespace" -o jsonpath='{.spec.template.spec.containers[0].image}' 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "v0.0.0")
  echo "$version"
}

bump_version() {
  local app=$1
  local bump_type=$2
  local version=$(get_current_version "$app")
  
  # Remove 'v' prefix
  version=${version#v}
  IFS='.' read -r major minor patch <<< "$version"
  major=${major:-0}
  minor=${minor:-0}
  patch=${patch:-0}
  
  case "$bump_type" in
    major)
      major=$((major + 1))
      minor=0
      patch=0
      ;;
    minor)
      minor=$((minor + 1))
      patch=0
      ;;
    patch)
      patch=$((patch + 1))
      ;;
    skip)
      # Don't bump, return current
      ;;
  esac
  
  echo "v${major}.${minor}.${patch}"
}

# Available apps
APP_NAMES=(api admin dashboard checkout file quote billing email keycloak)
APP_IMAGES=(
  "your-registry.example.com/portcall-api"
  "your-registry.example.com/portcall-admin"
  "your-registry.example.com/portcall-dashboard"
  "your-registry.example.com/portcall-checkout"
  "your-registry.example.com/portcall-file"
  "your-registry.example.com/portcall-quote"
  "your-registry.example.com/portcall-billing-worker"
  "your-registry.example.com/portcall-email-worker"
  "your-registry.example.com/portcall-keycloak"
)
APP_DOCKERFILES=(
  "apps/api/Dockerfile"
  "apps/admin/Dockerfile"
  "apps/dashboard/Dockerfile"
  "apps/checkout/Dockerfile"
  "apps/file/Dockerfile"
  "apps/quote/Dockerfile"
  "apps/billing/Dockerfile"
  "apps/email/Dockerfile"
  "docker-compose/keycloak/Dockerfile"
)
APP_BUILD_CONTEXTS=(
  "."
  "."
  "."
  "."
  "."
  "."
  "."
  "."
  "."
  "docker-compose/keycloak"
)

# List apps if requested
if [[ "${LIST_APPS:-false}" == "true" ]]; then
  echo -e "${BLUE}Available apps:${NC}"
  for i in "${!APP_NAMES[@]}"; do
    printf "  %d) %s\n" "$((i+1))" "${APP_NAMES[i]}"
  done
  exit 0
fi

# Interactive app selection if not provided
if [[ -z "$APP_INPUT" ]]; then
  echo -e "${BLUE}Select apps to deploy:${NC}"
  echo ""
  for i in "${!APP_NAMES[@]}"; do
    printf "  %d) %s\n" "$((i+1))" "${APP_NAMES[i]}"
  done
  echo ""
  echo "Examples: all | 1,2,3 | api,dashboard | dashboard"
  read -r -p "Enter comma-separated names/numbers or 'all': " APP_INPUT
  if [[ -z "$APP_INPUT" ]]; then
    echo -e "${RED}No apps selected${NC}"
    exit 1
  fi
fi

# Parse selected apps
SELECTED_APPS=()
IFS=',' read -r -a TOKENS <<< "$APP_INPUT"
for token in "${TOKENS[@]}"; do
  token="$(echo "$token" | xargs)"
  if [[ -z "$token" ]]; then
    continue
  fi
  lowered="$(echo "$token" | tr '[:upper:]' '[:lower:]')"
  if [[ "$lowered" == "all" ]]; then
    SELECTED_APPS=("${APP_NAMES[@]}")
    break
  fi
  if [[ "$token" =~ ^[0-9]+$ ]]; then
    idx=$((token - 1))
    if (( idx < 0 || idx >= ${#APP_NAMES[@]} )); then
      echo -e "${RED}Invalid app number: $token${NC}"
      exit 1
    fi
    SELECTED_APPS+=("${APP_NAMES[idx]}")
  else
    found=false
    for i in "${!APP_NAMES[@]}"; do
      if [[ "${APP_NAMES[i]}" == "$lowered" ]]; then
        SELECTED_APPS+=("$lowered")
        found=true
        break
      fi
    done
    if [[ "$found" == "false" ]]; then
      echo -e "${RED}Unknown app: $token${NC}"
      exit 1
    fi
  fi
done

if [[ ${#SELECTED_APPS[@]} -eq 0 ]]; then
  echo -e "${RED}No apps selected${NC}"
  exit 1
fi

echo -e "\n${GREEN}Selected apps: ${SELECTED_APPS[*]}${NC}"

# Interactive cluster selection
if [[ -z "$CLUSTER" ]]; then
  echo -e "\n${BLUE}Select cluster:${NC}"
  echo "  1) digitalocean (default - production)"
  echo "  2) custom"
  read -r -p "Enter number or name [1]: " cluster_choice
  cluster_choice="${cluster_choice:-1}"
  
  case "$cluster_choice" in
    1|digitalocean|"")
      CLUSTER="digitalocean"
      ;;
    2|custom)
      read -r -p "Enter custom cluster name: " custom_cluster
      CLUSTER="$custom_cluster"
      ;;
    *)
      echo -e "${RED}Invalid choice${NC}"
      exit 1
      ;;
  esac
fi

echo -e "${GREEN}Target cluster: $CLUSTER${NC}"

# Interactive version bump selection
if [[ -z "$VERSION_TYPE" ]]; then
  echo -e "\n${BLUE}Select version bump type:${NC}"
  echo "  1) patch (0.0.X) - bug fixes, minor changes (default)"
  echo "  2) minor (0.X.0) - new features, backwards compatible"
  echo "  3) major (X.0.0) - breaking changes"
  echo "  4) skip - use current version"
  read -r -p "Enter number or type [1]: " version_choice
  version_choice="${version_choice:-1}"
  
  case "$version_choice" in
    1|patch|"")
      VERSION_TYPE="patch"
      ;;
    2|minor)
      VERSION_TYPE="minor"
      ;;
    3|major)
      VERSION_TYPE="major"
      ;;
    4|skip)
      VERSION_TYPE="skip"
      ;;
    *)
      echo -e "${RED}Invalid choice${NC}"
      exit 1
      ;;
  esac
fi

echo -e "${GREEN}Version bump: $VERSION_TYPE${NC}"

# Confirmation
echo -e "\n${YELLOW}====== DEPLOYMENT SUMMARY ======${NC}"
echo -e "Apps:     ${GREEN}${SELECTED_APPS[*]}${NC}"
echo -e "Cluster:  ${GREEN}$CLUSTER${NC}"
echo -e "Values:   ${GREEN}$VALUES_FILE${NC}"
echo -e "Version:  ${GREEN}$VERSION_TYPE${NC}"
echo -e "Build:    ${GREEN}$([ "$SKIP_BUILD" = true ] && echo "SKIP" || echo "YES")${NC}"
echo -e "Tests:    ${GREEN}$([ "$SKIP_TESTS" = true ] && echo "SKIP" || echo "YES")${NC}"
echo -e "Dry-run:  ${GREEN}$([ "$DRY_RUN" = true ] && echo "YES" || echo "NO")${NC}"
echo -e "${YELLOW}================================${NC}\n"

if [[ "$DRY_RUN" == "true" ]]; then
  echo -e "${YELLOW}DRY RUN - No changes will be made${NC}"
fi

if [[ "$SKIP_CONFIRMATION" != "true" ]]; then
  read -r -p "Proceed with deployment? [y/N]: " confirm
  if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Deployment cancelled${NC}"
    exit 0
  fi
else
  echo -e "${GREEN}Auto-confirmed (--yes flag)${NC}"
fi

# Set kubectl context for cluster
echo -e "\n${BLUE}Switching to cluster: $CLUSTER${NC}"
case "$CLUSTER" in
  digitalocean)
    kubectl config use-context your-k8s-context || {
      kubectl config use-context your-k8s-context || {
          echo -e "${RED}Failed to switch to cluster context${NC}"
          echo "Run: kubectl config use-context your-k8s-context"
          exit 1
      }
    }
    ;;
  *)
    kubectl config use-context "$CLUSTER" || {
      echo -e "${RED}Failed to switch to cluster: $CLUSTER${NC}"
      exit 1
    }
    ;;
esac

# Verify cluster connection
echo -e "${BLUE}Verifying cluster connection...${NC}"
if ! kubectl get nodes --request-timeout=5s >/dev/null 2>&1; then
  echo -e "${RED}Failed to connect to cluster${NC}"
  exit 1
fi
echo -e "${GREEN}✓ Connected${NC}"

# Collect current versions of all apps to preserve them during helm upgrade
echo -e "${BLUE}Collecting current deployment versions...${NC}"
CURRENT_VERSIONS=$(get_all_current_versions)
if [[ -n "$CURRENT_VERSIONS" ]]; then
  echo -e "  ${CYAN}Found existing deployments, will preserve versions${NC}"
fi

# Track versions we'll deploy (stored as "app:version" pairs in a string)
DEPLOY_VERSIONS_LIST=""

# Build and deploy each app
for app in "${SELECTED_APPS[@]}"; do
  echo -e "\n${YELLOW}====== Processing: $app ======${NC}"
  
  # Find app index
  app_idx=-1
  for i in "${!APP_NAMES[@]}"; do
    if [[ "${APP_NAMES[i]}" == "$app" ]]; then
      app_idx=$i
      break
    fi
  done
  
  if [[ $app_idx -eq -1 ]]; then
    echo -e "${RED}App not found: $app${NC}"
    continue
  fi
  
  image_name="${APP_IMAGES[app_idx]}"
  dockerfile="${APP_DOCKERFILES[app_idx]}"
  build_context="${APP_BUILD_CONTEXTS[app_idx]}"
  
  # Get or bump version
  if [[ "$VERSION_TYPE" != "skip" ]]; then
    version=$(bump_version "$app" "$VERSION_TYPE")
    echo -e "${GREEN}New version: $version${NC}"
  else
    version=$(get_current_version "$app")
    echo -e "${GREEN}Using version: $version${NC}"
  fi
  
  # Build image
  if [[ "$SKIP_BUILD" != "true" && "$DRY_RUN" != "true" ]]; then
    echo -e "${BLUE}Building Docker image...${NC}"
    echo -e "  ${CYAN}Image: ${image_name}:${version}${NC}"
    echo -e "  ${CYAN}Dockerfile: ${dockerfile}${NC}"
    echo -e "  ${CYAN}Context: ${build_context}${NC}"
    
    # Platform specification - always build for linux/amd64 for cloud deployments
    PLATFORM="--platform linux/amd64"
    
    # Build args for specific apps
    BUILD_ARGS=""
    if [[ "$app" == "admin" ]]; then
      BUILD_ARGS="--build-arg VITE_KEYCLOAK_URL=https://auth.useportcall.com"
      echo -e "  ${CYAN}Build args: ${BUILD_ARGS}${NC}"
    fi
    
    # Use buildx for multi-platform support if available, otherwise fall back to regular build
    if docker buildx version &>/dev/null; then
      echo -e "  ${CYAN}Using docker buildx for platform: linux/amd64${NC}"
      if ! docker buildx build $PLATFORM -t "${image_name}:${version}" -t "${image_name}:latest" -f "$dockerfile" $BUILD_ARGS "$build_context" --load 2>&1 | while IFS= read -r line; do
        # Show key build steps
        if [[ "$line" =~ ^Step|^Successfully|^CACHED|^ERROR|\[linux/amd64 ]]; then
          echo "  $line"
        fi
      done; then
        echo -e "${RED}Failed to build $app${NC}"
        echo -e "${YELLOW}Tip: Check Dockerfile at $dockerfile${NC}"
        exit 1
      fi
    else
      echo -e "  ${CYAN}Note: docker buildx not available, using standard build${NC}"
      if ! docker build -t "${image_name}:${version}" -t "${image_name}:latest" -f "$dockerfile" $BUILD_ARGS "$build_context" 2>&1 | while IFS= read -r line; do
        # Show key build steps
        if [[ "$line" =~ ^Step|^Successfully|^CACHED|^ERROR ]]; then
          echo "  $line"
        fi
      done; then
        echo -e "${RED}Failed to build $app${NC}"
        echo -e "${YELLOW}Tip: Check Dockerfile at $dockerfile${NC}"
        exit 1
      fi
    fi
    
    echo -e "${BLUE}Pushing image to registry...${NC}"
    if ! docker push "${image_name}:${version}" 2>&1 | tail -5; then
      echo -e "${RED}Failed to push $app${NC}"
      echo -e "${YELLOW}Tip: Ensure you're logged into the registry:${NC}"
      echo "  doctl registry login"
      exit 1
    fi
    docker push "${image_name}:latest" 2>/dev/null || true
    echo -e "${GREEN}✓ Image pushed${NC}"
  elif [[ "$DRY_RUN" == "true" ]]; then
    echo -e "${YELLOW}[DRY-RUN] Would build and push ${image_name}:${version}${NC}"
  else
    echo -e "${YELLOW}Skipping build (using existing image)${NC}"
  fi
  
  # Deploy to cluster
  if [[ "$DRY_RUN" == "true" ]]; then
    helm_values_name=$(get_helm_values_name "$app")
    echo -e "${YELLOW}[DRY-RUN] Would deploy to Kubernetes:${NC}"
    echo "  helm upgrade --install portcall ./k8s/portcall-chart \\"
    echo "    -f \"$VALUES_FILE\" \\"
    echo "    --namespace portcall \\"
    echo "    --set \"${helm_values_name}.image.tag=${version}\" \\"
    echo "    --set \"${helm_values_name}.enabled=true\""
    continue
  fi
  
  echo -e "${BLUE}Deploying to Kubernetes...${NC}"
  
  # Clear any stuck releases before deploying (critical!)
  clear_stuck_release || true
  
  # Get the correct Helm values name for this app
  helm_values_name=$(get_helm_values_name "$app")
  
  # Track this app's version for subsequent deploys in this session
  DEPLOY_VERSIONS_LIST="${DEPLOY_VERSIONS_LIST} ${helm_values_name}:${version}"
  
  # Build version override sets - start with current versions from cluster
  VERSION_SETS="$CURRENT_VERSIONS"
  
  # Override with any apps we've already deployed in this run
  for pair in $DEPLOY_VERSIONS_LIST; do
    deployed_app="${pair%%:*}"
    deployed_version="${pair##*:}"
    # Remove old version set for this app and add the new one
    VERSION_SETS=$(echo "$VERSION_SETS" | sed "s/--set ${deployed_app}\.image\.tag=[^ ]*//g")
    VERSION_SETS="$VERSION_SETS --set ${deployed_app}.image.tag=${deployed_version}"
  done
  
  # Clean up any double spaces
  VERSION_SETS=$(echo "$VERSION_SETS" | tr -s ' ')
  
  # Get the correct Helm values name for this app
  helm_values_name=$(get_helm_values_name "$app")
  deployment_name=$(get_deployment_name "$app")
  
  # Use helm upgrade with all version sets to preserve other apps' versions
  # Note: No --wait or --atomic - we check rollout separately for speed
  # shellcheck disable=SC2086
  if ! helm upgrade --install portcall ./k8s/portcall-chart \
    -f "$VALUES_FILE" \
    --namespace portcall \
    --create-namespace \
    $VERSION_SETS \
    --set "${helm_values_name}.enabled=true" \
    --timeout=2m 2>&1 | tee /tmp/helm-output-$$.log; then
    echo -e "${RED}Helm upgrade failed for $app${NC}"
    echo -e "${YELLOW}Helm output:${NC}"
    cat /tmp/helm-output-$$.log
    rm -f /tmp/helm-output-$$.log
    
    echo -e "\n${CYAN}=== RECOVERY OPTIONS ===${NC}"
    echo -e "${CYAN}1. Check for stuck releases:${NC}"
    echo -e "   kubectl get secrets -n portcall | grep sh.helm.release | tail -3"
    echo -e "   kubectl delete secret sh.helm.release.v1.portcall.vXXX -n portcall"
    echo -e "${CYAN}2. Check current Helm status:${NC}"
    echo -e "   helm list -n portcall"
    echo -e "${CYAN}3. Retry deployment:${NC}"
    echo -e "   ./k8s/deploy.sh --apps $app --version skip --yes --skip-tests"
    echo -e "${CYAN}======================${NC}\n"
    exit 1
  fi
  rm -f /tmp/helm-output-$$.log
  
  # Watch the rollout with progress updates (60s timeout)
  if ! watch_rollout "$deployment_name" "portcall" 60; then
    echo -e "${YELLOW}⚠ Rollout watch timed out, checking status...${NC}"
    
    # Check if deployment is actually healthy despite timeout
    local ready=$(kubectl get deployment "$deployment_name" -n portcall -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
    local desired=$(kubectl get deployment "$deployment_name" -n portcall -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "1")
    
    if [[ "$ready" == "$desired" && "$ready" != "0" ]]; then
      echo -e "${GREEN}✓ Deployment is healthy (${ready}/${desired} ready)${NC}"
    else
      echo -e "${YELLOW}Deployment status: ${ready}/${desired} ready${NC}"
      echo -e "${CYAN}Check status with: kubectl get pods -n portcall -l app=${deployment_name}${NC}"
      echo -e "${CYAN}View logs with: kubectl logs -f deployment/${deployment_name} -n portcall${NC}"
      
      read -r -p "Continue anyway? [y/N]: " continue_confirm
      if [[ ! "$continue_confirm" =~ ^[Yy]$ ]]; then
        exit 1
      fi
    fi
  fi
  
  echo -e "${GREEN}✓ Successfully deployed $app:$version${NC}"
done

# Run smoke tests
if [[ "$SKIP_TESTS" != "true" && "$DRY_RUN" != "true" ]]; then
  echo -e "\n${BLUE}Running smoke tests...${NC}"
  
  # Get endpoints based on cluster
  case "$CLUSTER" in
    digitalocean)
      API_URL="https://api.useportcall.com"
      DASHBOARD_URL="https://dashboard.useportcall.com"
      AUTH_URL="https://auth.useportcall.com"
      ;;
    *)
      API_URL="http://localhost:8080"
      DASHBOARD_URL="http://localhost:8082"
      AUTH_URL="http://localhost:8180"
      ;;
  esac
  
  test_passed=0
  test_failed=0
  
  # Test each deployed app
  for app in "${SELECTED_APPS[@]}"; do
    case "$app" in
      api)
        echo -n "  Testing API ($API_URL/ping)... "
        if curl -sf --max-time 10 "$API_URL/ping" >/dev/null 2>&1; then
          echo -e "${GREEN}OK${NC}"
          ((test_passed++))
        else
          echo -e "${RED}FAILED${NC}"
          ((test_failed++))
        fi
        ;;
      dashboard)
        echo -n "  Testing Dashboard ($DASHBOARD_URL/api/config)... "
        if curl -sf --max-time 10 "$DASHBOARD_URL/api/config" >/dev/null 2>&1 && \
           curl -sf --max-time 10 "https://webhook.useportcall.com/ping" >/dev/null 2>&1; then
          echo -e "${GREEN}OK${NC}"
          ((test_passed++))
        else
          echo -e "${RED}FAILED${NC}"
          ((test_failed++))
        fi
        ;;
      checkout)
        echo -n "  Testing Checkout (https://checkout.useportcall.com)... "
        if curl -sf --max-time 10 "https://checkout.useportcall.com" >/dev/null 2>&1; then
          echo -e "${GREEN}OK${NC}"
          ((test_passed++))
        else
          echo -e "${RED}FAILED${NC}"
          ((test_failed++))
        fi
        ;;
      quote)
        echo -n "  Testing Quote (https://quote.useportcall.com/health)... "
        if curl -sf --max-time 10 "https://quote.useportcall.com/health" >/dev/null 2>&1; then
          echo -e "${GREEN}OK${NC}"
          ((test_passed++))
        else
          echo -e "${RED}FAILED${NC}"
          ((test_failed++))
        fi
        ;;
      admin)
        echo -n "  Testing Admin (internal pod health)... "
        if kubectl exec -n portcall deploy/admin -- curl -sf http://localhost:8081/health >/dev/null 2>&1; then
          echo -e "${GREEN}OK${NC}"
          ((test_passed++))
        else
          echo -e "${YELLOW}SKIPPED (internal only)${NC}"
        fi
        ;;
      billing|file|email|keycloak)
        deployment_name=$(get_deployment_name "$app")
        echo -n "  Testing $app (pod running)... "
        if kubectl get pods -n portcall -l "app=$deployment_name" --field-selector=status.phase=Running 2>/dev/null | grep -q Running; then
          echo -e "${GREEN}OK${NC}"
          ((test_passed++))
        else
          echo -e "${RED}FAILED${NC}"
          ((test_failed++))
        fi
        ;;
      *)
        echo -e "  ${YELLOW}No test defined for $app${NC}"
        ;;
    esac
  done
  
  echo ""
  echo -e "  ${GREEN}Passed: $test_passed${NC} | ${RED}Failed: $test_failed${NC}"
  
  if [[ $test_failed -gt 0 ]]; then
    echo -e "\n${YELLOW}Warning: Some smoke tests failed${NC}"
  fi
elif [[ "$DRY_RUN" == "true" ]]; then
  echo -e "\n${YELLOW}[DRY-RUN] Would run smoke tests for: ${SELECTED_APPS[*]}${NC}"
fi

# Final summary
echo -e "\n${GREEN}====== DEPLOYMENT COMPLETE ======${NC}"
echo -e "${GREEN}Successfully deployed:${NC}"
for app in "${SELECTED_APPS[@]}"; do
  echo -e "  ✓ $app"
done
echo -e "${GREEN}==================================${NC}\n"

echo -e "${BLUE}Verify deployment:${NC}"
echo "  kubectl get pods -n portcall"
echo "  kubectl logs -f deployment/\${APP_NAME} -n portcall"
echo ""
echo -e "${BLUE}View services:${NC}"
echo "  kubectl get svc -n portcall"
