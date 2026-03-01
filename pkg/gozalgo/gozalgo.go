package gozalgo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"runtime"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"
	lz4 "github.com/pierrec/lz4/v4"
)

// Simple block-based container format:
// 4 bytes magic: GOZ1
// then sequence of blocks:
// 1 byte algorithm id (0=zstd,1=lz4,2=brotli,3=snappy)
// 4 bytes BE compressed length
// 4 bytes BE uncompressed length
// compressed payload

const (
	Magic      = "GOZ2" // new format with version + CRC per block
	Version    = 1
	AlgoZstd   = 0
	AlgoLz4    = 1
	AlgoBrotli = 2
	AlgoSnappy = 3
)

// DefaultParallel controls the number of workers used by Compress when
// the caller does not set a parallelism value explicitly. Tests and callers
// may set this before invoking Compress to benchmark different settings.
var DefaultParallel int

// Compress reads from r and writes a GOZ archive to w.
// blockSize controls the size of chunks to try (e.g., 64*1024).
func Compress(r io.Reader, w io.Writer, blockSize int) error {
	if blockSize <= 0 {
		blockSize = 64 * 1024
	}
	parallel := DefaultParallel
	if parallel <= 0 {
		parallel = runtime.NumCPU()
	}

	type job struct {
		idx  int
		data []byte
	}
	type res struct {
		idx  int
		algo byte
		data []byte
		crc  uint32
		ulen uint32
		err  error
	}

	jobs := make(chan job, parallel*2)
	results := make(chan res, parallel*2)

	// workers
	var wg sync.WaitGroup
	wg.Add(parallel)
	for wkr := 0; wkr < parallel; wkr++ {
		go func() {
			defer wg.Done()
			zstdEnc, _ := zstd.NewWriter(nil)
			if zstdEnc != nil {
				defer zstdEnc.Close()
			}
			for j := range jobs {
				block := j.data
				var bestAlgo byte = AlgoSnappy
				var bestData []byte
				// zstd
				if zstdEnc != nil {
					zbuf := zstdEnc.EncodeAll(block, nil)
					bestAlgo = AlgoZstd
					bestData = append([]byte(nil), zbuf...)
				}
				// lz4
				lb := bytes.NewBuffer(nil)
				lw := lz4.NewWriter(lb)
				if _, err := lw.Write(block); err == nil {
					lw.Close()
					if bestData == nil || len(lb.Bytes()) < len(bestData) {
						bestAlgo = AlgoLz4
						bestData = lb.Bytes()
					}
				}
				// brotli
				bb := bytes.NewBuffer(nil)
				bw := brotli.NewWriterLevel(bb, brotli.BestCompression)
				if _, err := bw.Write(block); err == nil {
					bw.Close()
					if bestData == nil || len(bb.Bytes()) < len(bestData) {
						bestAlgo = AlgoBrotli
						bestData = bb.Bytes()
					}
				}
				// snappy
				sb := snappy.Encode(nil, block)
				if bestData == nil || len(sb) < len(bestData) {
					bestAlgo = AlgoSnappy
					bestData = sb
				}
				crc := crc32.ChecksumIEEE(block)
				results <- res{idx: j.idx, algo: bestAlgo, data: bestData, crc: crc, ulen: uint32(len(block)), err: nil}
			}
		}()
	}

	// start a writer goroutine that writes blocks in order as results arrive
	writeErrC := make(chan error, 1)
	go func() {
		// write header
		if _, err := w.Write([]byte(Magic)); err != nil {
			writeErrC <- err
			return
		}
		if _, err := w.Write([]byte{Version}); err != nil {
			writeErrC <- err
			return
		}

		pending := make(map[int]res)
		next := 0
		for {
			r, ok := <-results
			if !ok {
				break
			}
			if r.err != nil {
				writeErrC <- r.err
				return
			}
			pending[r.idx] = r
			for {
				rr, has := pending[next]
				if !has {
					break
				}
				// write block rr
				if _, err := w.Write([]byte{rr.algo}); err != nil {
					writeErrC <- err
					return
				}
				if err := binary.Write(w, binary.BigEndian, uint32(len(rr.data))); err != nil {
					writeErrC <- err
					return
				}
				if err := binary.Write(w, binary.BigEndian, rr.ulen); err != nil {
					writeErrC <- err
					return
				}
				if err := binary.Write(w, binary.BigEndian, rr.crc); err != nil {
					writeErrC <- err
					return
				}
				if _, err := w.Write(rr.data); err != nil {
					writeErrC <- err
					return
				}
				delete(pending, next)
				next++
			}
		}
		writeErrC <- nil
	}()

	// feed jobs
	idx := 0
	for {
		buf := make([]byte, blockSize)
		n, err := io.ReadFull(r, buf)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			close(jobs)
			wg.Wait()
			close(results)
			<-writeErrC
			return err
		}
		if n == 0 {
			break
		}
		block := append([]byte(nil), buf[:n]...)
		jobs <- job{idx: idx, data: block}
		idx++
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}
	close(jobs)
	wg.Wait()
	close(results)
	if err := <-writeErrC; err != nil {
		return err
	}
	return nil
}

// Decompress reads a GOZ archive from r and writes the decompressed data to w.
func Decompress(r io.Reader, w io.Writer) error {
	// read magic + version
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return err
	}
	magic := string(header)
	oldFormat := false
	if magic == Magic {
		// read version
		ver := make([]byte, 1)
		if _, err := io.ReadFull(r, ver); err != nil {
			return err
		}
		if ver[0] != Version {
			return errors.New("unsupported goz version")
		}
	} else if magic == "GOZ1" {
		// compatibility: old format has no version or CRC
		oldFormat = true
	} else {
		return errors.New("invalid goz magic")
	}
	for {
		h := make([]byte, 1)
		if _, err := io.ReadFull(r, h); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		algo := h[0]
		var clen uint32
		if err := binary.Read(r, binary.BigEndian, &clen); err != nil {
			return err
		}
		var ulen uint32
		if err := binary.Read(r, binary.BigEndian, &ulen); err != nil {
			return err
		}
		var crc uint32
		if !oldFormat {
			if err := binary.Read(r, binary.BigEndian, &crc); err != nil {
				return err
			}
		}
		data := make([]byte, clen)
		if _, err := io.ReadFull(r, data); err != nil {
			return err
		}
		var out []byte
		switch algo {
		case AlgoZstd:
			dec, err := zstd.NewReader(nil)
			if err != nil {
				return err
			}
			out, err = dec.DecodeAll(data, nil)
			dec.Close()
			if err != nil {
				return err
			}
		case AlgoLz4:
			lb := bytes.NewReader(data)
			lr := lz4.NewReader(lb)
			buf := make([]byte, ulen)
			if _, err := io.ReadFull(lr, buf); err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			out = buf
		case AlgoBrotli:
			br := brotli.NewReader(bytes.NewReader(data))
			buf := make([]byte, ulen)
			if _, err := io.ReadFull(br, buf); err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			out = buf
		case AlgoSnappy:
			dec, err := snappy.Decode(nil, data)
			if err != nil {
				return err
			}
			out = dec
		default:
			return errors.New("unknown algorithm in goz archive")
		}
		if !oldFormat {
			// verify CRC32 of uncompressed data
			if crc != crc32.ChecksumIEEE(out) {
				return errors.New("crc mismatch in goz archive block")
			}
		}
		if _, err := w.Write(out); err != nil {
			return err
		}
	}
}
