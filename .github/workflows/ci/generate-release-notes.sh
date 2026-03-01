#!/usr/bin/env bash
set -euo pipefail

# Generate RELEASE_NOTES.md based on tags
TAG=${1:-${GITHUB_REF_NAME:-}}
PREV_TAG=${2:-${PREV_TAG:-}}
if [ -z "$TAG" ]; then
  echo "TAG not set; exiting"
  exit 1
fi
if [ -z "$PREV_TAG" ]; then
  RANGE=""
else
  RANGE="$PREV_TAG..$TAG"
fi
DATE=$(date +%Y-%m-%d)
echo "## [$TAG] - $DATE" > RELEASE_NOTES.md
echo >> RELEASE_NOTES.md
if [ -z "$RANGE" ]; then
  git log --pretty=format:"- %s (%an)" --no-merges >> RELEASE_NOTES.md || true
else
  git log "$RANGE" --pretty=format:"- %s (%an)" --no-merges >> RELEASE_NOTES.md || true
fi
echo >> RELEASE_NOTES.md
cat RELEASE_NOTES.md
