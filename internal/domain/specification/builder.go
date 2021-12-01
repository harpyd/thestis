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
		stories     []storyFactory
	}

	storyFactory func() (Story, error)

	StoryBuilder struct {
		description string
		asA         string
		inOrderTo   string
		wantTo      string
		scenarios   []scenarioFactory
	}

	scenarioFactory func() (Scenario, error)

	ScenarioBuilder struct {
		description string
		theses      []thesisFactory
	}

	thesisFactory func() (Thesis, error)

	ThesisBuilder struct {
		after            []string
		keyword          string
		behavior         string
		httpBuilder      *HTTPBuilder
		assertionBuilder *AssertionBuilder
	}

	AssertionBuilder struct {
		method  string
		asserts []Assert
	}

	HTTPBuilder struct {
		requestBuilder  *HTTPRequestBuilder
		responseBuilder *HTTPResponseBuilder
	}

	HTTPRequestBuilder struct {
		method      string
		url         string
		contentType string
		body        map[string]interface{}
	}

	HTTPResponseBuilder struct {
		allowedCodes       []int
		allowedContentType string
	}
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() (*Specification, error) {
	spec := &Specification{
		id:          b.id,
		author:      b.author,
		title:       b.title,
		description: b.description,
		stories:     make(map[string]Story, len(b.stories)),
	}

	var err error

	for _, stryFactory := range b.stories {
		stry, stryErr := stryFactory()
		if _, ok := spec.stories[stry.Slug()]; ok {
			err = multierr.Append(err, NewStorySlugAlreadyExistsError(stry.Slug()))

			continue
		}

		err = multierr.Append(err, stryErr)

		spec.stories[stry.Slug()] = stry
	}

	return spec, NewBuildSpecificationError(err)
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

	b.stories = append(b.stories, func() (Story, error) {
		return sb.Build(slug)
	})

	return b
}

func NewStoryBuilder() *StoryBuilder {
	return &StoryBuilder{}
}

func (b *StoryBuilder) Build(slug string) (Story, error) {
	if slug == "" {
		return Story{}, NewStoryEmptySlugError()
	}

	stry := Story{
		slug:        slug,
		description: b.description,
		asA:         b.asA,
		inOrderTo:   b.inOrderTo,
		wantTo:      b.wantTo,
		scenarios:   make(map[string]Scenario),
	}

	var err error

	for _, scnFactory := range b.scenarios {
		scn, scnErr := scnFactory()
		if _, ok := stry.scenarios[scn.Slug()]; ok {
			err = multierr.Append(err, NewScenarioSlugAlreadyExistsError(scn.Slug()))

			continue
		}

		err = multierr.Append(err, scnErr)

		stry.scenarios[scn.Slug()] = scn
	}

	return stry, NewBuildStoryError(err, slug)
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

	b.scenarios = append(b.scenarios, func() (Scenario, error) {
		return sb.Build(slug)
	})

	return b
}

func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{}
}

func (b *ScenarioBuilder) Build(slug string) (Scenario, error) {
	if slug == "" {
		return Scenario{}, NewScenarioEmptySlugError()
	}

	scn := Scenario{
		slug:        slug,
		description: b.description,
		theses:      make(map[string]Thesis),
	}

	var err error

	for _, thsisFactory := range b.theses {
		thsis, thsisErr := thsisFactory()
		if _, ok := scn.theses[thsis.Slug()]; ok {
			err = multierr.Append(err, NewThesisSlugAlreadyExistsError(thsis.Slug()))

			continue
		}

		err = multierr.Append(err, thsisErr)

		scn.theses[thsis.Slug()] = thsis
	}

	return scn, NewBuildScenarioError(err, slug)
}

func (b *ScenarioBuilder) WithDescription(description string) *ScenarioBuilder {
	b.description = description

	return b
}

func (b *ScenarioBuilder) WithThesis(slug string, buildFn func(b *ThesisBuilder)) *ScenarioBuilder {
	tb := NewThesisBuilder()
	buildFn(tb)

	b.theses = append(b.theses, func() (Thesis, error) {
		return tb.Build(slug)
	})

	return b
}

func NewThesisBuilder() *ThesisBuilder {
	return &ThesisBuilder{
		assertionBuilder: NewAssertionBuilder(),
		httpBuilder:      NewHTTPBuilder(),
	}
}

func (b *ThesisBuilder) Build(slug string) (Thesis, error) {
	if slug == "" {
		return Thesis{}, NewThesisEmptySlugError()
	}

	kw, keywordErr := newKeywordFromString(b.keyword)
	http, httpErr := b.httpBuilder.Build()
	assertion, assertionErr := b.assertionBuilder.Build()

	thsis := Thesis{
		slug:  slug,
		after: make([]string, len(b.after)),
		statement: Statement{
			keyword:  kw,
			behavior: b.behavior,
		},
		http:      http,
		assertion: assertion,
	}

	copy(thsis.after, b.after)

	return thsis, NewBuildThesisError(multierr.Combine(
		keywordErr,
		httpErr,
		assertionErr,
	), slug)
}

func (b *ThesisBuilder) WithAfter(after string) *ThesisBuilder {
	b.after = append(b.after, after)

	return b
}

func (b *ThesisBuilder) WithStatement(keyword string, behavior string) *ThesisBuilder {
	b.keyword = keyword
	b.behavior = behavior

	return b
}

func (b *ThesisBuilder) WithAssertion(buildFn func(b *AssertionBuilder)) *ThesisBuilder {
	b.assertionBuilder.Reset()
	buildFn(b.assertionBuilder)

	return b
}

func (b *ThesisBuilder) WithHTTP(buildFn func(b *HTTPBuilder)) *ThesisBuilder {
	b.httpBuilder.Reset()
	buildFn(b.httpBuilder)

	return b
}

func NewAssertionBuilder() *AssertionBuilder {
	return &AssertionBuilder{}
}

func (b *AssertionBuilder) Build() (Assertion, error) {
	method, err := newAssertionMethodFromString(b.method)

	assertion := Assertion{
		method:  method,
		asserts: make([]Assert, len(b.asserts)),
	}

	copy(assertion.asserts, b.asserts)

	return assertion, err
}

func (b *AssertionBuilder) Reset() {
	b.method = ""
	b.asserts = nil
}

func (b *AssertionBuilder) WithMethod(method string) *AssertionBuilder {
	b.method = method

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
	return &HTTPBuilder{
		requestBuilder:  NewHTTPRequestBuilder(),
		responseBuilder: NewHTTPResponseBuilder(),
	}
}

func (b *HTTPBuilder) Build() (HTTP, error) {
	request, requestErr := b.requestBuilder.Build()
	response, responseErr := b.responseBuilder.Build()

	return HTTP{
		request:  request,
		response: response,
	}, multierr.Combine(requestErr, responseErr)
}

func (b *HTTPBuilder) Reset() {
	b.requestBuilder.Reset()
	b.responseBuilder.Reset()
}

func (b *HTTPBuilder) WithRequest(buildFn func(b *HTTPRequestBuilder)) *HTTPBuilder {
	b.requestBuilder.Reset()
	buildFn(b.requestBuilder)

	return b
}

func (b *HTTPBuilder) WithResponse(buildFn func(b *HTTPResponseBuilder)) *HTTPBuilder {
	b.responseBuilder.Reset()
	buildFn(b.responseBuilder)

	return b
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{}
}

func (b *HTTPRequestBuilder) Build() (HTTPRequest, error) {
	method, methodErr := newHTTPMethodFromString(b.method)
	ctype, ctypeErr := newContentTypeFromString(b.contentType)

	return HTTPRequest{
		method:      method,
		url:         b.url,
		contentType: ctype,
		body:        deepcopy.StringInterfaceMap(b.body),
	}, multierr.Combine(methodErr, ctypeErr)
}

func (b *HTTPRequestBuilder) Reset() {
	b.method = ""
	b.url = ""
	b.contentType = ""
	b.body = nil
}

func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
	b.method = method

	return b
}

func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
	b.url = url

	return b
}

func (b *HTTPRequestBuilder) WithContentType(contentType string) *HTTPRequestBuilder {
	b.contentType = contentType

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
	allowedContentType, err := newContentTypeFromString(b.allowedContentType)

	return HTTPResponse{
		allowedCodes:       deepcopy.IntSlice(b.allowedCodes),
		allowedContentType: allowedContentType,
	}, err
}

func (b *HTTPResponseBuilder) Reset() {
	b.allowedCodes = nil
	b.allowedContentType = ""
}

func (b *HTTPResponseBuilder) WithAllowedCodes(allowedCodes []int) *HTTPResponseBuilder {
	b.allowedCodes = allowedCodes

	return b
}

func (b *HTTPResponseBuilder) WithAllowedContentType(allowedContentType string) *HTTPResponseBuilder {
	b.allowedContentType = allowedContentType

	return b
}
