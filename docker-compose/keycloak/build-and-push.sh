#!/bin/bash
set -e

# Build custom Keycloak image with Portcall theme
cd "$(dirname "$0")"

IMAGE_NAME="your-registry.example.com/portcall-keycloak"
TAG="v1.0.0"

echo "Building Keycloak image with Portcall theme for linux/amd64..."
docker buildx build --platform linux/amd64 -t ${IMAGE_NAME}:${TAG} . --push

echo "Done! Image: ${IMAGE_NAME}:${TAG}"
