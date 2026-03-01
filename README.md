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
- `-block-size` : block size in bytes for the `goz` algorithm (default: `65536`)
- `-parallel` : number of parallel workers for the `goz` algorithm (default: `runtime.NumCPU()`)

Recommended defaults

- Default `-block-size=65536` (64 KiB) is a reasonable balance between latency and throughput.
- For higher throughput on multi-core machines, try `-block-size=262144` (256 KiB) and `-parallel` equal to your CPU core count; benchmarks in `pkg/gozalgo/bench_test.go` show this combination typically gives the best MB/s on modern CPUs.

Notes

- The CLI enforces the `.goz` extension for outputs and requires `.goz` for input archives when decompressing.
- The `goz` tool uses the `pkg/gozalgo` implementation which performs per-block "best-of" compression: each block is compressed with multiple algorithms (zstd, lz4, brotli, snappy) and the smallest compressed result is stored. This favors compressed size; tune `-block-size`/`-parallel` for throughput as needed.

Project layout

- `cmd/goz` — CLI
- `pkg/` — reusable packages and helpers
- `examples/` — example usage

Contributing

- Install local git hooks with `sh scripts/install-hooks.sh` (Unix/macOS) or `./scripts/install-hooks.ps1` (PowerShell).
- Run linters and tests locally: `gofmt -w .`, `go vet ./...`, `go test ./...`.

License

MIT — see `LICENSE`.


Acknowledgements

- Compression libraries used internally by `pkg/gozalgo`:
	- `github.com/klauspost/compress` (zstd)
	- `github.com/pierrec/lz4/v4` (lz4)
	- `github.com/andybalholm/brotli` (brotli)
	- `github.com/golang/snappy` (snappy)


Contributing

If you'd like to contribute, see `CONTRIBUTING.md` for development setup, tests, and pull request guidelines. Install local hooks with `sh scripts/install-hooks.sh` or `./scripts/install-hooks.ps1` on Windows to enable pre-commit checks.

## GOZ container format

go-zip uses a simple container format for `.goz` archives. There are two versions in use:

- GOZ1: original format (legacy compatibility). Structure:
	- 4 bytes: ASCII `GOZ1` magic
	- repeated blocks:
		- 1 byte: algorithm id (0=zstd,1=lz4,2=brotli,3=snappy)
		- 4 bytes BE: compressed length
		- 4 bytes BE: uncompressed length
		- compressed payload

- GOZ2: current format with version and per-block CRC. Structure:
	- 4 bytes: ASCII `GOZ2` magic
	- 1 byte: version (currently `1`)
	- repeated blocks:
		- 1 byte: algorithm id (0=zstd,1=lz4,2=brotli,3=snappy)
		- 4 bytes BE: compressed length
		- 4 bytes BE: uncompressed length
		- 4 bytes BE: CRC32 of uncompressed block
		- compressed payload

Notes:

- The `gozalgo` implementation writes the GOZ2 format and verifies CRCs on read. The decompressor is compatible with GOZ1 archives and will accept and decode them.
- The default algorithm selection strategy for `Goz.Algo` is per-block "best-of": each block is compressed with zstd, lz4, brotli, and snappy, and the smallest compressed result is stored. This favors size over CPU and can be tuned for performance.

Benchmarks

You can measure `gozalgo` throughput and CPU usage using the included benchmarks. From the repository root run:

```bash
go test ./pkg/gozalgo -bench . -benchmem
```

To run focused benchmarks with custom `-benchtime`:

```bash
go test ./pkg/gozalgo -bench . -benchmem -benchtime=3s
```

The benchmark file `pkg/gozalgo/bench_test.go` exercises different block sizes and worker counts to help choose sensible defaults for `--block-size` and `--parallel`.

