package mock

import (
	"io"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/domain/specification"
)

type SpecificationParserService struct {
	withErr bool
}

func NewSpecificationParserService(withErr bool) SpecificationParserService {
	return SpecificationParserService{withErr: withErr}
}

func (m SpecificationParserService) ParseSpecification(
	_ io.Reader,
	opts ...app.ParserOption,
) (*specification.Specification, error) {
	var b specification.Builder

	for _, opt := range opts {
		opt(&b)
	}

	if m.withErr {
		return b.Build()
	}

	return b.ErrlessBuild(), nil
}
