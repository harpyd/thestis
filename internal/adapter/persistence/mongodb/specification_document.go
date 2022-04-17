package mongodb

import (
	"time"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	specificationDocument struct {
		_              string          `bson:"_id,omitempty"`
		ID             string          `bson:"id,omitempty"`
		OwnerID        string          `bson:"ownerId"`
		TestCampaignID string          `bson:"testCampaignId"`
		LoadedAt       time.Time       `bson:"loadedAt"`
		Author         string          `bson:"author"`
		Title          string          `bson:"title"`
		Description    string          `bson:"description"`
		Stories        []storyDocument `bson:"stories"`
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
		Stage    specification.Stage `bson:"stage"`
		Behavior string              `bson:"behavior"`
	}

	httpDocument struct {
		Request  httpRequestDocument  `bson:"request"`
		Response httpResponseDocument `bson:"response"`
	}

	httpRequestDocument struct {
		Method      specification.HTTPMethod  `bson:"method"`
		URL         string                    `bson:"url"`
		ContentType specification.ContentType `bson:"contentType"`
		Body        map[string]interface{}    `bson:"body"`
	}

	httpResponseDocument struct {
		AllowedCodes       []int                     `bson:"allowedCodes"`
		AllowedContentType specification.ContentType `bson:"allowedContentType"`
	}

	assertionDocument struct {
		Method  specification.AssertionMethod `bson:"method"`
		Asserts []assertDocument              `bson:"asserts"`
	}

	assertDocument struct {
		Actual   string      `bson:"actual"`
		Expected interface{} `bson:"expected"`
	}
)

func newSpecificationDocument(spec *specification.Specification) specificationDocument {
	stories := spec.Stories()

	return specificationDocument{
		ID:             spec.ID(),
		OwnerID:        spec.OwnerID(),
		TestCampaignID: spec.TestCampaignID(),
		LoadedAt:       spec.LoadedAt(),
		Author:         spec.Author(),
		Title:          spec.Title(),
		Description:    spec.Description(),
		Stories:        newStoryDocuments(stories),
	}
}

func newStoryDocuments(stories []specification.Story) []storyDocument {
	documents := make([]storyDocument, 0, len(stories))

	for _, story := range stories {
		documents = append(documents, storyDocument{
			Slug:        story.Slug().Story(),
			Description: story.Description(),
			AsA:         story.AsA(),
			InOrderTo:   story.InOrderTo(),
			WantTo:      story.WantTo(),
			Scenarios:   newScenarioDocuments(story.Scenarios()),
		})
	}

	return documents
}

func newScenarioDocuments(scenarios []specification.Scenario) []scenarioDocument {
	documents := make([]scenarioDocument, 0, len(scenarios))

	for _, scenario := range scenarios {
		documents = append(documents, scenarioDocument{
			Slug:        scenario.Slug().Scenario(),
			Description: scenario.Description(),
			Theses:      newThesisDocuments(scenario.Theses()),
		})
	}

	return documents
}

func newThesisDocuments(theses []specification.Thesis) []thesisDocument {
	documents := make([]thesisDocument, 0, len(theses))
	for _, thesis := range theses {
		documents = append(documents, newThesisDocument(thesis))
	}

	return documents
}

func newThesisDocument(thesis specification.Thesis) thesisDocument {
	return thesisDocument{
		Slug:  thesis.Slug().Thesis(),
		After: mapSlugsToStrings(thesis.Dependencies()),
		Statement: statementDocument{
			Stage:    thesis.Stage(),
			Behavior: thesis.Behavior(),
		},
		HTTP:      newHTTPDocument(thesis.HTTP()),
		Assertion: newAssertionDocument(thesis.Assertion()),
	}
}

func mapSlugsToStrings(slugs []specification.Slug) []string {
	res := make([]string, 0, len(slugs))

	for _, s := range slugs {
		res = append(res, s.Thesis())
	}

	return res
}

func newHTTPDocument(http specification.HTTP) httpDocument {
	return httpDocument{
		Request: httpRequestDocument{
			Method:      http.Request().Method(),
			URL:         http.Request().URL(),
			ContentType: http.Request().ContentType(),
			Body:        http.Request().Body(),
		},
		Response: httpResponseDocument{
			AllowedCodes:       http.Response().AllowedCodes(),
			AllowedContentType: http.Response().AllowedContentType(),
		},
	}
}

func newAssertionDocument(assertion specification.Assertion) assertionDocument {
	return assertionDocument{
		Method:  assertion.Method(),
		Asserts: newAssertDocuments(assertion.Asserts()),
	}
}

func newAssertDocuments(asserts []specification.Assert) []assertDocument {
	documents := make([]assertDocument, 0, len(asserts))
	for _, assert := range asserts {
		documents = append(documents, assertDocument{
			Expected: assert.Expected(),
			Actual:   assert.Actual(),
		})
	}

	return documents
}

func (d specificationDocument) toSpecification() *specification.Specification {
	var b specification.Builder

	b.
		WithID(d.ID).
		WithOwnerID(d.OwnerID).
		WithTestCampaignID(d.TestCampaignID).
		WithLoadedAt(d.LoadedAt).
		WithAuthor(d.Author).
		WithTitle(d.Title).
		WithDescription(d.Description)

	for _, story := range d.Stories {
		b.WithStory(story.Slug, story.toStoryBuildFn())
	}

	return b.ErrlessBuild()
}

func (d storyDocument) toStoryBuildFn() func(builder *specification.StoryBuilder) {
	return func(builder *specification.StoryBuilder) {
		builder.
			WithDescription(d.Description).
			WithAsA(d.AsA).
			WithInOrderTo(d.InOrderTo).
			WithWantTo(d.WantTo)

		for _, scenario := range d.Scenarios {
			builder.WithScenario(scenario.Slug, scenario.toScenarioBuildFn())
		}
	}
}

func (d scenarioDocument) toScenarioBuildFn() func(builder *specification.ScenarioBuilder) {
	return func(builder *specification.ScenarioBuilder) {
		builder.WithDescription(d.Description)

		for _, thesis := range d.Theses {
			builder.WithThesis(thesis.Slug, thesis.toThesisBuildFn())
		}
	}
}

func (d thesisDocument) toThesisBuildFn() func(builder *specification.ThesisBuilder) {
	return func(builder *specification.ThesisBuilder) {
		builder.
			WithStatement(d.Statement.Stage, d.Statement.Behavior).
			WithHTTP(d.HTTP.toHTTPBuildFn()).
			WithAssertion(d.Assertion.toAssertionBuildFn())

		for _, after := range d.After {
			builder.WithDependency(after)
		}
	}
}

func (d httpDocument) toHTTPBuildFn() func(builder *specification.HTTPBuilder) {
	return func(builder *specification.HTTPBuilder) {
		builder.
			WithRequest(d.Request.toHTTPRequestBuildFn()).
			WithResponse(d.Response.toHTTPResponseBuildFn())
	}
}

func (d httpRequestDocument) toHTTPRequestBuildFn() func(builder *specification.HTTPRequestBuilder) {
	return func(builder *specification.HTTPRequestBuilder) {
		builder.
			WithMethod(d.Method).
			WithURL(d.URL).
			WithContentType(d.ContentType).
			WithBody(d.Body)
	}
}

func (d httpResponseDocument) toHTTPResponseBuildFn() func(builder *specification.HTTPResponseBuilder) {
	return func(builder *specification.HTTPResponseBuilder) {
		builder.
			WithAllowedCodes(d.AllowedCodes).
			WithAllowedContentType(d.AllowedContentType)
	}
}

func (d assertionDocument) toAssertionBuildFn() func(builder *specification.AssertionBuilder) {
	return func(builder *specification.AssertionBuilder) {
		builder.WithMethod(d.Method)

		for _, assert := range d.Asserts {
			builder.WithAssert(assert.Actual, assert.Expected)
		}
	}
}

func (d specificationDocument) toSpecificSpecification() app.SpecificSpecification {
	spec := app.SpecificSpecification{
		ID:             d.ID,
		TestCampaignID: d.TestCampaignID,
		LoadedAt:       d.LoadedAt,
		Author:         d.Author,
		Title:          d.Title,
		Description:    d.Description,
		Stories:        make([]app.Story, 0, len(d.Stories)),
	}

	for _, s := range d.Stories {
		spec.Stories = append(spec.Stories, s.toStory())
	}

	return spec
}

func (d storyDocument) toStory() app.Story {
	story := app.Story{
		Slug:        d.Slug,
		Description: d.Description,
		AsA:         d.AsA,
		InOrderTo:   d.InOrderTo,
		WantTo:      d.WantTo,
		Scenarios:   make([]app.Scenario, 0, len(d.Scenarios)),
	}

	for _, s := range d.Scenarios {
		story.Scenarios = append(story.Scenarios, s.toScenario())
	}

	return story
}

func (d scenarioDocument) toScenario() app.Scenario {
	scenario := app.Scenario{
		Slug:        d.Slug,
		Description: d.Description,
		Theses:      make([]app.Thesis, 0, len(d.Theses)),
	}

	for _, t := range d.Theses {
		scenario.Theses = append(scenario.Theses, t.toThesis())
	}

	return scenario
}

func (d thesisDocument) toThesis() app.Thesis {
	return app.Thesis{
		Slug:      d.Slug,
		After:     d.After,
		Statement: d.Statement.toStatement(),
		HTTP:      d.HTTP.toHTTP(),
		Assertion: d.Assertion.toAssertion(),
	}
}

func (d statementDocument) toStatement() app.Statement {
	return app.Statement{
		Stage:    d.Stage.String(),
		Behavior: d.Behavior,
	}
}

func (d httpDocument) toHTTP() app.HTTP {
	return app.HTTP{
		Request:  d.Request.toHTTPRequest(),
		Response: d.Response.toHTTPResponse(),
	}
}

func (d httpRequestDocument) toHTTPRequest() app.HTTPRequest {
	return app.HTTPRequest{
		Method:      d.Method.String(),
		URL:         d.URL,
		ContentType: d.ContentType.String(),
		Body:        d.Body,
	}
}

func (d httpResponseDocument) toHTTPResponse() app.HTTPResponse {
	return app.HTTPResponse{
		AllowedCodes:       d.AllowedCodes,
		AllowedContentType: d.AllowedContentType.String(),
	}
}

func (d assertionDocument) toAssertion() app.Assertion {
	assert := app.Assertion{
		Method:  d.Method.String(),
		Asserts: make([]app.Assert, 0, len(d.Asserts)),
	}

	for _, a := range d.Asserts {
		assert.Asserts = append(assert.Asserts, a.toAssert())
	}

	return assert
}

func (d assertDocument) toAssert() app.Assert {
	return app.Assert{
		Actual:   d.Actual,
		Expected: d.Expected,
	}
}
