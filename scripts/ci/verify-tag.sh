#!/usr/bin/env bash
set -euo pipefail

# Best-effort tag signature verification
TAG=${1:-${GITHUB_REF_NAME:-}}
if [ -z "$TAG" ]; then
  echo "No tag provided; exiting"
  exit 0
fi
echo "Attempting to verify tag signature for $TAG"
if ! command -v gpg >/dev/null 2>&1; then
  if command -v apt-get >/dev/null 2>&1; then
    sudo apt-get update && sudo apt-get install -y gnupg
  fi
fi
if git tag -v "$TAG" 2>&1 | tee /tmp/tag-verify.out; then
  echo "Tag $TAG: signature verified"
  exit 0
else
  echo "Warning: tag verification failed or public key missing. Proceeding anyway."
  cat /tmp/tag-verify.out || true
  exit 0
fi
