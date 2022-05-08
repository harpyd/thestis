package query

import "time"

type SpecificTestCampaignView struct {
	ID        string
	ViewName  string
	Summary   string
	CreatedAt time.Time
}

type (
	SpecificSpecificationView struct {
		ID             string
		TestCampaignID string
		LoadedAt       time.Time
		Author         string
		Title          string
		Description    string
		Stories        []StoryView
	}

	StoryView struct {
		Slug        string
		Description string
		AsA         string
		InOrderTo   string
		WantTo      string
		Scenarios   []ScenarioView
	}

	ScenarioView struct {
		Slug        string
		Description string
		Theses      []ThesisView
	}

	ThesisView struct {
		Slug      string
		After     []string
		Statement StatementView
		HTTP      HTTPView
		Assertion AssertionView
	}

	StatementView struct {
		Stage    string
		Behavior string
	}

	HTTPView struct {
		Request  HTTPRequestView
		Response HTTPResponseView
	}

	HTTPRequestView struct {
		Method      string
		URL         string
		ContentType string
		Body        map[string]interface{}
	}

	HTTPResponseView struct {
		AllowedCodes       []int
		AllowedContentType string
	}

	AssertionView struct {
		Method  string
		Asserts []AssertView
	}

	AssertView struct {
		Actual   string
		Expected interface{}
	}
)

func (h HTTPView) IsZero() bool {
	return h.Request.IsZero() && h.Response.IsZero()
}

func (r HTTPRequestView) IsZero() bool {
	return r.Method == "" &&
		r.URL == "" &&
		r.ContentType == "" &&
		len(r.Body) == 0
}

func (r HTTPResponseView) IsZero() bool {
	return r.AllowedContentType == "" && len(r.AllowedCodes) == 0
}

func (a AssertionView) IsZero() bool {
	return a.Method == "" && len(a.Asserts) == 0
}

type (
	SpecificPerformanceView struct {
		ID              string
		SpecificationID string
		StartedAt       time.Time
		Flows           []FlowView
	}

	FlowView struct {
		StartedAt    time.Time
		OverallState string
		Statuses     []StatusView
	}

	StatusView struct {
		Slug           ScenarioSlugView
		State          string
		ThesisStatuses []ThesisStatusView
	}

	ScenarioSlugView struct {
		Story    string
		Scenario string
	}

	ThesisStatusView struct {
		ThesisSlug   string
		State        string
		OccurredErrs []string
	}
)
