package app

import (
	"io"
	"time"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	SpecificationParserService interface {
		ParseSpecification(reader io.Reader, opts ...ParserOption) (*specification.Specification, error)
	}

	ParserOption func(b *specification.Builder)
)

func WithSpecificationID(specID string) ParserOption {
	return func(b *specification.Builder) {
		b.WithID(specID)
	}
}

func WithSpecificationOwnerID(ownerID string) ParserOption {
	return func(b *specification.Builder) {
		b.WithOwnerID(ownerID)
	}
}

func WithSpecificationLoadedAt(loadedAt time.Time) ParserOption {
	return func(b *specification.Builder) {
		b.WithLoadedAt(loadedAt)
	}
}

type MetricsService interface {
	IncRequestsCount(status, method, path string)
}
