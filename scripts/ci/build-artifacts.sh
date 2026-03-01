#!/usr/bin/env bash
set -euo pipefail
#!/usr/bin/env bash
set -euo pipefail

# Build artifacts for the given OS matrix (called from GitHub Actions)
# Usage: build-artifacts.sh <matrix-os>
OS_MATRIX="$1"
mkdir -p artifacts
if [ "${OS_MATRIX}" = "windows-latest" ]; then
  GOOS=windows
  OUT=goz.exe
else
  GOOS=linux
  OUT=goz
fi
echo "Building for GOOS=$GOOS -> artifacts/${OUT}"
GOOS=$GOOS go build -ldflags "-s -w" -o artifacts/${OUT} ./cmd/goz

echo "Built: artifacts/${OUT}"
