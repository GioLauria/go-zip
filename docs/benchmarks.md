# gozalgo Benchmarks

This page summarizes local benchmark runs of `pkg/gozalgo` and provides recommended presets for common scenarios.

Environment used for the measurements

- CPU: AMD Ryzen 7 5825U (test machine)
- OS: Windows
- Command run:

```bash
go test ./pkg/gozalgo -bench . -benchmem
```

Representative results (throughput in MB/s)

- block-size=16 KiB, parallel=1  -> ~2.58 MB/s
- block-size=64 KiB, parallel=1  -> ~6.92 MB/s
- block-size=64 KiB, parallel=4  -> ~10.30 MB/s
- block-size=256 KiB, parallel=4 -> ~13.41 MB/s
- block-size=256 KiB, parallel=16 -> ~14.87 MB/s

Interpretation

- Larger block sizes (>= 256 KiB) reduce per-block overhead and generally increase throughput on multi-core systems.
- Increasing `-parallel` improves throughput up to a point; choose a value near your CPU core count for best results.
- Smaller block sizes (16 KiB) give lower latency for small writes but much lower throughput.

Recommended presets

- Default (balanced): `--block-size=65536 --parallel=runtime.NumCPU()`
- Throughput (high speed on multi-core): `--block-size=262144 --parallel=<num-cores or 4>`
- Low-latency (single-threaded, minimal memory): `--block-size=16384 --parallel=1`

How to run focused benchmarks

- Increase bench time to get stable numbers:

```bash
go test ./pkg/gozalgo -bench . -benchmem -benchtime=3s
```

- Run a single configuration (example):

```bash
GOBENCH_BLOCK=262144 GOBENCH_PAR=4 go test ./pkg/gozalgo -bench . -benchmem -run=^$ -benchtime=3s
```

(Replace the `GOBENCH_*` env vars with whatever harness you prefer; `pkg/gozalgo/bench_test.go` iterates several sizes and parallels by default.)

Notes

- Results vary by CPU, OS, and available memory. Use the supplied benchmark file to reproduce numbers on your target hardware and pick the preset that best fits your workload.
