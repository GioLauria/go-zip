# GOZ Container Format

This document describes the GOZ container used by the `goz` tool (`.goz` files).

Overview

- File extension: `.goz`
- Primary implementation: `pkg/gozalgo` (GOZ2 format with backward compatibility for GOZ1)

GOZ2 layout (current)

- 4 bytes: ASCII magic `GOZ2`
- 1 byte: version (current value: `1`)
- Then a sequence of blocks. Each block has the following structure:
  - 1 byte: algorithm id
    - `0` = zstd
    - `1` = lz4
    - `2` = brotli
    - `3` = snappy
  - 4 bytes: compressed length (unsigned 32-bit, big-endian)
  - 4 bytes: uncompressed length (unsigned 32-bit, big-endian)
  - 4 bytes: CRC32 of the uncompressed payload (IEEE polynomial, big-endian)
  - N bytes: compressed payload (N == compressed length)

Notes

- Integer fields use big-endian encoding.
- CRC32 verifies the uncompressed payload when decompressing; a mismatch indicates data corruption.
- Blocks are independent and can be decompressed in order and streamed to the output.

GOZ1 compatibility

- GOZ1 uses the same per-block layout but omits the version byte and the CRC field. The decompressor accepts GOZ1 for backward compatibility.

Design rationale

- Per-block best-of: each block is compressed with multiple algorithms (zstd, lz4, brotli, snappy) and the smallest result is chosen to optimize overall archive size.
- Streaming and parallelism: `goz` splits input into blocks, compresses them in parallel workers, and writes them in order to the output stream without buffering the whole file in memory.

Extending the format

- New algorithm ids may be added in future versions. When adding new fields to the container header prefer a new version number and keep the reader backward/forward compatible.

See `pkg/gozalgo` for the implementation details and `cmd/goz` for CLI usage.
