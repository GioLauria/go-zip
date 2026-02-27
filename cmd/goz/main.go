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
    compress := flag.Bool("C", false, "Compress a single file")
    decompress := flag.Bool("D", false, "Decompress a .goz archive to a folder")
    level := flag.Int("level", gzip.BestCompression, "gzip compression level (1-9). Higher = smaller, slower")
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
        if err := compressFile(src, outPath, *level); err != nil {
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

func compressFile(src, dest string, level int) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
        return err
    }
    out, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer out.Close()

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
    return nil
}

func decompressFile(archive, outDir string) error {
    f, err := os.Open(archive)
    if err != nil {
        return err
    }
    defer f.Close()

    gr, err := gzip.NewReader(f)
    if err != nil {
        return err
    }
    defer gr.Close()

    name := gr.Name
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

    if _, err := io.Copy(outFile, gr); err != nil {
        return err
    }
    return nil
}
