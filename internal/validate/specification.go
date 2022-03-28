package validate

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"

	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	"github.com/harpyd/thestis/internal/domain/specification"
)

const errorIndent = 2

func Specification(specPath string) {
	specFile, err := os.Open(specPath)
	if err != nil {
		log.Fatalf("%s: %s", specPath, err)
	}

	parser := yaml.NewSpecificationParser()

	if _, err := parser.ParseSpecification(specFile); err != nil {
		log.Fatalf("%s:\n%s", specFile.Name(), formatError(err, errorIndent))
	}
}

const (
	contextColor = color.FgLightMagenta
	errorColor   = color.FgRed
)

func formatError(err error, indent int) string {
	return formatErrorWithStartIndent(err, indent, indent)
}

func formatErrorWithStartIndent(err error, startIndent, indent int) string {
	if err == nil {
		return ""
	}

	var target *specification.BuildError

	if !errors.As(err, &target) {
		return errorColor.Render(err)
	}

	errs := target.Errors()
	if len(errs) == 0 {
		return ""
	}

	var b strings.Builder

	_, _ = fmt.Fprintf(&b, "%s:\n", contextColor.Render(target.Context()))

	for _, err := range errs {
		_, _ = fmt.Fprintf(
			&b,
			"%s%s",
			strings.Repeat(" ", indent),
			formatErrorWithStartIndent(err, startIndent, startIndent+indent),
		)
	}

	return b.String()
}
