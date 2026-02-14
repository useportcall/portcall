#!/usr/bin/env bash

# Test the version manager
source "$(dirname "$0")/version-manager.sh"

echo "Version Manager Tests"
echo "===================="
echo ""

# Test bump_version function
echo "Testing version bumping:"
echo "  v1.0.0 patch bump: $(bump_version "v1.0.0" "patch")"
echo "  v1.0.0 minor bump: $(bump_version "v1.0.0" "minor")"
echo "  v1.0.0 major bump: $(bump_version "v1.0.0" "major")"
echo "  v1.2.3 patch bump: $(bump_version "v1.2.3" "patch")"
echo "  v1.2.3 minor bump: $(bump_version "v1.2.3" "minor")"
echo "  v1.2.3 major bump: $(bump_version "v1.2.3" "major")"
echo ""

# Test determine_version function
echo "Testing version determination:"
echo "  With explicit version v2.0.0: $(determine_version "dashboard" "portcall" "v2.0.0")"
echo "  With --patch flag: $(determine_version "dashboard" "portcall" "--patch")"
echo "  With --minor flag: $(determine_version "dashboard" "portcall" "--minor")"
echo "  With --major flag: $(determine_version "dashboard" "portcall" "--major")"
echo "  With no argument (default patch): $(determine_version "dashboard" "portcall" "")"
echo ""

# Test current version detection
echo "Current deployed versions:"
echo "  Dashboard: $(get_current_version "dashboard" "portcall")"
echo "  API: $(get_current_version "api" "portcall")"
echo ""
