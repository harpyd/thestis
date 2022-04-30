package yaml_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/internal/core/infrastructure/parser/yaml"
)

const (
	fixturesPath                     = "./fixtures"
	validSpecPath                    = fixturesPath + "/valid-spec.yml"
	invalidAssertionMethodSpecPath   = fixturesPath + "/invalid-assertion-method-spec.yml"
	invalidHTTPMethodSpecPath        = fixturesPath + "/invalid-http-method-spec.yml"
	invalidNoKeywordSpecPath         = fixturesPath + "/invalid-no-keyword-spec.yml"
	invalidContentTypeSpecPath       = fixturesPath + "/invalid-content-type-spec.yml"
	invalidMixedErrorsSpecPath       = fixturesPath + "/invalid-mixed-errors-spec.yml"
	invalidNoHTTPOrAssertionSpecPath = fixturesPath + "/invalid-no-http-or-assertion-spec.yml"
	invalidNoStoriesSpecPath         = fixturesPath + "/invalid-no-stories-spec.yml"
	invalidNoScenariosSpecPath       = fixturesPath + "/invalid-no-scenarios-spec.yml"
	invalidNoThesesSpecPath          = fixturesPath + "/invalid-no-theses-spec.yml"
)

func TestParseSpecification(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	testCases := []struct {
		Name        string
		SpecPath    string
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "valid_specification",
			SpecPath:    validSpecPath,
			ShouldBeErr: false,
		},
		{
			Name:        "invalid_assertion_method_specification",
			SpecPath:    invalidAssertionMethodSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexAssertionMethodError,
		},
		{
			Name:        "invalid_http_method_specification",
			SpecPath:    invalidHTTPMethodSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexHTTPRequestMethodError,
		},
		{
			Name:        "invalid_no_keyword_specification",
			SpecPath:    invalidNoKeywordSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexKeywordError,
		},
		{
			Name:        "invalid_content_type_specification",
			SpecPath:    invalidContentTypeSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexHTTPResponseContentTypeError,
		},
		{
			Name:        "invalid_mixed_errors_specification",
			SpecPath:    invalidMixedErrorsSpecPath,
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return isComplexAssertionMethodError(err) &&
					isComplexKeywordError(err) &&
					isComplexHTTPRequestMethodError(err) &&
					isComplexHTTPResponseContentTypeError(err)
			},
		},
		{
			Name:        "invalid_no_http_or_assertion_specification",
			SpecPath:    invalidNoHTTPOrAssertionSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexUselessThesisError,
		},
		{
			Name:        "invalid_no_stories_specification",
			SpecPath:    invalidNoStoriesSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexNoStoriesError,
		},
		{
			Name:        "invalid_no_scenarios_specification",
			SpecPath:    invalidNoScenariosSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexNoScenariosError,
		},
		{
			Name:        "invalid_no_theses_specification",
			SpecPath:    invalidNoThesesSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexNoThesesError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			parser := yaml.NewSpecificationParser()

			specFile, err := os.Open(c.SpecPath)
			require.NoError(t, err)

			_, err = parser.ParseSpecification(specFile)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
		})
	}
}

func isComplexAssertionMethodError(err error) bool {
	var (
		berr *specification.BuildError
		nerr *specification.NotAllowedAssertionMethodError
	)

	return errors.As(err, &berr) && errors.As(err, &nerr)
}

func isComplexHTTPRequestMethodError(err error) bool {
	var (
		berr *specification.BuildError
		nerr *specification.NotAllowedHTTPMethodError
	)

	return errors.As(err, &berr) && errors.As(err, &nerr)
}

func isComplexKeywordError(err error) bool {
	var (
		berr *specification.BuildError
		nerr *specification.NotAllowedStageError
	)

	return errors.As(err, &berr) && errors.As(err, &nerr)
}

func isComplexHTTPResponseContentTypeError(err error) bool {
	var (
		berr *specification.BuildError
		nerr *specification.NotAllowedContentTypeError
	)

	return errors.As(err, &berr) && errors.As(err, &nerr)
}

func isComplexUselessThesisError(err error) bool {
	var berr *specification.BuildError

	return errors.As(err, &berr) &&
		errors.Is(err, specification.ErrUselessThesis)
}

func isComplexNoStoriesError(err error) bool {
	var berr *specification.BuildError

	return errors.As(err, &berr) &&
		errors.Is(err, specification.ErrNoSpecificationStories)
}

func isComplexNoScenariosError(err error) bool {
	var berr *specification.BuildError

	return errors.As(err, &berr) &&
		errors.Is(err, specification.ErrNoStoryScenarios)
}

func isComplexNoThesesError(err error) bool {
	var berr *specification.BuildError

	return errors.As(err, &berr) &&
		errors.Is(err, specification.ErrNoScenarioTheses)
}
