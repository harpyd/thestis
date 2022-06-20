package validate

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"

	"github.com/harpyd/thestis/internal/core/adapter/driven/parser/yaml"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

const errorIndent = "  "

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
	contextColor = color.FgGreen
	errorColor   = color.FgRed
)

func formatError(err error, indent string) string {
	return formatErrorWithStartIndent(err, indent, indent)
}

func formatErrorWithStartIndent(err error, startIndent, indent string) string {
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

	_, _ = fmt.Fprintf(&b, "%s:", contextColor.Render(prefix(target)))

	for _, err := range errs {
		_, _ = fmt.Fprintf(
			&b,
			"\n%s%s",
			indent,
			formatErrorWithStartIndent(err, startIndent, startIndent+indent),
		)
	}

	return b.String()
}

func prefix(err *specification.BuildError) string {
	if slug, ok := err.SlugContext(); ok {
		return slug.Partial()
	}

	msg, _ := err.StringContext()

	return msg
}
