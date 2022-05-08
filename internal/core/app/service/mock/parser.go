package mock

import (
	"io"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type SpecificationParserService struct {
	withErr bool
}

func NewSpecificationParserService(withErr bool) SpecificationParserService {
	return SpecificationParserService{withErr: withErr}
}

func (m SpecificationParserService) ParseSpecification(
	_ io.Reader,
	opts ...service.ParserOption,
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
