package specification

import (
	"github.com/harpyd/thestis/pkg/deepcopy"
	"go.uber.org/multierr"
)

type (
	Builder struct {
		Specification

		err error
	}

	StoryBuilder struct {
		Story

		err error
	}

	ScenarioBuilder struct {
		Scenario

		err error
	}

	ThesisBuilder struct {
		Thesis

		err error
	}

	HTTPBuilder struct {
		HTTP

		err error
	}

	AssertionBuilder struct {
		Assertion

		err error
	}

	HTTPRequestBuilder struct {
		HTTPRequest

		err error
	}

	HTTPResponseBuilder struct {
		HTTPResponse

		err error
	}
)

func NewBuilder() *Builder {
	return &Builder{
		Specification: Specification{
			stories: make(map[string]Story),
		},
	}
}

func (b *Builder) Build() (*Specification, error) {
	return &b.Specification, b.err
}

func (b *Builder) WithAuthor(author string) *Builder {
	b.author = author

	return b
}

func (b *Builder) WithTitle(title string) *Builder {
	b.title = title

	return b
}

func (b *Builder) WithDescription(description string) *Builder {
	b.description = description

	return b
}

func (b *Builder) WithStory(slug string, buildFn func(b *StoryBuilder)) *Builder {
	sb := NewStoryBuilder()
	buildFn(sb)

	story, err := sb.Build(slug)
	b.stories[slug] = story
	b.err = multierr.Append(b.err, err)

	return b
}

func NewStoryBuilder() *StoryBuilder {
	return &StoryBuilder{
		Story: Story{
			scenarios: make(map[string]Scenario),
		},
	}
}

func (b *StoryBuilder) Build(slug string) (Story, error) {
	b.slug = slug

	return b.Story, b.err
}

func (b *StoryBuilder) WithDescription(description string) *StoryBuilder {
	b.description = description

	return b
}

func (b *StoryBuilder) WithAsA(asA string) *StoryBuilder {
	b.asA = asA

	return b
}

func (b *StoryBuilder) WithInOrderTo(inOrderTo string) *StoryBuilder {
	b.inOrderTo = inOrderTo

	return b
}

func (b *StoryBuilder) WithWantTo(wantTo string) *StoryBuilder {
	b.wantTo = wantTo

	return b
}

func (b *StoryBuilder) WithScenario(slug string, buildFn func(b *ScenarioBuilder)) *StoryBuilder {
	sb := NewScenarioBuilder()
	buildFn(sb)

	scenario, err := sb.Build(slug)
	b.scenarios[slug] = scenario
	b.err = multierr.Append(b.err, err)

	return b
}

func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{
		Scenario: Scenario{
			theses: make(map[string]Thesis),
		},
	}
}

func (b *ScenarioBuilder) Build(slug string) (Scenario, error) {
	b.slug = slug

	return b.Scenario, b.err
}

func (b *ScenarioBuilder) WithDescription(description string) *ScenarioBuilder {
	b.description = description

	return b
}

func (b *ScenarioBuilder) WithThesis(slug string, buildFn func(b *ThesisBuilder)) *ScenarioBuilder {
	tb := NewThesisBuilder()
	buildFn(tb)

	thesis, err := tb.Build(slug)
	b.theses[slug] = thesis
	b.err = multierr.Append(b.err, err)

	return b
}

func NewThesisBuilder() *ThesisBuilder {
	return &ThesisBuilder{}
}

func (b *ThesisBuilder) Build(slug string) (Thesis, error) {
	b.slug = slug

	return b.Thesis, b.err
}

func (b *ThesisBuilder) WithStatement(keyword string, behavior string) *ThesisBuilder {
	kw, err := newKeywordFromString(keyword)
	b.statement = Statement{
		keyword:  kw,
		behavior: behavior,
	}
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *ThesisBuilder) WithAssertion(buildFn func(b *AssertionBuilder)) *ThesisBuilder {
	ab := NewAssertionBuilder()
	buildFn(ab)

	assertion, err := ab.Build()
	b.assertion = assertion
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *ThesisBuilder) WithHTTP(buildFn func(b *HTTPBuilder)) *ThesisBuilder {
	hb := NewHTTPBuilder()
	buildFn(hb)

	http, err := hb.Build()
	b.http = http
	b.err = multierr.Append(b.err, err)

	return b
}

func NewAssertionBuilder() *AssertionBuilder {
	return &AssertionBuilder{}
}

func (b *AssertionBuilder) Build() (Assertion, error) {
	return b.Assertion, b.err
}

func (b *AssertionBuilder) WithMethod(method string) *AssertionBuilder {
	m, err := newAssertionMethodFromString(method)
	b.method = m
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *AssertionBuilder) WithAssert(expected string, actual interface{}) *AssertionBuilder {
	b.asserts = append(b.asserts, Assert{
		expected: expected,
		actual:   deepcopy.Interface(actual),
	})

	return b
}

func NewHTTPBuilder() *HTTPBuilder {
	return &HTTPBuilder{}
}

func (b *HTTPBuilder) Build() (HTTP, error) {
	return b.HTTP, b.err
}

func (b *HTTPBuilder) WithMethod(method string) *HTTPBuilder {
	m, err := newHTTPMethodFromString(method)
	b.method = m
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *HTTPBuilder) WithURL(url string) *HTTPBuilder {
	b.url = url

	return b
}

func (b *HTTPBuilder) WithRequest(buildFn func(b *HTTPRequestBuilder)) *HTTPBuilder {
	rb := NewHTTPRequestBuilder()
	buildFn(rb)

	request, err := rb.Build()
	b.request = request
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *HTTPBuilder) WithResponse(buildFn func(b *HTTPResponseBuilder)) *HTTPBuilder {
	rb := NewHTTPResponseBuilder()
	buildFn(rb)

	response, err := rb.Build()
	b.response = response
	b.err = multierr.Append(b.err, err)

	return b
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{}
}

func (b *HTTPRequestBuilder) Build() (HTTPRequest, error) {
	return b.HTTPRequest, b.err
}

func (b *HTTPRequestBuilder) WithContentType(contentType string) *HTTPRequestBuilder {
	ct, err := newContentTypeFromString(contentType)
	b.contentType = ct
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *HTTPRequestBuilder) WithBody(body map[string]interface{}) *HTTPRequestBuilder {
	b.body = deepcopy.StringInterfaceMap(body)

	return b
}

func NewHTTPResponseBuilder() *HTTPResponseBuilder {
	return &HTTPResponseBuilder{}
}

func (b *HTTPResponseBuilder) Build() (HTTPResponse, error) {
	return b.HTTPResponse, b.err
}

func (b *HTTPResponseBuilder) WithAllowedCodes(allowedCodes []int) *HTTPResponseBuilder {
	b.allowedCodes = deepcopy.IntSlice(allowedCodes)

	return b
}

func (b *HTTPResponseBuilder) WithAllowedContentType(allowedContentType string) *HTTPResponseBuilder {
	act, err := newContentTypeFromString(allowedContentType)
	b.allowedContentType = act
	b.err = multierr.Append(b.err, err)

	return b
}
