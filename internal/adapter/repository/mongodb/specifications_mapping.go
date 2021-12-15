package mongodb

import "github.com/harpyd/thestis/internal/domain/specification"

type (
	specificationDocument struct {
		ID          string          `bson:"_id,omitempty"`
		Author      string          `bson:"author"`
		Title       string          `bson:"title"`
		Description string          `bson:"description"`
		Stories     []storyDocument `bson:"stories"`
	}

	storyDocument struct {
		Slug        string             `bson:"slug"`
		Description string             `bson:"description"`
		AsA         string             `bson:"asA"`
		InOrderTo   string             `bson:"inOrderTo"`
		WantTo      string             `bson:"wantTo"`
		Scenarios   []scenarioDocument `bson:"scenarios"`
	}

	scenarioDocument struct {
		Slug        string           `bson:"slug"`
		Description string           `bson:"description"`
		Theses      []thesisDocument `bson:"theses"`
	}

	thesisDocument struct {
		Slug      string            `bson:"slug"`
		After     []string          `bson:"after"`
		Statement statementDocument `bson:"statement"`
		HTTP      httpDocument      `bson:"http"`
		Assertion assertionDocument `bson:"assertion"`
	}

	statementDocument struct {
		Keyword   string `bson:"keyword"`
		Behaviour string `bson:"behaviour"`
	}

	httpDocument struct {
		Request  httpRequestDocument  `bson:"request"`
		Response httpResponseDocument `bson:"response"`
	}

	httpRequestDocument struct {
		Method      string                 `bson:"method"`
		URL         string                 `bson:"url"`
		ContentType string                 `bson:"contentType"`
		Body        map[string]interface{} `bson:"body"`
	}

	httpResponseDocument struct {
		AllowedCodes       []int  `bson:"allowedCodes"`
		AllowedContentType string `bson:"allowedContentType"`
	}

	assertionDocument struct {
		Method  string           `bson:"method"`
		Asserts []assertDocument `bson:"asserts"`
	}

	assertDocument struct {
		Actual   string      `bson:"actual"`
		Expected interface{} `bson:"expected"`
	}
)

func marshalToSpecificationDocument(spec *specification.Specification) specificationDocument {
	stories, _ := spec.Stories()

	return specificationDocument{
		ID:          spec.ID(),
		Author:      spec.Author(),
		Title:       spec.Title(),
		Description: spec.Description(),
		Stories:     marshalToStoryDocuments(stories),
	}
}

func marshalToStoryDocuments(stories []specification.Story) []storyDocument {
	documents := make([]storyDocument, 0, len(stories))

	for _, story := range stories {
		scenarios, _ := story.Scenarios()

		documents = append(documents, storyDocument{
			Slug:        story.Slug(),
			Description: story.Description(),
			AsA:         story.AsA(),
			InOrderTo:   story.InOrderTo(),
			WantTo:      story.WantTo(),
			Scenarios:   marshalToScenarioDocuments(scenarios),
		})
	}

	return documents
}

func marshalToScenarioDocuments(scenarios []specification.Scenario) []scenarioDocument {
	documents := make([]scenarioDocument, 0, len(scenarios))

	for _, scenario := range scenarios {
		theses, _ := scenario.Theses()

		documents = append(documents, scenarioDocument{
			Slug:        scenario.Slug(),
			Description: scenario.Description(),
			Theses:      marshalToThesisDocuments(theses),
		})
	}

	return documents
}

func marshalToThesisDocuments(theses []specification.Thesis) []thesisDocument {
	documents := make([]thesisDocument, 0, len(theses))
	for _, thesis := range theses {
		documents = append(documents, thesisDocument{
			Slug:  thesis.Slug(),
			After: thesis.After(),
			Statement: statementDocument{
				Keyword:   thesis.Statement().Keyword().String(),
				Behaviour: thesis.Statement().Behavior(),
			},
			HTTP:      marshalToHTTPDocument(thesis.HTTP()),
			Assertion: marshalToAssertionDocument(thesis.Assertion()),
		})
	}

	return documents
}

func marshalToHTTPDocument(http specification.HTTP) httpDocument {
	return httpDocument{
		Request: httpRequestDocument{
			Method:      http.Request().Method().String(),
			URL:         http.Request().URL(),
			ContentType: http.Request().ContentType().String(),
			Body:        http.Request().Body(),
		},
		Response: httpResponseDocument{
			AllowedCodes:       http.Response().AllowedCodes(),
			AllowedContentType: http.Response().AllowedContentType().String(),
		},
	}
}

func marshalToAssertionDocument(assertion specification.Assertion) assertionDocument {
	return assertionDocument{
		Method:  assertion.Method().String(),
		Asserts: marshalToAssertDocuments(assertion.Asserts()),
	}
}

func marshalToAssertDocuments(asserts []specification.Assert) []assertDocument {
	documents := make([]assertDocument, 0, len(asserts))
	for _, assert := range asserts {
		documents = append(documents, assertDocument{
			Expected: assert.Expected(),
			Actual:   assert.Actual(),
		})
	}

	return documents
}

func (d specificationDocument) unmarshalToSpecification() *specification.Specification {
	builder := specification.NewBuilder().
		WithID(d.ID).
		WithAuthor(d.Author).
		WithTitle(d.Title).
		WithDescription(d.Description)

	for _, story := range d.Stories {
		builder.WithStory(story.Slug, story.unmarshalToStoryBuildFn())
	}

	return builder.ErrlessBuild()
}

func (d storyDocument) unmarshalToStoryBuildFn() func(builder *specification.StoryBuilder) {
	return func(builder *specification.StoryBuilder) {
		builder.
			WithDescription(d.Description).
			WithAsA(d.AsA).
			WithInOrderTo(d.InOrderTo).
			WithWantTo(d.WantTo)

		for _, scenario := range d.Scenarios {
			builder.WithScenario(scenario.Slug, scenario.unmarshalToScenarioBuildFn())
		}
	}
}

func (d scenarioDocument) unmarshalToScenarioBuildFn() func(builder *specification.ScenarioBuilder) {
	return func(builder *specification.ScenarioBuilder) {
		builder.WithDescription(d.Description)

		for _, thesis := range d.Theses {
			builder.WithThesis(thesis.Slug, thesis.unmarshalToThesisBuildFn())
		}
	}
}

func (d thesisDocument) unmarshalToThesisBuildFn() func(builder *specification.ThesisBuilder) {
	return func(builder *specification.ThesisBuilder) {
		builder.
			WithStatement(d.Statement.Keyword, d.Statement.Behaviour).
			WithHTTP(d.HTTP.unmarshalToHTTPBuildFn()).
			WithAssertion(d.Assertion.unmarshalToAssertionBuildFn())

		for _, after := range d.After {
			builder.WithAfter(after)
		}
	}
}

func (d httpDocument) unmarshalToHTTPBuildFn() func(builder *specification.HTTPBuilder) {
	return func(builder *specification.HTTPBuilder) {
		builder.
			WithRequest(d.Request.unmarshalToHTTPRequestBuildFn()).
			WithResponse(d.Response.unmarshalToHTTPResponseBuildFn())
	}
}

func (d httpRequestDocument) unmarshalToHTTPRequestBuildFn() func(builder *specification.HTTPRequestBuilder) {
	return func(builder *specification.HTTPRequestBuilder) {
		builder.
			WithMethod(d.Method).
			WithURL(d.URL).
			WithContentType(d.ContentType).
			WithBody(d.Body)
	}
}

func (d httpResponseDocument) unmarshalToHTTPResponseBuildFn() func(builder *specification.HTTPResponseBuilder) {
	return func(builder *specification.HTTPResponseBuilder) {
		builder.
			WithAllowedCodes(d.AllowedCodes).
			WithAllowedContentType(d.AllowedContentType)
	}
}

func (d assertionDocument) unmarshalToAssertionBuildFn() func(builder *specification.AssertionBuilder) {
	return func(builder *specification.AssertionBuilder) {
		builder.WithMethod(d.Method)

		for _, assert := range d.Asserts {
			builder.WithAssert(assert.Actual, assert.Expected)
		}
	}
}
