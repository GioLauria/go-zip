package greetings

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Giova")
	want := "Hello, Giova!"
	if got != want {
		t.Fatalf("Hello(\"Giova\") = %q; want %q", got, want)
	}

	got = Hello("")
	want = "Hello, world!"
	if got != want {
		t.Fatalf("Hello(\"\") = %q; want %q", got, want)
	}
}
