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

Extensibility

- The project currently supports a single, custom container format (`.goz`) implemented by the `pkg/gozalgo` package.
- If you want to add an alternative compression backend in the future, follow these guidelines:
	1. Implement a small adapter that satisfies a compressor interface such as:

		 ```go
		 type Compressor interface {
				 Compress(io.Reader, io.Writer, int, int) error // (r, w, blockSize, parallel)
				 Decompress(io.Reader, io.Writer) error
		 }
		 ```

	2. Register or call that adapter from `cmd/goz` where appropriate.
	3. Keep `.goz` as the on-disk package extension; implement any new container metadata in a backwards-compatible way.

Note: the current CLI and documentation assume `goz` is the sole supported method.

Tag signing

- Releases (git tags) are expected to be GPG-signed. The local `pre-push` hook rejects pushing unsigned tags.
- Create an annotated signed tag locally with:

```bash
git tag -s vX.Y.Z -m "chore(release): vX.Y.Z"
git push origin vX.Y.Z
```

Ensure you have GPG configured for Git (see `git help gpg` / `git config user.signingkey`).

Testing and CI

- Unit tests: `go test ./...`
- Formatting: `gofmt -w .`
- CI: `.github/workflows` contains a CI workflow that runs the checks.
