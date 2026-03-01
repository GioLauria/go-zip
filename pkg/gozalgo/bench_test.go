package gozalgo

import (
	"bytes"
	"io"
	"runtime"
	"strconv"
	"testing"
)

// Benchmarks Compress performance for various block sizes and parallel workers.
// Run with:
//   go test ./pkg/gozalgo -bench . -benchmem

func BenchmarkCompress_Varied(b *testing.B) {
	// generate a sample input (4 MiB)
	data := make([]byte, 4*1024*1024)
	for i := range data {
		data[i] = byte(i)
	}

	blockSizes := []int{16 * 1024, 64 * 1024, 256 * 1024}
	parallels := []int{1, 2, 4, runtime.NumCPU()}

	for _, bs := range blockSizes {
		b.Run("bs="+itoa(bs), func(b *testing.B) {
			for _, p := range parallels {
				name := "par=" + itoa(p)
				b.Run(name, func(b *testing.B) {
					b.SetBytes(int64(len(data)))
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						r := bytes.NewReader(data)
						// discard output to avoid memory growth
						if err := Compress(r, io.Discard, bs, p); err != nil {
							b.Fatalf("compress failed: %v", err)
						}
					}
				})
			}
		})
	}
}

// small itoa helper for bench names
func itoa(v int) string {
	return strconv.Itoa(v)
}
