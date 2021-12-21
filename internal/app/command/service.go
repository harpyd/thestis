package command

import (
	"io"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type specificationParserService interface {
	ParseSpecification(reader io.Reader, opts ...app.ParserOption) (*specification.Specification, error)
}
