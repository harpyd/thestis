package v1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/pkg/http/httperr"
)

func unmarshalToSpecificationSourceCommand(
	w http.ResponseWriter,
	r *http.Request,
	testCampaignID string,
) (cmd app.LoadSpecificationCommand, ok bool) {
	user, ok := unmarshalUser(w, r)
	if !ok {
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		httperr.BadRequest(string(ErrorSlugBadRequest), err, w, r)

		return
	}

	return app.LoadSpecificationCommand{
		Content:        content,
		TestCampaignID: testCampaignID,
		LoadedByID:     user.UUID,
	}, true
}

func unmarshalToSpecificSpecificationQuery(
	w http.ResponseWriter,
	r *http.Request,
	specificationID string,
) (qry app.SpecificSpecificationQuery, ok bool) {
	user, ok := unmarshalUser(w, r)
	if !ok {
		return
	}

	return app.SpecificSpecificationQuery{
		SpecificationID: specificationID,
		UserID:          user.UUID,
	}, true
}

func marshalToSpecificationResponse(w http.ResponseWriter, r *http.Request, spec app.SpecificSpecification) {
	response := SpecificationResponse{
		Specification: marshalToSpecification(spec),
		SourceUri:     "",
	}

	render.Respond(w, r, response)
}

func marshalToSpecification(spec app.SpecificSpecification) Specification {
	res := Specification{
		Id:          spec.ID,
		LoadedAt:    spec.LoadedAt,
		Author:      &spec.Author,
		Title:       &spec.Title,
		Description: &spec.Description,
		Stories:     make([]Story, 0, len(spec.Stories)),
	}

	for _, s := range spec.Stories {
		res.Stories = append(res.Stories, marshalToStory(s))
	}

	return res
}

func marshalToStory(story app.Story) Story {
	res := Story{
		Slug:        story.Slug,
		Description: &story.Description,
		AsA:         &story.AsA,
		InOrderTo:   &story.InOrderTo,
		WantTo:      &story.WantTo,
		Scenarios:   make([]Scenario, 0, len(story.Scenarios)),
	}

	for _, s := range story.Scenarios {
		res.Scenarios = append(res.Scenarios, marshalToScenario(s))
	}

	return res
}

func marshalToScenario(scenario app.Scenario) Scenario {
	res := Scenario{
		Slug:        scenario.Slug,
		Description: &scenario.Description,
		Theses:      make([]Thesis, 0, len(scenario.Theses)),
	}

	for _, t := range scenario.Theses {
		res.Theses = append(res.Theses, marshalToThesis(t))
	}

	return res
}

func marshalToThesis(thesis app.Thesis) Thesis {
	return Thesis{
		Slug:      thesis.Slug,
		After:     thesis.After,
		Statement: marshalToStatement(thesis.Statement),
		Http:      marshalToHTTP(thesis.HTTP),
		Assertion: marshalToAssertion(thesis.Assertion),
	}
}

func marshalToStatement(statement app.Statement) Statement {
	return Statement{
		Keyword:  statement.Keyword,
		Behavior: statement.Behavior,
	}
}

func marshalToHTTP(http app.HTTP) *Http {
	if http.IsZero() {
		return nil
	}

	return &Http{
		Request:  marshalToHTTPRequest(http.Request),
		Response: marshalToHTTPResponse(http.Response),
	}
}

func marshalToHTTPRequest(request app.HTTPRequest) *HttpRequest {
	if request.IsZero() {
		return nil
	}

	return &HttpRequest{
		Method:      HttpMethod(request.Method),
		Url:         request.URL,
		ContentType: &request.ContentType,
		Body:        unmarshalBody(request.Body),
	}
}

func unmarshalBody(body map[string]interface{}) *map[string]interface{} {
	if len(body) == 0 {
		return nil
	}

	return &body
}

func marshalToHTTPResponse(response app.HTTPResponse) *HttpResponse {
	if response.IsZero() {
		return nil
	}

	return &HttpResponse{
		AllowedCodes:       response.AllowedCodes,
		AllowedContentType: &response.AllowedContentType,
	}
}

func marshalToAssertion(assertion app.Assertion) *Assertion {
	if assertion.IsZero() {
		return nil
	}

	return &Assertion{
		With:   AssertionMethod(assertion.Method),
		Assert: marshalToAsserts(assertion.Asserts),
	}
}

func marshalToAsserts(asserts []app.Assert) []Assert {
	res := make([]Assert, 0, len(asserts))

	for _, a := range asserts {
		res = append(res, Assert{
			Actual:   a.Actual,
			Expected: fmt.Sprintf("%v", a.Expected),
		})
	}

	return res
}
