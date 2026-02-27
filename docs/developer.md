# Go Zip — Developer notes

Repository layout

- `cmd/goz` — CLI entrypoint and main
- `pkg/` — reusable packages and helpers
- `docs/` — project documentation (this folder)
- `scripts/` — helper scripts (install git hooks)

Hooks

- Local hooks are provided in `.githooks/`. Install them with:

```bash
sh scripts/install-hooks.sh
```

or on Windows PowerShell:

```powershell
.\scripts\install-hooks.ps1
```

The repository includes `pre-commit` checks and `post-commit` hooks that update the `CHANGELOG.md` automatically.

Extension enforcement

- The CLI enforces `.goz` as the package extension. `-out` will be normalized to end with `.goz`.
- Decompression rejects files that do not have the `.goz` extension.

Adding zstd backend (suggested)

To add zstd as an optional backend:

1. Add dependency: `github.com/klauspost/compress/zstd`.
2. Add a `-method` flag (e.g., `gzip` or `zstd`).
3. Implement codecs behind an interface `type Compressor interface { Compress(io.Reader, io.Writer) error; Decompress(io.Reader, io.Writer) error }`.
4. Make sure `.goz` remains the package extension (non-negotiable).

Testing and CI

- Unit tests: `go test ./...`
- Formatting: `gofmt -w .`
- CI: `.github/workflows` contains a CI workflow that runs the checks.
