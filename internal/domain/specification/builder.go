package specification

type (
	Builder struct {
		Specification
	}

	StoryBuilder struct {
		Story
	}

	ScenarioBuilder struct {
		Scenario
	}

	ThesisBuilder struct {
		Thesis
	}

	HTTPBuilder struct {
		HTTP
	}

	AssertionBuilder struct {
		Assertion
	}

	HTTPRequestBuilder struct {
		HTTPRequest
	}

	HTTPResponseBuilder struct {
		HTTPResponse
	}
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() *Specification {
	return &b.Specification
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

	b.stories[slug] = sb.Build(slug)

	return b
}

func NewStoryBuilder() *StoryBuilder {
	return &StoryBuilder{}
}

func (b *StoryBuilder) Build(slug string) Story {
	b.slug = slug

	return b.Story
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

	b.scenarios[slug] = sb.Build(slug)

	return b
}

func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{}
}

func (b *ScenarioBuilder) Build(slug string) Scenario {
	b.slug = slug

	return b.Scenario
}

func (b *ScenarioBuilder) WithDescription(description string) *ScenarioBuilder {
	b.description = description

	return b
}

func (b *ScenarioBuilder) WithThesis(slug string, buildFn func(b *ThesisBuilder)) *ScenarioBuilder {
	tb := NewThesisBuilder()
	buildFn(tb)

	b.theses[slug] = tb.Build(slug)

	return b
}

func NewThesisBuilder() *ThesisBuilder {
	return &ThesisBuilder{}
}

func (b *ThesisBuilder) Build(slug string) Thesis {
	b.slug = slug

	return b.Thesis
}

func (b *ThesisBuilder) WithStatement(keyword string, behavior string) *ThesisBuilder {
	kw, _ := newKeywordFromString(keyword)
	b.statement = Statement{
		keyword:  kw,
		behavior: behavior,
	}

	return b
}

func (b *ThesisBuilder) WithHTTP(buildFn func(b *HTTPBuilder)) *ThesisBuilder {
	hb := NewHTTPBuilder()
	buildFn(hb)

	b.http = hb.Build()

	return b
}

func NewHTTPBuilder() *HTTPBuilder {
	return &HTTPBuilder{}
}

func (b *HTTPBuilder) Build() HTTP {
	return b.HTTP
}

func (b *HTTPBuilder) WithMethod(method string) *HTTPBuilder {
	b.method, _ = newHTTPMethodFromString(method)

	return b
}

func (b *HTTPBuilder) WithURL(url string) *HTTPBuilder {
	b.url = url

	return b
}

func (b *HTTPBuilder) WithRequest(buildFn func(b *HTTPRequestBuilder)) *HTTPBuilder {
	rb := NewHTTPRequestBuilder()
	buildFn(rb)

	b.request = rb.Build()

	return b
}

func (b *HTTPBuilder) WithResponse(buildFn func(b *HTTPResponseBuilder)) *HTTPBuilder {
	rb := NewHTTPResponseBuilder()
	buildFn(rb)

	b.response = rb.Build()

	return b
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{}
}

func (b *HTTPRequestBuilder) Build() HTTPRequest {
	return b.HTTPRequest
}

func (b *HTTPRequestBuilder) WithContentType(contentType string) *HTTPRequestBuilder {
	b.contentType, _ = newContentTypeFromString(contentType)

	return b
}

func NewHTTPResponseBuilder() *HTTPResponseBuilder {
	return &HTTPResponseBuilder{}
}

func (b *HTTPResponseBuilder) Build() HTTPResponse {
	return b.HTTPResponse
}
