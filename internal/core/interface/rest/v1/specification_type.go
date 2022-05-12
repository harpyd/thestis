package v1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"

	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/query"
	"github.com/harpyd/thestis/internal/core/interface/rest"
)

func decodeSpecificationSourceCommand(
	w http.ResponseWriter,
	r *http.Request,
	testCampaignID string,
) (cmd command.LoadSpecification, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		rest.BadRequest(string(ErrorSlugBadRequest), err, w, r)

		return
	}

	return command.LoadSpecification{
		Content:        content,
		TestCampaignID: testCampaignID,
		LoadedByID:     user.UUID,
	}, true
}

func decodeSpecificSpecificationQuery(
	w http.ResponseWriter,
	r *http.Request,
	specificationID string,
) (qry query.SpecificSpecification, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return query.SpecificSpecification{
		SpecificationID: specificationID,
		UserID:          user.UUID,
	}, true
}

func renderSpecificationResponse(
	w http.ResponseWriter,
	r *http.Request,
	spec query.SpecificSpecificationView,
) {
	response := SpecificationResponse{
		Specification: newSpecification(spec),
		SourceUri:     "",
	}

	render.Respond(w, r, response)
}

func newSpecification(spec query.SpecificSpecificationView) Specification {
	res := Specification{
		Id:             spec.ID,
		TestCampaignId: spec.TestCampaignID,
		LoadedAt:       spec.LoadedAt,
		Author:         &spec.Author,
		Title:          &spec.Title,
		Description:    &spec.Description,
		Stories:        make([]Story, 0, len(spec.Stories)),
	}

	for _, s := range spec.Stories {
		res.Stories = append(res.Stories, newStory(s))
	}

	return res
}

func newStory(story query.StoryView) Story {
	res := Story{
		Slug:        story.Slug,
		Description: &story.Description,
		AsA:         &story.AsA,
		InOrderTo:   &story.InOrderTo,
		WantTo:      &story.WantTo,
		Scenarios:   make([]Scenario, 0, len(story.Scenarios)),
	}

	for _, s := range story.Scenarios {
		res.Scenarios = append(res.Scenarios, newScenario(s))
	}

	return res
}

func newScenario(scenario query.ScenarioView) Scenario {
	res := Scenario{
		Slug:        scenario.Slug,
		Description: &scenario.Description,
		Theses:      make([]Thesis, 0, len(scenario.Theses)),
	}

	for _, t := range scenario.Theses {
		res.Theses = append(res.Theses, newThesis(t))
	}

	return res
}

func newThesis(thesis query.ThesisView) Thesis {
	return Thesis{
		Slug:      thesis.Slug,
		After:     thesis.After,
		Statement: newStatement(thesis.Statement),
		Http:      newHTTP(thesis.HTTP),
		Assertion: newAssertion(thesis.Assertion),
	}
}

func newStatement(statement query.StatementView) Statement {
	return Statement{
		Stage:    statement.Stage,
		Behavior: statement.Behavior,
	}
}

func newHTTP(http query.HTTPView) *Http {
	if http.IsZero() {
		return nil
	}

	return &Http{
		Request:  newHTTPRequest(http.Request),
		Response: newHTTPResponse(http.Response),
	}
}

func newHTTPRequest(request query.HTTPRequestView) *HttpRequest {
	if request.IsZero() {
		return nil
	}

	return &HttpRequest{
		Method:      HttpMethod(request.Method),
		Url:         request.URL,
		ContentType: &request.ContentType,
		Body:        newBody(request.Body),
	}
}

func newBody(body map[string]interface{}) *map[string]interface{} {
	if len(body) == 0 {
		return nil
	}

	return &body
}

func newHTTPResponse(response query.HTTPResponseView) *HttpResponse {
	if response.IsZero() {
		return nil
	}

	return &HttpResponse{
		AllowedCodes:       response.AllowedCodes,
		AllowedContentType: &response.AllowedContentType,
	}
}

func newAssertion(assertion query.AssertionView) *Assertion {
	if assertion.IsZero() {
		return nil
	}

	return &Assertion{
		With:   AssertionMethod(assertion.Method),
		Assert: newAsserts(assertion.Asserts),
	}
}

func newAsserts(asserts []query.AssertView) []Assert {
	res := make([]Assert, 0, len(asserts))

	for _, a := range asserts {
		res = append(res, Assert{
			Actual:   a.Actual,
			Expected: fmt.Sprintf("%v", a.Expected),
		})
	}

	return res
}