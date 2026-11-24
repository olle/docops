// Command minifycss writes a minified copy of a CSS file.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/olle/docops/internal/cssmin"
)

func main() {
	input := flag.String("in", "", "CSS source file")
	output := flag.String("out", "", "minified output file")
	flag.Parse()

	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(2)
	}

	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	minified := cssmin.Minify(string(data))
	if err := os.WriteFile(*output, []byte(minified), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%s -> %s (%d -> %d bytes)\n", *input, *output, len(data), len(minified))
}
