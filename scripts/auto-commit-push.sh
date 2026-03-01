#!/bin/sh
set -e

# Helper: stage all changes, require documentation update for code changes,
# sign the commit and push.

msg="$1"
if [ -z "$msg" ]; then
  echo "Usage: $0 \"commit message\""
  exit 1
fi

git add -A

staged_files=$(git diff --cached --name-only)
code_changed=0
docs_present=0
for f in $staged_files; do
  case "$f" in
    docs/*|README.md|CHANGELOG.md|*.md)
      docs_present=1
      ;;
    *)
      code_changed=1
      ;;
  esac
done

if [ "$code_changed" -eq 1 ] && [ "$docs_present" -eq 0 ]; then
  echo "Error: code changes detected but no documentation updated. Aborting."
  echo "Please update README.md, docs/, or CHANGELOG.md and stage them, or run this script after updating docs."
  exit 1
fi

git commit -S -m "$msg"
git push
