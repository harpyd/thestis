package specification

import (
	"go.uber.org/multierr"

	"github.com/harpyd/thestis/pkg/deepcopy"
)

type (
	Builder struct {
		id          string
		author      string
		title       string
		description string
		stories     map[string]Story

		err error
	}

	StoryBuilder struct {
		description string
		asA         string
		inOrderTo   string
		wantTo      string
		scenarios   map[string]Scenario

		err error
	}

	ScenarioBuilder struct {
		description string
		theses      map[string]Thesis

		err error
	}

	ThesisBuilder struct {
		statement Statement
		http      HTTP
		assertion Assertion

		err error
	}

	AssertionBuilder struct {
		method  AssertionMethod
		asserts []Assert

		err error
	}

	HTTPBuilder struct {
		request  HTTPRequest
		response HTTPResponse

		err error
	}

	HTTPRequestBuilder struct {
		method      HTTPMethod
		url         string
		contentType ContentType
		body        map[string]interface{}

		err error
	}

	HTTPResponseBuilder struct {
		allowedCodes       []int
		allowedContentType ContentType

		err error
	}
)

func NewBuilder() *Builder {
	return &Builder{
		stories: make(map[string]Story),
	}
}

func (b *Builder) Build() (*Specification, error) {
	spec := &Specification{
		id:          b.id,
		author:      b.author,
		title:       b.title,
		description: b.description,
		stories:     make(map[string]Story, len(b.stories)),
	}

	for slug, story := range b.stories {
		spec.stories[slug] = story
	}

	return spec, NewBuildSpecificationError(b.err)
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
		scenarios: make(map[string]Scenario),
	}
}

func (b *StoryBuilder) Build(slug string) (Story, error) {
	if slug == "" {
		return Story{}, NewStoryEmptySlugError()
	}

	story := Story{
		slug:        slug,
		description: b.description,
		asA:         b.asA,
		inOrderTo:   b.inOrderTo,
		wantTo:      b.wantTo,
		scenarios:   make(map[string]Scenario),
	}

	for slug, scenario := range b.scenarios {
		story.scenarios[slug] = scenario
	}

	return story, NewBuildStoryError(b.err, slug)
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
		theses: make(map[string]Thesis),
	}
}

func (b *ScenarioBuilder) Build(slug string) (Scenario, error) {
	if slug == "" {
		return Scenario{}, NewScenarioEmptySlugError()
	}

	scenario := Scenario{
		slug:        slug,
		description: b.description,
		theses:      make(map[string]Thesis),
	}

	for slug, thesis := range b.theses {
		scenario.theses[slug] = thesis
	}

	return scenario, NewBuildScenarioError(b.err, slug)
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
	if slug == "" {
		return Thesis{}, NewThesisEmptySlugError()
	}

	return Thesis{
		slug:      slug,
		statement: b.statement,
		http:      b.http,
		assertion: b.assertion,
	}, NewBuildThesisError(b.err, slug)
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
	assertion := Assertion{
		method:  b.method,
		asserts: make([]Assert, len(b.asserts)),
	}

	copy(assertion.asserts, b.asserts)

	return assertion, b.err
}

func (b *AssertionBuilder) WithMethod(method string) *AssertionBuilder {
	m, err := newAssertionMethodFromString(method)
	b.method = m
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *AssertionBuilder) WithAssert(actual string, expected interface{}) *AssertionBuilder {
	b.asserts = append(b.asserts, Assert{
		actual:   actual,
		expected: expected,
	})

	return b
}

func NewHTTPBuilder() *HTTPBuilder {
	return &HTTPBuilder{}
}

func (b *HTTPBuilder) Build() (HTTP, error) {
	return HTTP{
		request:  b.request,
		response: b.response,
	}, b.err
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
	return HTTPRequest{
		method:      b.method,
		url:         b.url,
		contentType: b.contentType,
		body:        deepcopy.StringInterfaceMap(b.body),
	}, b.err
}

func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
	m, err := newHTTPMethodFromString(method)
	b.method = m
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
	b.url = url

	return b
}

func (b *HTTPRequestBuilder) WithContentType(contentType string) *HTTPRequestBuilder {
	ct, err := newContentTypeFromString(contentType)
	b.contentType = ct
	b.err = multierr.Append(b.err, err)

	return b
}

func (b *HTTPRequestBuilder) WithBody(body map[string]interface{}) *HTTPRequestBuilder {
	b.body = body

	return b
}

func NewHTTPResponseBuilder() *HTTPResponseBuilder {
	return &HTTPResponseBuilder{}
}

func (b *HTTPResponseBuilder) Build() (HTTPResponse, error) {
	return HTTPResponse{
		allowedCodes:       deepcopy.IntSlice(b.allowedCodes),
		allowedContentType: b.allowedContentType,
	}, b.err
}

func (b *HTTPResponseBuilder) WithAllowedCodes(allowedCodes []int) *HTTPResponseBuilder {
	b.allowedCodes = allowedCodes

	return b
}

func (b *HTTPResponseBuilder) WithAllowedContentType(allowedContentType string) *HTTPResponseBuilder {
	act, err := newContentTypeFromString(allowedContentType)
	b.allowedContentType = act
	b.err = multierr.Append(b.err, err)

	return b
}
