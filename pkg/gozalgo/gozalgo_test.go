package gozalgo

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestRoundTripSmall(t *testing.T) {
	data := []byte("the quick brown fox jumps over the lazy dog")
	var b bytes.Buffer
	if err := Compress(bytes.NewReader(data), &b, 32); err != nil {
		t.Fatalf("compress failed: %v", err)
	}
	// header magic + version
	if b.Len() < 5 {
		t.Fatalf("compressed too small")
	}
	if got := b.Bytes()[:4]; string(got) != Magic {
		t.Fatalf("bad magic: %q", got)
	}
	var out bytes.Buffer
	if err := Decompress(bytes.NewReader(b.Bytes()), &out); err != nil {
		t.Fatalf("decompress failed: %v", err)
	}
	if !bytes.Equal(out.Bytes(), data) {
		t.Fatalf("roundtrip mismatch")
	}
}

func TestRoundTripRandom(t *testing.T) {
	data := make([]byte, 256*1024)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		t.Fatalf("rand failed: %v", err)
	}
	var b bytes.Buffer
	if err := Compress(bytes.NewReader(data), &b, 64*1024); err != nil {
		t.Fatalf("compress failed: %v", err)
	}
	var out bytes.Buffer
	if err := Decompress(bytes.NewReader(b.Bytes()), &out); err != nil {
		t.Fatalf("decompress failed: %v", err)
	}
	if !bytes.Equal(out.Bytes(), data) {
		t.Fatalf("roundtrip mismatch for random data")
	}
}
