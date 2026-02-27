package greetings

import "fmt"

// Hello returns a greeting for the provided name.
// If name is empty it returns a greeting for "world".
func Hello(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
