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
- `-level` : gzip compression level (1-9), higher is slower and smaller
 - `-level` : gzip compression level (1-9), higher is slower and smaller
 - `-method` : compression method to use (`gzip` or `zstd`). `zstd` offers stronger compression; default is `gzip`.

When using `-out` for compression the CLI will normalize the filename to end with `.goz`.

Decompress

Decompress a `.goz` archive into a directory. The archive must have a `.goz` extension.

```bash
./goz -D /path/to/file.txt.goz /path/to/outdir
# Windows PowerShell
.\goz.exe -D C:\path\to\file.txt.goz C:\path\to\outdir
```

Notes

- The CLI stores the original filename in the gzip header; decompression restores that name inside the specified output directory.
- The `.goz` extension is mandatory and enforced by the CLI.
