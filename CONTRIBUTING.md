# Contributing to go-zip

Thanks for your interest in contributing to **go-zip**. This document explains how to set up a development environment, run checks, and submit changes.

## Ways to contribute

- Report bugs and request features using the issue templates.
- Improve or expand documentation in `README.md` or `docs/`.
- Submit code changes via pull requests.
- Review open pull requests and provide feedback.

## Code of Conduct

Please follow the project's Code of Conduct: [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## Development setup

Prerequisites:

- Install Go (this repository declares Go 1.23 in `go.mod`).

Clone and prepare:

1. Fork the repository and clone your fork:

```bash
git clone git@github.com:<your-username>/go-zip.git
cd go-zip
git remote add upstream git@github.com:GioLauria/go-zip.git
```

2. Ensure dependencies are present and code is formatted:

```bash
go mod tidy
gofmt -w .
```

3. Run static checks and tests:

```bash
go vet ./...
go test ./...
```

## Local git hooks

This repository includes local hooks to run checks automatically. Install them with:

Unix / macOS:

```bash
sh scripts/install-hooks.sh
```

Windows (PowerShell):

```powershell
.\scripts\install-hooks.ps1
```

## Branching and commits

- Use short, descriptive branch names (e.g. `feat/add-compression-flag`, `fix/handle-empty-archive`).
- Keep commits focused and small.
- Use clear, imperative commit messages. We recommend Conventional Commits (e.g. `feat:`, `fix:`, `chore:`).

## Pull requests

- Open a pull request against `main` (or the repository's default branch).
- In the PR description explain the change, rationale, and how it was tested.
- Link related issues where applicable.
- Add or update tests for new behavior when relevant.

### Pull request checklist

- [ ] Tests added or updated where applicable
- [ ] Code formatted with `gofmt`
- [ ] CI passes (lint/test/build)
- [ ] Documentation updated (README or `docs/`)

## Reporting issues

Use the issue templates in `.github/ISSUE_TEMPLATE/` when filing bugs or feature requests. Provide steps to reproduce, expected vs actual behavior, and relevant logs or stack traces.

## Changelog and releases

This project follows Keep a Changelog (`CHANGELOG.md`). Add unreleased changes to the `Unreleased` section; maintainers will cut releases and update versioned headings.

## Additional notes

- If you're planning a larger design change, open an issue or discussion first to gather feedback.
- For maintainers: ensure CI, changelog, and release notes are updated when merging significant changes.
