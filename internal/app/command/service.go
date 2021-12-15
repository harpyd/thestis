package command

import (
	"io"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type specificationParserService interface {
	ParseSpecification(specID string, reader io.Reader) (*specification.Specification, error)
}
