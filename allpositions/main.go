package main

import (
	"fmt"
	"os"

	"github.com/donniet/tictacgo"
)

func main() {
	pos := tictacgo.AllValidPositionsEquiv()

	fmt.Fprintf(os.Stderr, "number of positions: %d\n", len(pos))

	// newline delimited output to stdout
	for _, p := range pos {
		fmt.Printf("%s\n", p)
	}
}
