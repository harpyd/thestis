package specification

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrUnknownKeyword = errors.New("unknown keyword")

type noElemWithSlugError struct {
	elemName string
	slug     string
}

func NewNoStoryError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: "story",
		slug:     slug,
	})
}

func IsNoStoryError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == "story"
}

func NewNoScenarioError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: "scenario",
		slug:     slug,
	})
}

func IsNoScenarioError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == "scenario"
}

func NewNoThesisError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: "thesis",
		slug:     slug,
	})
}

func IsNoThesisError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == "thesis"
}

func (e noElemWithSlugError) Error() string {
	return fmt.Sprintf("no %s with slug %s", e.elemName, e.slug)
}
