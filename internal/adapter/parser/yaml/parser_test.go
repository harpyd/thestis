package yaml_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	"github.com/harpyd/thestis/internal/domain/specification"
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
)

func TestSpecificationParserService_ParseSpecification(t *testing.T) {
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
			IsErr:       isComplexNoThesisHTTPOrAssertionError,
		},
		{
			Name:        "invalid_no_stories_specification",
			SpecPath:    invalidNoStoriesSpecPath,
			ShouldBeErr: true,
			IsErr:       isComplexNoStoriesError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			parser := yaml.NewSpecificationParserService()

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
	return specification.IsBuildSpecificationError(err) &&
		specification.IsBuildStoryError(err) &&
		specification.IsBuildScenarioError(err) &&
		specification.IsBuildThesisError(err) &&
		specification.IsBuildAssertionError(err) &&
		specification.IsNotAllowedAssertionMethodError(err)
}

func isComplexHTTPRequestMethodError(err error) bool {
	return specification.IsBuildSpecificationError(err) &&
		specification.IsBuildStoryError(err) &&
		specification.IsBuildScenarioError(err) &&
		specification.IsBuildThesisError(err) &&
		specification.IsBuildHTTPError(err) &&
		specification.IsBuildHTTPRequestError(err) &&
		specification.IsNotAllowedHTTPMethodError(err)
}

func isComplexKeywordError(err error) bool {
	return specification.IsBuildSpecificationError(err) &&
		specification.IsBuildStoryError(err) &&
		specification.IsBuildScenarioError(err) &&
		specification.IsBuildThesisError(err) &&
		specification.IsNotAllowedKeywordError(err)
}

func isComplexHTTPResponseContentTypeError(err error) bool {
	return specification.IsBuildSpecificationError(err) &&
		specification.IsBuildStoryError(err) &&
		specification.IsBuildScenarioError(err) &&
		specification.IsBuildThesisError(err) &&
		specification.IsBuildHTTPError(err) &&
		specification.IsBuildHTTPResponseError(err) &&
		specification.IsNotAllowedContentTypeError(err)
}

func isComplexNoThesisHTTPOrAssertionError(err error) bool {
	return specification.IsBuildSpecificationError(err) &&
		specification.IsBuildStoryError(err) &&
		specification.IsBuildScenarioError(err) &&
		specification.IsBuildThesisError(err) &&
		specification.IsNoThesisHTTPOrAssertionError(err)
}

func isComplexNoStoriesError(err error) bool {
	return specification.IsBuildSpecificationError(err) &&
		specification.IsNoSpecificationStoriesError(err)
}
