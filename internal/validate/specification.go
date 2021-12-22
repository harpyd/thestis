package validate

import (
	"log"
	"os"

	"github.com/gookit/color"

	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	"github.com/harpyd/thestis/internal/format"
)

const errorIndent = 2

func Specification(specPath string) {
	specFile, err := os.Open(specPath)
	if err != nil {
		log.Fatalf("%s: %s", specPath, err)
	}

	parser := yaml.NewSpecificationParserService()

	if _, err := parser.ParseSpecification(specFile); err != nil {
		fmtErr := format.SpecificationError(err, errorIndent)
		log.Fatalf("%s:\n%s", specFile.Name(), color.FgRed.Render(fmtErr))
	}
}
