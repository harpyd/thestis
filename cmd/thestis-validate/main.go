package main

import (
	"flag"

	"github.com/harpyd/thestis/internal/validate"
)

const exampleSpecPath = "./examples/specification/horns-and-hooves-test.yml"

func main() {
	flag.Parse()

	specPath := flag.Arg(0)
	if specPath == "" {
		specPath = exampleSpecPath
	}

	validate.Specification(specPath)
}
