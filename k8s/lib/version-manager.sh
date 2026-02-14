#!/usr/bin/env bash

# Version Manager for Portcall Deployments
# Automatically bumps semantic versions based on flags

# Get current deployed version for a service
get_current_version() {
  local service=$1
  local namespace=$2
  
  # Try to get the current image tag from the deployment
  local current_version=$(kubectl get deployment "$service" -n "$namespace" -o jsonpath='{.spec.template.spec.containers[0].image}' 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "")
  
  if [ -z "$current_version" ]; then
    # Fallback: check values file
    current_version=$(grep -A 2 "^${service}:" "$(dirname "$0")/deploy/digitalocean/values.yaml" | grep "tag:" | awk '{print $2}' | tr -d '"' || echo "v0.0.0")
  fi
  
  echo "$current_version"
}

# Bump version based on type (patch, minor, major)
bump_version() {
  local version=$1
  local bump_type=${2:-patch}
  
  # Remove 'v' prefix if present
  version=${version#v}
  
  # Split version into components
  IFS='.' read -r major minor patch <<< "$version"
  
  # Default to 0.0.0 if parsing fails
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
    patch|*)
      patch=$((patch + 1))
      ;;
  esac
  
  echo "v${major}.${minor}.${patch}"
}

# Parse version flags and determine the version to use
determine_version() {
  local service=$1
  local namespace=$2
  local provided_version=$3
  local bump_type="patch"
  
  # Check for bump type flags in provided_version
  case "$provided_version" in
    --major|-M)
      bump_type="major"
      provided_version=""
      ;;
    --minor|-m)
      bump_type="minor"
      provided_version=""
      ;;
    --patch|-p|--auto|-a|"")
      bump_type="patch"
      provided_version=""
      ;;
  esac
  
  # If version is explicitly provided (and not a flag), use it
  if [ -n "$provided_version" ] && [[ "$provided_version" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # Ensure it has 'v' prefix
    [[ "$provided_version" =~ ^v ]] || provided_version="v${provided_version}"
    echo "$provided_version"
    return
  fi
  
  # Otherwise, auto-bump
  current_version=$(get_current_version "$service" "$namespace")
  new_version=$(bump_version "$current_version" "$bump_type")
  echo "$new_version"
}

# Main function to be called from deploy scripts
get_version() {
  local service=$1
  local namespace=$2
  local provided_version=$3
  
  determine_version "$service" "$namespace" "$provided_version"
}
