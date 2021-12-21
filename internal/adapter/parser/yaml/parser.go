package yaml

import (
	"io"

	"gopkg.in/yaml.v3"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationParserService struct{}

func NewSpecificationParserService() SpecificationParserService {
	return SpecificationParserService{}
}

func (s SpecificationParserService) ParseSpecification(reader io.Reader, opts ...app.ParserOption) (*specification.Specification, error) {
	decoder := yaml.NewDecoder(reader)

	var spec specificationSchema
	if err := decoder.Decode(&spec); err != nil {
		return nil, err
	}

	return build(spec, opts)
}

func build(spec specificationSchema, opts []app.ParserOption) (*specification.Specification, error) {
	builder := specification.NewBuilder().
		WithAuthor(spec.Author).
		WithTitle(spec.Title).
		WithDescription(spec.Description)

	for _, opt := range opts {
		opt(builder)
	}

	for slug, story := range spec.Stories {
		builder.WithStory(slug, buildStory(story))
	}

	return builder.Build()
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
			builder.WithStatement("given", thesis.Given)
		case len(thesis.When) > 0:
			builder.WithStatement("when", thesis.When)
		case len(thesis.Then) > 0:
			builder.WithStatement("then", thesis.Then)
		default:
			builder.WithStatement("", "")
		}

		builder.
			WithAssertion(buildAssertion(thesis.Assertion)).
			WithHTTP(buildHTTP(thesis.HTTP))

		for _, after := range thesis.After {
			builder.WithAfter(after)
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
