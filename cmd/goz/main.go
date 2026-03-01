package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GioLauria/go-zip/pkg/gozalgo"
)

func main() {
	compress := flag.Bool("C", false, "Compress a single file")
	decompress := flag.Bool("D", false, "Decompress a .goz archive to a folder")
	out := flag.String("out", "", "output file or directory (optional)")
	blockSize := flag.Int("block-size", 64*1024, "block size in bytes for Goz.Algo")
	parallel := flag.Int("parallel", runtime.NumCPU(), "number of parallel workers for Goz.Algo (0 = auto)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  goz -C <input> [flags]        Compress a single file (writes <input>.goz by default)")
		fmt.Fprintln(os.Stderr, "  goz -D <archive.goz> <outdir>  Decompress a .goz archive into the specified directory")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  goz -C file.txt")
		fmt.Fprintln(os.Stderr, "  goz -C file.txt -out /tmp/archive.go z")
		fmt.Fprintln(os.Stderr, "  goz -D file.txt.goz /tmp/outdir")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Presets (recommended):")
		fmt.Fprintln(os.Stderr, "  Balanced:   --block-size=65536 --parallel=<num-cores>")
		fmt.Fprintln(os.Stderr, "  Throughput: --block-size=262144 --parallel=4")
		fmt.Fprintln(os.Stderr, "  Low-latency: --block-size=16384 --parallel=1")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Notes:")
		fmt.Fprintln(os.Stderr, "  - The CLI enforces the .goz extension for archives.")
		fmt.Fprintln(os.Stderr, "  - `goz` uses a per-block best-of strategy (zstd,lz4,brotli,snappy) internally; tune block size and parallelism for your workload.")
	}
	flag.Parse()

	if *compress == *decompress {
		fmt.Fprintln(os.Stderr, "Specify exactly one of -C (compress) or -D (decompress)")
		flag.Usage()
		os.Exit(2)
	}

	if *compress {
		if flag.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "Usage: goz -C /path/to/file [flags]")
			os.Exit(2)
		}
		src := flag.Arg(0)
		outPath := *out
		if outPath == "" {
			outPath = src + ".goz"
		}
		// always enforce .goz extension
		if !strings.HasSuffix(strings.ToLower(outPath), ".goz") {
			outPath = outPath + ".goz"
		}
		if err := compressFile(src, outPath, *blockSize, *parallel); err != nil {
			fmt.Fprintln(os.Stderr, "compress error:", err)
			os.Exit(1)
		}
		fmt.Println("Compressed ->", outPath)
		return
	}

	if *decompress {
		if flag.NArg() < 1 || flag.NArg() > 2 {
			fmt.Fprintln(os.Stderr, "Usage: goz -D /path/to/file.goz /path/to/outdir")
			os.Exit(2)
		}
		archive := flag.Arg(0)
		// require .goz extension
		if strings.ToLower(filepath.Ext(archive)) != ".goz" {
			fmt.Fprintln(os.Stderr, "decompress error: archive must have .goz extension")
			os.Exit(2)
		}
		outDir := *out
		if outDir == "" {
			if flag.NArg() == 2 {
				outDir = flag.Arg(1)
			} else {
				fmt.Fprintln(os.Stderr, "Must specify output directory with -out or as second arg for -D")
				os.Exit(2)
			}
		}
		if err := decompressFile(archive, outDir); err != nil {
			fmt.Fprintln(os.Stderr, "decompress error:", err)
			os.Exit(1)
		}
		fmt.Println("Decompressed ->", outDir)
		return
	}
}

func compressFile(src, dest string, blockSize int, parallel int) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	// get uncompressed size if available
	var uncompressedSize int64
	if fi, err := in.Stat(); err == nil {
		uncompressedSize = fi.Size()
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	var compressedSize int64
	// Always use the custom Goz.Algo
	// set package default parallelism and call Compress
	gozalgo.DefaultParallel = parallel
	if err := gozalgo.Compress(in, out, blockSize); err != nil {
		return err
	}
	if fi, err := out.Stat(); err == nil {
		compressedSize = fi.Size()
	}
	// print sizes and compression ratio (MB)
	fmt.Printf("Uncompressed: %.2f MB\n", float64(uncompressedSize)/(1024*1024))
	fmt.Printf("Compressed:   %.2f MB\n", float64(compressedSize)/(1024*1024))
	if uncompressedSize > 0 {
		saved := 100 * (1.0 - float64(compressedSize)/float64(uncompressedSize))
		fmt.Printf("Reduction:    %.2f%%\n", saved)
	}
	return nil
}

func decompressFile(archive, outDir string) error {
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()
	// compressed size (archive)
	var compressedSize int64
	if fi, err := f.Stat(); err == nil {
		compressedSize = fi.Size()
	}

	// Always treat input as GOZ custom format
	name := filepath.Base(archive)
	if ext := filepath.Ext(name); ext == ".goz" {
		name = name[:len(name)-len(ext)]
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	outPath := filepath.Join(outDir, name)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	if err := gozalgo.Decompress(f, outFile); err != nil {
		return err
	}
	// get uncompressed size
	var uncompressedSize int64
	if fi, err := outFile.Stat(); err == nil {
		uncompressedSize = fi.Size()
	}
	fmt.Printf("Compressed:   %.2f MB\n", float64(compressedSize)/(1024*1024))
	fmt.Printf("Uncompressed: %.2f MB\n", float64(uncompressedSize)/(1024*1024))
	if uncompressedSize > 0 {
		saved := 100 * (1.0 - float64(compressedSize)/float64(uncompressedSize))
		fmt.Printf("Reduction:    %.2f%%\n", saved)
	}
	return nil
}
