package specification_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuilder_WithAuthor(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithAuthor("author")

	spec, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, "author", spec.Author())
}

func TestBuilder_WithTitle(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithTitle("specification")

	spec, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, "specification", spec.Title())
}

func TestBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithDescription("description")

	spec, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, "description", spec.Description())
}

func TestBuilder_WithStory(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithStory("firstStory", func(b *specification.StoryBuilder) {
		b.WithDescription("this is a first story")
	})
	builder.WithStory("secondStory", func(b *specification.StoryBuilder) {
		b.WithDescription("this is a second story")
	})

	spec, err := builder.Build()
	require.NoError(t, err)

	expectedFirstStory, err := specification.NewStoryBuilder().
		WithDescription("this is a first story").
		Build("firstStory")
	require.NoError(t, err)

	actualFirstStory, ok := spec.Story("firstStory")
	require.True(t, ok)
	require.Equal(t, expectedFirstStory, actualFirstStory)

	expectedSecondStory, err := specification.NewStoryBuilder().
		WithDescription("this is a second story").
		Build("secondStory")
	require.NoError(t, err)

	actualSecondStory, ok := spec.Story("secondStory")
	require.True(t, ok)
	require.Equal(t, expectedSecondStory, actualSecondStory)
}

func TestStoryBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithDescription("description")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "description", story.Description())
}

func TestStoryBuilder_WithAsA(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithAsA("author")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "author", story.AsA())
}

func TestStoryBuilder_WithInOrderTo(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithInOrderTo("to do something")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "to do something", story.InOrderTo())
}

func TestStoryBuilder_WithWantTo(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithWantTo("do work")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "do work", story.WantTo())
}

func TestStoryBuilder_WithScenario(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithScenario("firstScenario", func(b *specification.ScenarioBuilder) {
		b.WithDescription("this is a first scenario")
	})
	builder.WithScenario("secondScenario", func(b *specification.ScenarioBuilder) {
		b.WithDescription("this is a second scenario")
	})

	story, err := builder.Build("someStory")

	require.NoError(t, err)

	expectedFirstScenario, err := specification.NewScenarioBuilder().
		WithDescription("this is a first scenario").
		Build("firstScenario")
	require.NoError(t, err)

	actualFirstScenario, ok := story.Scenario("firstScenario")
	require.True(t, ok)
	require.Equal(t, expectedFirstScenario, actualFirstScenario)

	expectedSecondScenario, err := specification.NewScenarioBuilder().
		WithDescription("this is a second scenario").
		Build("secondScenario")
	require.NoError(t, err)

	actualSecondScenario, ok := story.Scenario("secondScenario")
	require.True(t, ok)
	require.Equal(t, expectedSecondScenario, actualSecondScenario)
}

func TestScenarioBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()
	builder.WithDescription("description")

	scenario, err := builder.Build("someScenario")

	require.NoError(t, err)
	require.Equal(t, "description", scenario.Description())
}

func TestScenarioBuilder_WithThesis(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()
	builder.WithThesis("getBeer", func(b *specification.ThesisBuilder) {
		b.WithStatement("when", "get beer")
		b.WithHTTP(func(b *specification.HTTPBuilder) {
			b.WithRequest(func(b *specification.HTTPRequestBuilder) {
				b.WithMethod("GET")
				b.WithURL("https://api/v1/products")
			})
			b.WithResponse(func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedCodes([]int{200})
				b.WithAllowedContentType("application/json")
			})
		})
	})
	builder.WithThesis("checkBeer", func(b *specification.ThesisBuilder) {
		b.WithStatement("then", "check beer")
		b.WithAssertion(func(b *specification.AssertionBuilder) {
			b.WithMethod("JSONPATH")
			b.WithAssert("getSomeBody.response.body.product", "beer")
		})
	})

	scenario, err := builder.Build("someScenario")

	require.NoError(t, err)

	expectedGetBeerThesis, err := specification.NewThesisBuilder().
		WithStatement("when", "get beer").
		WithHTTP(func(b *specification.HTTPBuilder) {
			b.WithRequest(func(b *specification.HTTPRequestBuilder) {
				b.WithMethod("GET")
				b.WithURL("https://api/v1/products")
			})
			b.WithResponse(func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedCodes([]int{200})
				b.WithAllowedContentType("application/json")
			})
		}).
		Build("getBeer")
	require.NoError(t, err)

	actualGetBeerThesis, ok := scenario.Thesis("getBeer")
	require.True(t, ok)
	require.Equal(t, expectedGetBeerThesis, actualGetBeerThesis)

	expectedCheckBeerThesis, err := specification.NewThesisBuilder().
		WithStatement("then", "check beer").
		WithAssertion(func(b *specification.AssertionBuilder) {
			b.WithMethod("jsonpath")
			b.WithAssert("getSomeBody.response.body.product", "beer")
		}).
		Build("checkBeer")
	require.NoError(t, err)

	actualCheckBeerThesis, ok := scenario.Thesis("checkBeer")
	require.True(t, ok)
	require.Equal(t, expectedCheckBeerThesis, actualCheckBeerThesis)
}

func TestThesisBuilder_WithStatement(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Keyword     string
		Behavior    string
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "build_with_allowed_given_keyword",
			Keyword:     "given",
			Behavior:    "hooves delivered to the warehouse",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_when_keyword",
			Keyword:     "when",
			Behavior:    "selling hooves",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_then_keyword",
			Keyword:     "then",
			Behavior:    "check that hooves are sold",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_keyword",
			Keyword:     "zen",
			Behavior:    "zen du dust",
			ShouldBeErr: true,
			IsErr:       specification.IsNotAllowedKeywordError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewThesisBuilder()
			builder.WithStatement(c.Keyword, c.Behavior)

			thesis, err := builder.Build("sellHooves")

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.ToLower(c.Keyword), thesis.Statement().Keyword().String())
			require.Equal(t, c.Behavior, thesis.Statement().Behavior())
		})
	}
}

func TestThesisBuilder_WithAssertion(t *testing.T) {
	t.Parallel()

	builder := specification.NewThesisBuilder()
	builder.WithAssertion(func(b *specification.AssertionBuilder) {
		b.WithMethod("JSONPATH")
		b.WithAssert("getSomeBody.response.body.type", "product")
	})

	thesis, err := builder.Build("someThesis")

	require.NoError(t, err)
	expectedAssertion, err := specification.NewAssertionBuilder().
		WithMethod("JSONPATH").
		WithAssert("getSomeBody.response.body.type", "product").
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedAssertion, thesis.Assertion())
}

func TestThesisBuilder_WithHTTP(t *testing.T) {
	t.Parallel()

	builder := specification.NewThesisBuilder()
	builder.WithHTTP(func(b *specification.HTTPBuilder) {
		b.WithRequest(func(b *specification.HTTPRequestBuilder) {
			b.WithMethod("GET")
			b.WithURL("https://some-api/v1/endpoint")
		})
		b.WithResponse(func(b *specification.HTTPResponseBuilder) {
			b.WithAllowedCodes([]int{200})
			b.WithAllowedContentType("application/json")
		})
	})

	thesis, err := builder.Build("someThesis")

	require.NoError(t, err)
	expectedHTTP, err := specification.NewHTTPBuilder().
		WithRequest(func(b *specification.HTTPRequestBuilder) {
			b.WithMethod("GET")
			b.WithURL("https://some-api/v1/endpoint")
		}).
		WithResponse(func(b *specification.HTTPResponseBuilder) {
			b.WithAllowedCodes([]int{200})
			b.WithAllowedContentType("application/json")
		}).
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedHTTP, thesis.HTTP())
}

func TestAssertionBuilder_WithMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Method      string
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "build_with_allowed_jsonpath_assertion_method",
			Method:      "JSONPATH",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_assertion_method",
			Method:      "JAYZ",
			ShouldBeErr: true,
			IsErr:       specification.IsNotAllowedAssertionMethodError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewAssertionBuilder()
			builder.WithMethod(c.Method)

			assertion, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.ToLower(c.Method), assertion.Method().String())
		})
	}
}

func TestAssertionBuilder_WithAssert(t *testing.T) {
	t.Parallel()

	builder := specification.NewAssertionBuilder()
	builder.WithAssert("getSomeBody.response.body.type", "product")
	builder.WithAssert("getSomeBody.response.body.items..price", []int{2100, 1100})
	builder.WithAssert("getSomeBody.response.body.items..amount", []int{10, 33})

	assertion, err := builder.Build()

	require.NoError(t, err)

	asserts := assertion.Asserts()

	require.Equal(t, []string{
		"getSomeBody.response.body.type",
		"getSomeBody.response.body.items..price",
		"getSomeBody.response.body.items..amount",
	}, mapAssertsToActual(asserts))

	require.Equal(t, []interface{}{
		"product",
		[]int{2100, 1100},
		[]int{10, 33},
	}, mapAssertsToExpected(asserts))
}

func mapAssertsToActual(asserts []specification.Assert) []string {
	expected := make([]string, 0, len(asserts))
	for _, a := range asserts {
		expected = append(expected, a.Actual())
	}

	return expected
}

func mapAssertsToExpected(asserts []specification.Assert) []interface{} {
	actual := make([]interface{}, 0, len(asserts))
	for _, a := range asserts {
		actual = append(actual, a.Expected())
	}

	return actual
}

func TestHTTPBuilder_WithRequest(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPBuilder()
	builder.WithRequest(func(b *specification.HTTPRequestBuilder) {
		b.WithMethod("GET")
		b.WithURL("https://api.warehouse/v1/horns")
	})

	http, err := builder.Build()

	require.NoError(t, err)
	expectedRequest, err := specification.NewHTTPRequestBuilder().
		WithMethod("GET").
		WithURL("https://api.warehouse/v1/horns").
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedRequest, http.Request())
}

func TestHTTPBuilder_WithResponse(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPBuilder()
	builder.WithResponse(func(b *specification.HTTPResponseBuilder) {
		b.WithAllowedCodes([]int{200})
		b.WithAllowedContentType("application/json")
	})

	http, err := builder.Build()

	require.NoError(t, err)
	expectedResponse, err := specification.NewHTTPResponseBuilder().
		WithAllowedCodes([]int{200}).
		WithAllowedContentType("application/json").
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, http.Response())
}

func TestHTTPRequestBuilder_WithURL(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPRequestBuilder()
	builder.WithURL("https://api.warehouse/v1/hooves")

	request, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, "https://api.warehouse/v1/hooves", request.URL())
}

func TestHTTPRequestBuilder_WithMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Method      string
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "build_with_allowed_get_method",
			Method:      "GET",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_post_method",
			Method:      "POST",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_put_method",
			Method:      "PUT",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_patch_method",
			Method:      "PATCH",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_delete_method",
			Method:      "DELETE",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_options_method",
			Method:      "OPTIONS",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_trace_method",
			Method:      "TRACE",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_connect_method",
			Method:      "CONNECT",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_head_method",
			Method:      "HEAD",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_past_method",
			Method:      "PAST",
			ShouldBeErr: true,
			IsErr:       specification.IsNotAllowedHTTPMethodError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewHTTPRequestBuilder()
			builder.WithMethod(c.Method)

			request, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.ToUpper(c.Method), request.Method().String())
		})
	}
}

func TestHTTPRequestBuilder_WithContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ContentType string
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "build_with_allowed_content_type_application/json",
			ContentType: "application/json",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_content_type_application/xml",
			ContentType: "application/xml",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_content_type",
			ContentType: "content/type",
			ShouldBeErr: true,
			IsErr:       specification.IsNotAllowedContentTypeError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewHTTPRequestBuilder()
			builder.WithContentType(c.ContentType)

			request, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.ToLower(c.ContentType), request.ContentType().String())
		})
	}
}

func TestHTTPRequestBuilder_WithBody(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPRequestBuilder()
	builder.WithBody(map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": map[string]interface{}{
			"d": 1,
			"e": map[string]interface{}{
				"f": 3,
			},
		},
	})

	request, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": map[string]interface{}{
			"d": 1,
			"e": map[string]interface{}{
				"f": 3,
			},
		},
	}, request.Body())
}

func TestHTTPResponseBuilder_WithAllowedCodes(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPResponseBuilder()
	builder.WithAllowedCodes([]int{200, 404})

	request, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, []int{200, 404}, request.AllowedCodes())
}

func TestHTTPResponseBuilder_WithAllowedContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ContentType string
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "build_with_allowed_content_type_application/json",
			ContentType: "application/json",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_content_type_application/xml",
			ContentType: "application/xml",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_content_type",
			ContentType: "some/content",
			ShouldBeErr: true,
			IsErr:       specification.IsNotAllowedContentTypeError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewHTTPResponseBuilder()
			builder.WithAllowedContentType(c.ContentType)

			request, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.ToLower(c.ContentType), request.AllowedContentType().String())
		})
	}
}
