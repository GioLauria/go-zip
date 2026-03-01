package main

import (
	"fmt"
	"os"
)

func main() {
	b, _ := os.ReadFile("pkg/gozalgo/gozalgo_test.go")
	lines := bytesSplit(b, '\n')
	for i, l := range lines {
		if i == 12 || i == 37 {
			fmt.Printf("Line %d: %s\n", i+1, string(l))
			for j, r := range []rune(string(l)) {
				fmt.Printf("%d:%q(%U) ", j, r, r)
			}
			fmt.Println()
		}
	}
}

func bytesSplit(b []byte, sep byte) [][]byte {
	var out [][]byte
	cur := []byte{}
	for _, c := range b {
		if c == sep {
			out = append(out, cur)
			cur = []byte{}
			continue
		}
		cur = append(cur, c)
	}
	out = append(out, cur)
	return out
}
