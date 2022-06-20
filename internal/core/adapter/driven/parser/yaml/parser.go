package yaml

import (
	"io"

	"gopkg.in/yaml.v3"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type SpecificationParser struct{}

func NewSpecificationParser() SpecificationParser {
	return SpecificationParser{}
}

func (s SpecificationParser) ParseSpecification(
	reader io.Reader,
	opts ...service.ParserOption,
) (*specification.Specification, error) {
	decoder := yaml.NewDecoder(reader)

	var spec specificationSchema
	if err := decoder.Decode(&spec); err != nil {
		return nil, service.WrapWithParseError(err)
	}

	return build(spec, opts)
}

func build(spec specificationSchema, opts []service.ParserOption) (*specification.Specification, error) {
	var b specification.Builder

	b.
		WithAuthor(spec.Author).
		WithTitle(spec.Title).
		WithDescription(spec.Description)

	for _, opt := range opts {
		opt(&b)
	}

	for slug, story := range spec.Stories {
		b.WithStory(slug, buildStory(story))
	}

	return b.Build()
}

func buildStory(story storySchema) func(builder *specification.StoryBuilder) {
	return func(builder *specification.StoryBuilder) {
		builder.
			WithDescription(story.Description).
			WithAsA(story.AsA).
			WithInOrderTo(story.InOrderTo).
			WithWantTo(story.WantTo)

		for slug, scenario := range story.Scenarios {
			builder.WithScenario(slug, buildScenario(scenario))
		}
	}
}

func buildScenario(scenario scenarioSchema) func(builder *specification.ScenarioBuilder) {
	return func(builder *specification.ScenarioBuilder) {
		builder.WithDescription(scenario.Description)

		for slug, thesis := range scenario.Theses {
			builder.WithThesis(slug, buildThesis(thesis))
		}
	}
}

func buildThesis(thesis thesisSchema) func(builder *specification.ThesisBuilder) {
	return func(builder *specification.ThesisBuilder) {
		switch {
		case len(thesis.Given) > 0:
			builder.WithStatement(specification.Given, thesis.Given)
		case len(thesis.When) > 0:
			builder.WithStatement(specification.When, thesis.When)
		case len(thesis.Then) > 0:
			builder.WithStatement(specification.Then, thesis.Then)
		default:
			builder.WithStatement("", "")
		}

		builder.
			WithAssertion(buildAssertion(thesis.Assertion)).
			WithHTTP(buildHTTP(thesis.HTTP))

		for _, after := range thesis.After {
			builder.WithDependency(after)
		}
	}
}

func buildAssertion(assertion assertionSchema) func(builder *specification.AssertionBuilder) {
	return func(builder *specification.AssertionBuilder) {
		builder.WithMethod(assertion.Method)

		for _, assert := range assertion.Assert {
			builder.WithAssert(assert.Actual, assert.Expected)
		}
	}
}

func buildHTTP(http httpSchema) func(builder *specification.HTTPBuilder) {
	return func(builder *specification.HTTPBuilder) {
		builder.
			WithRequest(buildHTTPRequest(http.Request)).
			WithResponse(buildHTTPResponse(http.Response))
	}
}

func buildHTTPRequest(request httpRequestSchema) func(builder *specification.HTTPRequestBuilder) {
	return func(builder *specification.HTTPRequestBuilder) {
		builder.
			WithMethod(request.Method).
			WithURL(request.URL).
			WithContentType(request.ContentType).
			WithBody(request.Body)
	}
}

func buildHTTPResponse(response httpResponseSchema) func(builder *specification.HTTPResponseBuilder) {
	return func(builder *specification.HTTPResponseBuilder) {
		builder.
			WithAllowedCodes(response.AllowedCodes).
			WithAllowedContentType(response.AllowedContentType)
	}
}
