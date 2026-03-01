# Go Zip — Usage

This document describes the user-facing CLI for the `goz` tool.

Overview

- Executable: `goz` (Windows builds produce `goz.exe`).
- Package extension: `.goz` (required — the tool enforces this extension).

Build

Unix / macOS:

```bash
go build -o goz ./cmd/goz
```

Windows (PowerShell):

```powershell
go build -o goz.exe ./cmd/goz
```

Compress

Compress a single file. By default the output is written to `<input>.goz`.

```bash
./goz -C /path/to/file.txt
# Windows PowerShell
.\goz.exe -C C:\path\to\file.txt
```

Flags

- `-C` : compress mode
- `-D` : decompress mode
- `-out` : explicitly set output file (for compress) or output directory (for decompress)
- `-block-size` : block size in bytes for the `goz` algorithm (default: `65536`)
- `-parallel` : number of parallel workers for the `goz` algorithm (default: `runtime.NumCPU()`)

When using `-out` for compression the CLI will normalize the filename to end with `.goz`.

Decompress

Decompress a `.goz` archive into a directory. The archive must have a `.goz` extension.

```bash
./goz -D /path/to/file.txt.goz /path/to/outdir
# Windows PowerShell
.\goz.exe -D C:\path\to\file.txt.goz C:\path\to\outdir
```

Notes

- The `goz` tool stores compressed data in the project's custom GOZ container format. The archive preserves the original filename semantics; when you decompress a `.goz` archive the original base name is restored inside the specified output directory.
- The `.goz` extension is mandatory and enforced by the CLI.

See `docs/goz-format.md` for details on the GOZ container format.
