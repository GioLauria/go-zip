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
echo >> RELEASE_NOTES.md

# Append built artifacts (if available) with sha256 checksums
echo "### Artifacts" >> RELEASE_NOTES.md
if [ -d artifacts ]; then
  # find files up to 2 levels (artifactName/file or direct file)
  find artifacts -type f -maxdepth 2 -print0 | while IFS= read -r -d '' file; do
    name=$(basename "${file}")
    if command -v sha256sum >/dev/null 2>&1; then
      sha=$(sha256sum "${file}" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
      sha=$(shasum -a 256 "${file}" | awk '{print $1}')
    else
      sha="(sha256 unavailable)"
    fi
    echo "- ${name} (sha256: ${sha})" >> RELEASE_NOTES.md
  done
else
  echo "No compiled artifacts found." >> RELEASE_NOTES.md
fi

echo >> RELEASE_NOTES.md
cat RELEASE_NOTES.md
