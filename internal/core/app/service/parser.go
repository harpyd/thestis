package service

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type (
	SpecificationParser interface {
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

func WithSpecificationTestCampaignID(testCampaignID string) ParserOption {
	return func(b *specification.Builder) {
		b.WithTestCampaignID(testCampaignID)
	}
}

func WithSpecificationLoadedAt(loadedAt time.Time) ParserOption {
	return func(b *specification.Builder) {
		b.WithLoadedAt(loadedAt)
	}
}

type ParseError struct {
	err error
}

func WrapWithParseError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(&ParseError{
		err: err,
	})
}

func (e *ParseError) Unwrap() error {
	return e.err
}

func (e *ParseError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return fmt.Sprintf("parsing specification: %s", e.err)
}
