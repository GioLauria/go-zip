# Project Name

Short description.

> Template note: replace placeholders before publishing.

## Highlights

- Core value proposition.
- Key feature.
template-go

This repository is a reusable Go project template with a small example app and a library package.

Quick start

```bash
make build
./bin/app
```

Or run directly:

```bash
go run ./cmd/app
```

Project layout

- `cmd/app` — application entrypoint
- `pkg/` — reusable packages for import by apps
- `examples/` — usage examples
- `scripts/` — helper scripts (hooks installer)

CI & tooling

- GitHub Actions workflow: [/.github/workflows/go-ci.yml](.github/workflows/go-ci.yml)
- Linter config: `./.golangci.yml`
- Use `make lint` to run `golangci-lint` (install `golangci-lint` first)

Hooks

Run the hook installer for your platform:

- Unix/macOS:

	```bash
	sh scripts/install-hooks.sh
	```

- Windows (PowerShell):

	```powershell
	.\scripts\install-hooks.ps1
	```


CI and Git Hooks

- This project includes a GitHub Actions workflow at `.github/workflows/go-ci.yml` that runs `gofmt` check, `go vet`, `go test`, and builds the project on push and PR.
- To enable the included local git hooks run the installer for your platform:

	- Unix/macOS:

		```bash
		sh scripts/install-hooks.sh
		```

	- Windows (PowerShell):

		```powershell
		.\scripts\install-hooks.ps1
		```

	The installer sets `core.hooksPath` to the included `.githooks` directory. The pre-commit hook will check formatting and run `go vet` and `go test`.
"# go-zip" 
