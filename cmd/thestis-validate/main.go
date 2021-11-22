package main

import (
	"flag"
	"log"
	"os"

	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
)

func main() {
	flag.Parse()
	specFilename := flag.Arg(0)

	specFile, err := os.Open(specFilename)
	if err != nil {
		log.Fatalf("%s: %s", specFilename, err)
	}

	parser := yaml.NewSpecificationParserService()

	if _, err := parser.ParseSpecification(specFile); err != nil {
		log.Fatalf("%s: %s", specFile.Name(), err)
	}
}
