package mongodb

import (
	"time"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/entity/specification"
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

func newSpecification(d specificationDocument) *specification.Specification {
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
		b.WithStory(story.Slug, newStoryBuildFn(story))
	}

	return b.ErrlessBuild()
}

func newStoryBuildFn(d storyDocument) func(builder *specification.StoryBuilder) {
	return func(builder *specification.StoryBuilder) {
		builder.
			WithDescription(d.Description).
			WithAsA(d.AsA).
			WithInOrderTo(d.InOrderTo).
			WithWantTo(d.WantTo)

		for _, scenario := range d.Scenarios {
			builder.WithScenario(scenario.Slug, newScenarioBuildFn(scenario))
		}
	}
}

func newScenarioBuildFn(d scenarioDocument) func(builder *specification.ScenarioBuilder) {
	return func(builder *specification.ScenarioBuilder) {
		builder.WithDescription(d.Description)

		for _, thesis := range d.Theses {
			builder.WithThesis(thesis.Slug, newThesisBuildFn(thesis))
		}
	}
}

func newThesisBuildFn(d thesisDocument) func(builder *specification.ThesisBuilder) {
	return func(builder *specification.ThesisBuilder) {
		builder.
			WithStatement(d.Statement.Stage, d.Statement.Behavior).
			WithHTTP(newHTTPBuildFn(d.HTTP)).
			WithAssertion(newAssertionBuildFn(d.Assertion))

		for _, after := range d.After {
			builder.WithDependency(after)
		}
	}
}

func newHTTPBuildFn(d httpDocument) func(builder *specification.HTTPBuilder) {
	return func(builder *specification.HTTPBuilder) {
		builder.
			WithRequest(newHTTPRequestBuildFn(d.Request)).
			WithResponse(newHTTPResponseBuildFn(d.Response))
	}
}

func newHTTPRequestBuildFn(d httpRequestDocument) func(builder *specification.HTTPRequestBuilder) {
	return func(builder *specification.HTTPRequestBuilder) {
		builder.
			WithMethod(d.Method).
			WithURL(d.URL).
			WithContentType(d.ContentType).
			WithBody(d.Body)
	}
}

func newHTTPResponseBuildFn(d httpResponseDocument) func(builder *specification.HTTPResponseBuilder) {
	return func(builder *specification.HTTPResponseBuilder) {
		builder.
			WithAllowedCodes(d.AllowedCodes).
			WithAllowedContentType(d.AllowedContentType)
	}
}

func newAssertionBuildFn(d assertionDocument) func(builder *specification.AssertionBuilder) {
	return func(builder *specification.AssertionBuilder) {
		builder.WithMethod(d.Method)

		for _, assert := range d.Asserts {
			builder.WithAssert(assert.Actual, assert.Expected)
		}
	}
}

func newAppSpecificSpecification(d specificationDocument) app.SpecificSpecification {
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
		spec.Stories = append(spec.Stories, newAppStory(s))
	}

	return spec
}

func newAppStory(d storyDocument) app.Story {
	story := app.Story{
		Slug:        d.Slug,
		Description: d.Description,
		AsA:         d.AsA,
		InOrderTo:   d.InOrderTo,
		WantTo:      d.WantTo,
		Scenarios:   make([]app.Scenario, 0, len(d.Scenarios)),
	}

	for _, s := range d.Scenarios {
		story.Scenarios = append(story.Scenarios, newAppScenario(s))
	}

	return story
}

func newAppScenario(d scenarioDocument) app.Scenario {
	scenario := app.Scenario{
		Slug:        d.Slug,
		Description: d.Description,
		Theses:      make([]app.Thesis, 0, len(d.Theses)),
	}

	for _, t := range d.Theses {
		scenario.Theses = append(scenario.Theses, newAppThesis(t))
	}

	return scenario
}

func newAppThesis(d thesisDocument) app.Thesis {
	return app.Thesis{
		Slug:      d.Slug,
		After:     d.After,
		Statement: newAppStatement(d.Statement),
		HTTP:      newAppHTTP(d.HTTP),
		Assertion: newAppAssertion(d.Assertion),
	}
}

func newAppStatement(d statementDocument) app.Statement {
	return app.Statement{
		Stage:    d.Stage.String(),
		Behavior: d.Behavior,
	}
}

func newAppHTTP(d httpDocument) app.HTTP {
	return app.HTTP{
		Request:  newAppHTTPRequest(d.Request),
		Response: newAppHTTPResponse(d.Response),
	}
}

func newAppHTTPRequest(d httpRequestDocument) app.HTTPRequest {
	return app.HTTPRequest{
		Method:      d.Method.String(),
		URL:         d.URL,
		ContentType: d.ContentType.String(),
		Body:        d.Body,
	}
}

func newAppHTTPResponse(d httpResponseDocument) app.HTTPResponse {
	return app.HTTPResponse{
		AllowedCodes:       d.AllowedCodes,
		AllowedContentType: d.AllowedContentType.String(),
	}
}

func newAppAssertion(d assertionDocument) app.Assertion {
	assertion := app.Assertion{
		Method:  d.Method.String(),
		Asserts: make([]app.Assert, 0, len(d.Asserts)),
	}

	for _, a := range d.Asserts {
		assertion.Asserts = append(assertion.Asserts, newAppAssert(a))
	}

	return assertion
}

func newAppAssert(d assertDocument) app.Assert {
	return app.Assert{
		Actual:   d.Actual,
		Expected: d.Expected,
	}
}
