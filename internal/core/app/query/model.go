package query

import "time"

type SpecificTestCampaignModel struct {
	ID        string
	ViewName  string
	Summary   string
	CreatedAt time.Time
}

type (
	SpecificSpecificationModel struct {
		ID             string
		TestCampaignID string
		LoadedAt       time.Time
		Author         string
		Title          string
		Description    string
		Stories        []StoryModel
	}

	StoryModel struct {
		Slug        string
		Description string
		AsA         string
		InOrderTo   string
		WantTo      string
		Scenarios   []ScenarioModel
	}

	ScenarioModel struct {
		Slug        string
		Description string
		Theses      []ThesisModel
	}

	ThesisModel struct {
		Slug      string
		After     []string
		Statement StatementModel
		HTTP      HTTPModel
		Assertion AssertionModel
	}

	StatementModel struct {
		Stage    string
		Behavior string
	}

	HTTPModel struct {
		Request  HTTPRequestModel
		Response HTTPResponseModel
	}

	HTTPRequestModel struct {
		Method      string
		URL         string
		ContentType string
		Body        map[string]interface{}
	}

	HTTPResponseModel struct {
		AllowedCodes       []int
		AllowedContentType string
	}

	AssertionModel struct {
		Method  string
		Asserts []AssertModel
	}

	AssertModel struct {
		Actual   string
		Expected interface{}
	}
)

func (h HTTPModel) IsZero() bool {
	return h.Request.IsZero() && h.Response.IsZero()
}

func (r HTTPRequestModel) IsZero() bool {
	return r.Method == "" &&
		r.URL == "" &&
		r.ContentType == "" &&
		len(r.Body) == 0
}

func (r HTTPResponseModel) IsZero() bool {
	return r.AllowedContentType == "" && len(r.AllowedCodes) == 0
}

func (a AssertionModel) IsZero() bool {
	return a.Method == "" && len(a.Asserts) == 0
}

type (
	SpecificPerformanceModel struct {
		ID              string
		SpecificationID string
		StartedAt       time.Time
		Flows           []FlowModel
	}

	FlowModel struct {
		StartedAt    time.Time
		OverallState string
		Statuses     []StatusModel
	}

	StatusModel struct {
		Slug           ScenarioSlugModel
		State          string
		ThesisStatuses []ThesisStatusModel
	}

	ScenarioSlugModel struct {
		Story    string
		Scenario string
	}

	ThesisStatusModel struct {
		ThesisSlug   string
		State        string
		OccurredErrs []string
	}
)
