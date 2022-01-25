package app

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"

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

func WithSpecificationTestCampaignID(performanceID string) ParserOption {
	return func(b *specification.Builder) {
		b.WithTestCampaignID(performanceID)
	}
}

func WithSpecificationLoadedAt(loadedAt time.Time) ParserOption {
	return func(b *specification.Builder) {
		b.WithLoadedAt(loadedAt)
	}
}

type parsingError struct {
	err error
}

func NewParsingError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(parsingError{
		err: err,
	})
}

func IsParsingError(err error) bool {
	var target parsingError

	return errors.As(err, &target)
}

func (e parsingError) Cause() error {
	return e.err
}

func (e parsingError) Unwrap() error {
	return e.err
}

func (e parsingError) Error() string {
	return fmt.Sprintf("parsing specification: %s", e.err)
}
