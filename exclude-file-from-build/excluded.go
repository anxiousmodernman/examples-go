// +build excluded

package main

import "fmt"

// ExcludedFunction will not be included in the build. Note that the build tag
// must be followed by a blank line.
func ExcludedFunction() {
	fmt.Println("I won't be included unless you run `go build -tags=excluded`")
}
