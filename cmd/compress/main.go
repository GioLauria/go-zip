package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	level := flag.Int("level", gzip.BestCompression, "gzip compression level (1-9). Higher = smaller, slower")
	out := flag.String("out", "", "output file path (optional). If empty, input + extension is used")
	ext := flag.String("ext", ".zpkg", "extension to append when -out not set (default: .zpkg)")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: compress [flags] <input-file>")
		flag.PrintDefaults()
		os.Exit(2)
	}

	inPath := flag.Arg(0)
	inFile, err := os.Open(inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "open input:", err)
		os.Exit(1)
	}
	defer inFile.Close()

	outPath := *out
	if outPath == "" {
		outPath = inPath + *ext
	}

	tmpOut, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create output:", err)
		os.Exit(1)
	}
	defer tmpOut.Close()

	gw, err := gzip.NewWriterLevel(tmpOut, *level)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create gzip writer:", err)
		os.Exit(1)
	}
	gw.Name = filepath.Base(inPath)

	before, _ := fileSize(inPath)
	copied, err := io.Copy(gw, inFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "compress copy:", err)
		os.Exit(1)
	}
	if err := gw.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "close gzip writer:", err)
		os.Exit(1)
	}
	after, _ := fileSize(outPath)

	fmt.Printf("compressed %s -> %s (%d bytes -> %d bytes)", inPath, outPath, before, after)
	if before > 0 {
		reduction := 100.0 * float64(before-after) / float64(before)
		fmt.Printf(" — %.1f%% reduction\n", reduction)
	} else {
		fmt.Println()
	}
	_ = copied
}

func fileSize(p string) (int64, error) {
	fi, err := os.Stat(p)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}
