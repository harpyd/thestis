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

func marshalToSpecificationDocument(spec *specification.Specification) specificationDocument {
	stories := spec.Stories()

	return specificationDocument{
		ID:             spec.ID(),
		OwnerID:        spec.OwnerID(),
		TestCampaignID: spec.TestCampaignID(),
		LoadedAt:       spec.LoadedAt(),
		Author:         spec.Author(),
		Title:          spec.Title(),
		Description:    spec.Description(),
		Stories:        marshalToStoryDocuments(stories),
	}
}

func marshalToStoryDocuments(stories []specification.Story) []storyDocument {
	documents := make([]storyDocument, 0, len(stories))

	for _, story := range stories {
		documents = append(documents, storyDocument{
			Slug:        story.Slug().Story(),
			Description: story.Description(),
			AsA:         story.AsA(),
			InOrderTo:   story.InOrderTo(),
			WantTo:      story.WantTo(),
			Scenarios:   marshalToScenarioDocuments(story.Scenarios()),
		})
	}

	return documents
}

func marshalToScenarioDocuments(scenarios []specification.Scenario) []scenarioDocument {
	documents := make([]scenarioDocument, 0, len(scenarios))

	for _, scenario := range scenarios {
		documents = append(documents, scenarioDocument{
			Slug:        scenario.Slug().Scenario(),
			Description: scenario.Description(),
			Theses:      marshalToThesisDocuments(scenario.Theses()),
		})
	}

	return documents
}

func marshalToThesisDocuments(theses []specification.Thesis) []thesisDocument {
	documents := make([]thesisDocument, 0, len(theses))
	for _, thesis := range theses {
		documents = append(documents, marshalToThesisDocument(thesis))
	}

	return documents
}

func marshalToThesisDocument(thesis specification.Thesis) thesisDocument {
	return thesisDocument{
		Slug:  thesis.Slug().Thesis(),
		After: marshalSlugsToStrings(thesis.Dependencies()),
		Statement: statementDocument{
			Stage:    thesis.Stage(),
			Behavior: thesis.Behavior(),
		},
		HTTP:      marshalToHTTPDocument(thesis.HTTP()),
		Assertion: marshalToAssertionDocument(thesis.Assertion()),
	}
}

func marshalSlugsToStrings(slugs []specification.Slug) []string {
	res := make([]string, 0, len(slugs))

	for _, s := range slugs {
		res = append(res, s.Thesis())
	}

	return res
}

func marshalToHTTPDocument(http specification.HTTP) httpDocument {
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

func marshalToAssertionDocument(assertion specification.Assertion) assertionDocument {
	return assertionDocument{
		Method:  assertion.Method(),
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
		b.WithStory(story.Slug, story.unmarshalToStoryBuildFn())
	}

	return b.ErrlessBuild()
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
			WithStatement(d.Statement.Stage, d.Statement.Behavior).
			WithHTTP(d.HTTP.unmarshalToHTTPBuildFn()).
			WithAssertion(d.Assertion.unmarshalToAssertionBuildFn())

		for _, after := range d.After {
			builder.WithDependency(after)
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

func (d specificationDocument) unmarshalToSpecificSpecification() app.SpecificSpecification {
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
		spec.Stories = append(spec.Stories, s.unmarshalToStory())
	}

	return spec
}

func (d storyDocument) unmarshalToStory() app.Story {
	story := app.Story{
		Slug:        d.Slug,
		Description: d.Description,
		AsA:         d.AsA,
		InOrderTo:   d.InOrderTo,
		WantTo:      d.WantTo,
		Scenarios:   make([]app.Scenario, 0, len(d.Scenarios)),
	}

	for _, s := range d.Scenarios {
		story.Scenarios = append(story.Scenarios, s.unmarshalToScenario())
	}

	return story
}

func (d scenarioDocument) unmarshalToScenario() app.Scenario {
	scenario := app.Scenario{
		Slug:        d.Slug,
		Description: d.Description,
		Theses:      make([]app.Thesis, 0, len(d.Theses)),
	}

	for _, t := range d.Theses {
		scenario.Theses = append(scenario.Theses, t.unmarshalToThesis())
	}

	return scenario
}

func (d thesisDocument) unmarshalToThesis() app.Thesis {
	return app.Thesis{
		Slug:      d.Slug,
		After:     d.After,
		Statement: d.Statement.unmarshalToStatement(),
		HTTP:      d.HTTP.unmarshalToHTTP(),
		Assertion: d.Assertion.unmarshalToAssertion(),
	}
}

func (d statementDocument) unmarshalToStatement() app.Statement {
	return app.Statement{
		Stage:    d.Stage.String(),
		Behavior: d.Behavior,
	}
}

func (d httpDocument) unmarshalToHTTP() app.HTTP {
	return app.HTTP{
		Request:  d.Request.unmarshalToHTTPRequest(),
		Response: d.Response.unmarshalToHTTPResponse(),
	}
}

func (d httpRequestDocument) unmarshalToHTTPRequest() app.HTTPRequest {
	return app.HTTPRequest{
		Method:      d.Method.String(),
		URL:         d.URL,
		ContentType: d.ContentType.String(),
		Body:        d.Body,
	}
}

func (d httpResponseDocument) unmarshalToHTTPResponse() app.HTTPResponse {
	return app.HTTPResponse{
		AllowedCodes:       d.AllowedCodes,
		AllowedContentType: d.AllowedContentType.String(),
	}
}

func (d assertionDocument) unmarshalToAssertion() app.Assertion {
	assert := app.Assertion{
		Method:  d.Method.String(),
		Asserts: make([]app.Assert, 0, len(d.Asserts)),
	}

	for _, a := range d.Asserts {
		assert.Asserts = append(assert.Asserts, a.unmarshalToAssert())
	}

	return assert
}

func (d assertDocument) unmarshalToAssert() app.Assert {
	return app.Assert{
		Actual:   d.Actual,
		Expected: d.Expected,
	}
}
