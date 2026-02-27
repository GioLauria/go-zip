#!/bin/sh
# Install git hooks by setting core.hooksPath to .githooks
set -e
git config core.hooksPath .githooks
echo "Git hooks installed. (core.hooksPath = .githooks)"
