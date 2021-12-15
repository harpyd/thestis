package mock

import (
	"io"

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
	specID string,
	_ io.Reader,
) (*specification.Specification, error) {
	builder := specification.NewBuilder().WithID(specID)

	if m.withErr {
		return builder.Build()
	}

	return builder.ErrlessBuild(), nil
}
