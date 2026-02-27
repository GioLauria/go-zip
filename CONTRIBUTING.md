# Contributing

Thanks for your interest in contributing!

> Note: please review project name, language, and tooling placeholders if present before publishing.

## Ways to contribute

- Report bugs and request features
- Improve documentation
- Submit code changes
- Review pull requests

## Code of Conduct

Please follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Development setup

1. Fork and clone the repo.
2. Install dependencies (TBD).
3. Run tests (TBD).

## Branching and commits

- Use short, descriptive branch names.
- Keep commits focused and small.
- Use clear commit messages (TBD).

## Pull requests

- Describe the change and the motivation.
- Add or update tests where appropriate.
# Contributing

Thank you for helping improve Go Zip — contributions are welcome. This document explains how to set up a development environment, run checks, and submit changes.

Getting started

1. Fork this repository and create a feature branch: `git checkout -b feat/my-change`.
2. Ensure you have Go 1.20+ installed.
3. From the repository root run:

```bash
go mod tidy
gofmt -w .
go vet ./...
go test ./...
```

Local git hooks

This repo includes local hooks in `.githooks/`. Install them to run checks automatically:

Unix/macOS:

```bash
sh scripts/install-hooks.sh
```

Windows (PowerShell):

```powershell
.\scripts\install-hooks.ps1
```

Coding standards

- Run `gofmt -w .` before committing.
- Keep functions small and tests focused.
- Add unit tests for new behavior and run `go test ./...`.

Commit messages and PRs

- Use clear, imperative commit messages (e.g. `feat: add zstd method flag`).
- Open a pull request describing the change, rationale, and testing performed.
- Link related issues in the PR description.

Pull request checklist

- [ ] Tests added or updated where applicable
- [ ] Code formatted with `gofmt`
- [ ] CI passes (lint/test/build)
- [ ] Documentation updated (README or `docs/`)

Reporting issues

- Use the issue templates in `.github/ISSUE_TEMPLATE/` when filing bugs or feature requests.

Changelog and releases

- The project follows Keep a Changelog (`CHANGELOG.md`). Use the `Unreleased` section for ongoing changes; maintainers will cut releases and move entries into versioned headings.

Code of Conduct

- Please follow the project's Code of Conduct (see `CODE_OF_CONDUCT.md`).
