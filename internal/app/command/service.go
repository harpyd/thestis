package command

import (
	"io"

	"github.com/harpyd/thestis/internal/domain/specification"
)

// nolint
type specificationParserService interface {
	ParseSpecification(reader io.Reader) (*specification.Specification, error)
}
