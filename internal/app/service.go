package app

import (
	"io"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationParserService interface {
	ParseSpecification(reader io.Reader, opts ...ParserOption) (*specification.Specification, error)
}

type MetricsService interface {
	IncRequestsCount(status, method, path string)
}
