#!/usr/bin/env bash
set -euo pipefail

# Build artifacts for the given OS matrix (called from GitHub Actions)
# Usage: build-artifacts.sh <matrix-os> <tag>
OS_MATRIX="$1"
TAG="${2:-unknown}"
mkdir -p artifacts

if [ "${OS_MATRIX}" = "windows-latest" ]; then
  GOOS=windows
  SUFFIX=".exe"
else
  GOOS=linux
  SUFFIX=""
fi

# Determine GOARCH from the environment (setup-go provides `go`)
GOARCH=$(go env GOARCH || echo "amd64")

# Artifact name includes tag, OS and arch to make releases clear
OUT=goz-${TAG}-${GOOS}-${GOARCH}${SUFFIX}
echo "Building for GOOS=$GOOS GOARCH=$GOARCH -> artifacts/${OUT}"
GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o artifacts/${OUT} ./cmd/goz

echo "Built: artifacts/${OUT}"
