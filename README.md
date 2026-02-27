# Go Zip

Go Zip (`goz`) is a small command-line tool and Go library for compressing single files into portable `.goz` packages.

IMPORTANT: the compressed package extension is always `.goz` — this is required and enforced by the CLI.

Features

- Compress a single file to a `.goz` package.
- Decompress a `.goz` package and restore the original filename and content.
- Cross-platform: builds for Windows, macOS, and Linux.

Build

Unix / macOS:

```bash
go build -o goz ./cmd/goz
```

Windows (PowerShell):

```powershell
go build -o goz.exe ./cmd/goz
```

Usage

Compress (writes `<input>.goz` by default):

```bash
./goz -C /path/to/file.txt
# Windows PowerShell:
.\goz.exe -C C:\path\to\file.txt
```

If you pass `-out`, the CLI will ensure the output filename ends with `.goz`.

Decompress (archive must have `.goz` extension):

```bash
./goz -D /path/to/file.txt.goz /path/to/outdir
# Windows PowerShell:
.\goz.exe -D C:\path\to\file.txt.goz C:\path\to\outdir
```

Flags

- `-C` : compress a single file (compress mode)
- `-D` : decompress archive into directory (decompress mode)
- `-out` : explicitly set output file (compress) or output directory (decompress)
- `-level` : gzip compression level (1-9)

Notes

- The CLI enforces the `.goz` extension for outputs and requires `.goz` for input archives when decompressing.
- Current backend: gzip. For higher compression, a zstd backend can be added.

Project layout

- `cmd/goz` — CLI
- `pkg/` — reusable packages and helpers
- `examples/` — example usage

Contributing

- Install local git hooks with `sh scripts/install-hooks.sh` (Unix/macOS) or `./scripts/install-hooks.ps1` (PowerShell).
- Run linters and tests locally: `gofmt -w .`, `go vet ./...`, `go test ./...`.

License

MIT — see `LICENSE`.
