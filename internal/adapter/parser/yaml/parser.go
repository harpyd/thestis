package yaml

import (
	"io"

	"gopkg.in/yaml.v2"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationParserService struct{}

func NewSpecificationParserService() SpecificationParserService {
	return SpecificationParserService{}
}

func (s *SpecificationParserService) ParseSpecification(reader io.Reader) (*specification.Specification, error) {
	decoder := yaml.NewDecoder(reader)

	var spec specificationSchema
	if err := decoder.Decode(&spec); err != nil {
		return nil, err
	}

	return specification.
		NewBuilder().
		WithAuthor(spec.Author).
		WithTitle(spec.Title).
		WithDescription(spec.Description).
		Build(), nil
}
