# Go Zip

Go Zip is a small, practical command-line tool and library for compressing single files with strong, easy-to-use defaults. The project provides a simple CLI at `cmd/goz` that produces compressed packages with the extension `.goz`.

Key points

- Purpose: fast, straightforward file compression for single files.
- CLI: `cmd/goz` — compress an input file to `<input>.goz` by default or decompress with `-D`.
- Default extension: `.goz` (customizable with `-out`).
- Notes: gzip is the current default implementation; zstd is recommended for higher compression and can be added later.

Quick start

Build the CLI:

```powershell
go build -o goz ./cmd/goz
```

Compress a file:

```powershell
./goz -C C:\path\to\file.txt
```

Decompress a `.goz` archive into a folder:

```powershell
./goz -D C:\path\to\file.txt.goz C:\path\to\outdir
```

Common flags:

- `-C` : compress a single file
- `-D` : decompress an archive to a directory
- `-out` : explicitly set output file or directory
- `-level` : gzip compression level (1-9)

What `goz` does

- Reads a single input file and writes a compressed package (single-file archive).
- By default writes `input + .goz` and stores the original filename in the gzip header.
- Prints original and compressed sizes and the percent reduction.

Why `.goz`?

`.goz` is a lightweight, project-specific extension that avoids colliding with common system extensions while being easy to identify as a Go Zip package. If/when zstd is adopted you may keep `.goz` or choose a zstd-specific marker such as `.gozst`.

Roadmap / Recommendations

- Add optional `zstd` backend (e.g. `github.com/klauspost/compress/zstd`) for higher compression ratios.
- Add multi-file packaging or archive+compress mode if multi-file distribution is needed.
# Go Zip

Go Zip is a small, practical command-line tool and library for compressing single files with strong, easy-to-use defaults. The project provides a simple CLI at `cmd/goz` that produces compressed packages with the extension `.goz`.

Key points

- Purpose: fast, straightforward file compression for single files.
- CLI: `cmd/goz` — compress an input file to `<input>.goz` by default or decompress with `-D`.
- Default extension: `.goz` (customizable with `-out` or `-ext`).
- Notes: gzip is the current default implementation; zstd is recommended for higher compression and can be added later.

Quick start

Build the CLI:

```powershell
go build -o goz.exe ./cmd/goz
```

Compress a file:

```powershell
./goz.exe -C C:\path\to\file.txt
```

Decompress a `.goz` archive into a folder:

```powershell
./goz.exe -D C:\path\to\file.txt.goz C:\path\to\outdir
```

Common flags:

- `-C` : compress a single file
- `-D` : decompress an archive to a directory
- `-out` : explicitly set output file or directory
- `-level` : gzip compression level (1-9)

What `goz` does

- Reads a single input file and writes a compressed package (single-file archive).
- By default writes `input + .goz` and stores the original filename in the gzip header.
- Prints original and compressed sizes and the percent reduction.

Why `.goz`?

`.goz` is a lightweight, project-specific extension that avoids colliding with common system extensions while being easy to identify as a Go Zip package. If/when zstd is adopted you may keep `.goz` or choose a zstd-specific marker such as `.gozst`.

Roadmap / Recommendations

- Add optional `zstd` backend (e.g. `github.com/klauspost/compress/zstd`) for higher compression ratios.
- Add multi-file packaging or archive+compress mode if multi-file distribution is needed.

Project layout

- `cmd/goz` — compressor/decompressor CLI
- `pkg/` — library code and helpers
- `examples/` — usage examples

Contributing

See the repository's tooling and CI configs for linting and tests. Use the hook installer in `scripts/` to enable local git hooks.

To enable the included local git hooks run the installer for your platform:

- Unix/macOS:

```bash
sh scripts/install-hooks.sh
```

- Windows (PowerShell):

```powershell
.\scripts\install-hooks.ps1
```
