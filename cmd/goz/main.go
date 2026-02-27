package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	zstd "github.com/klauspost/compress/zstd"
)

// zstdReadCloser adapts zstd.Decoder to io.ReadCloser (Close returns error)
type zstdReadCloser struct{ *zstd.Decoder }

func (z zstdReadCloser) Close() error { z.Decoder.Close(); return nil }

func main() {
	compress := flag.Bool("C", false, "Compress a single file")
	decompress := flag.Bool("D", false, "Decompress a .goz archive to a folder")
	method := flag.String("method", "gzip", "compression method: gzip or zstd")
	level := flag.Int("level", gzip.BestCompression, "gzip compression level (1-9). Higher = smaller, slower")
	max := flag.Bool("max", false, "enable maximum compression (slower)")
	out := flag.String("out", "", "output file or directory (optional)")
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
		if err := compressFile(src, outPath, *level, *method, *max); err != nil {
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
		if err := decompressFile(archive, outDir, *method); err != nil {
			fmt.Fprintln(os.Stderr, "decompress error:", err)
			os.Exit(1)
		}
		fmt.Println("Decompressed ->", outDir)
		return
	}
}

func compressFile(src, dest string, level int, method string, max bool) error {
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
	switch strings.ToLower(method) {
	case "gzip":
		// enforce max level for gzip when requested
		if max {
			level = gzip.BestCompression
		}
		gw, err := gzip.NewWriterLevel(out, level)
		if err != nil {
			return err
		}
		gw.Name = filepath.Base(src)

		if _, err := io.Copy(gw, in); err != nil {
			gw.Close()
			return err
		}
		if err := gw.Close(); err != nil {
			return err
		}
		if fi, err := out.Stat(); err == nil {
			compressedSize = fi.Size()
		}
	case "zstd":
		// use strong zstd compression; if max set, prefer best compression and low concurrency
		if max {
			zw, err := zstd.NewWriter(out, zstd.WithEncoderLevel(zstd.SpeedBestCompression), zstd.WithEncoderConcurrency(1))
			if err != nil {
				return err
			}
			if _, err := io.Copy(zw, in); err != nil {
				zw.Close()
				return err
			}
			if err := zw.Close(); err != nil {
				return err
			}
			if fi, err := out.Stat(); err == nil {
				compressedSize = fi.Size()
			}
			break
		}
		// use strong zstd compression by default
		zw, err := zstd.NewWriter(out, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
		if err != nil {
			return err
		}
		if _, err := io.Copy(zw, in); err != nil {
			zw.Close()
			return err
		}
		if err := zw.Close(); err != nil {
			return err
		}
		if fi, err := out.Stat(); err == nil {
			compressedSize = fi.Size()
		}
	default:
		return fmt.Errorf("unknown compression method: %s", method)
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

func decompressFile(archive, outDir string, method string) error {
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

	var reader io.ReadCloser
	switch strings.ToLower(method) {
	case "gzip":
		gr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		reader = gr
	case "zstd":
		zr, err := zstd.NewReader(f)
		if err != nil {
			return err
		}
		reader = zstdReadCloser{zr}
	default:
		return fmt.Errorf("unknown compression method: %s", method)
	}
	defer reader.Close()

	// if gzip header contains name use it; for zstd header name not available
	name := ""
	if strings.ToLower(method) == "gzip" {
		if gr, ok := reader.(*gzip.Reader); ok {
			name = gr.Name
		}
	}
	if name == "" {
		name = filepath.Base(archive)
		// try to strip .goz
		if ext := filepath.Ext(name); ext == ".goz" {
			name = name[:len(name)-len(ext)]
		}
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

	if _, err := io.Copy(outFile, reader); err != nil {
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
