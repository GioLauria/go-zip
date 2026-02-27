# Changelog Guidelines

This project follows Keep a Changelog: https://keepachangelog.com/en/1.0.0/

Process for contributors

- Add notable changes to `CHANGELOG.md` under the `## [Unreleased]` section.
- A local `post-commit` hook is provided to automatically append the last commit subject to the Unreleased section and amend the commit so the changelog is included. Hooks must be enabled locally via `scripts/install-hooks.sh` or `scripts/install-hooks.ps1`.

Maintainers

- When preparing a release, move entries from `Unreleased` into a new `## [x.y.z] - YYYY-MM-DD` section and update the changelog according to Keep a Changelog rules.
