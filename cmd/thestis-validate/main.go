package main

import (
	"flag"
	"log"
	"os"

	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
)

const exampleSpecification = "./examples/specification/horns-and-hooves-test.yml"

func main() {
	flag.Parse()

	specFilename := flag.Arg(0)
	if specFilename == "" {
		specFilename = exampleSpecification
	}

	specFile, err := os.Open(specFilename)
	if err != nil {
		log.Fatalf("%s: %s", specFilename, err)
	}

	parser := yaml.NewSpecificationParserService()

	if _, err := parser.ParseSpecification(specFile); err != nil {
		log.Fatalf("%s: %s", specFile.Name(), err)
	}
}
