package mock

import (
	"io"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationParserService struct {
	withErr bool
}

func NewSpecificationParserService(withErr bool) *SpecificationParserService {
	return &SpecificationParserService{
		withErr: withErr,
	}
}

func (m *SpecificationParserService) ParseSpecification(
	_ io.Reader,
	opts ...app.ParserOption,
) (*specification.Specification, error) {
	builder := specification.NewBuilder()

	for _, opt := range opts {
		opt(builder)
	}

	if m.withErr {
		return builder.Build()
	}

	return builder.ErrlessBuild(), nil
}
