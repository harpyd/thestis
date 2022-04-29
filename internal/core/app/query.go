package app

import "time"

type (
	SpecificTestCampaignQuery struct {
		TestCampaignID string
		UserID         string
	}

	SpecificSpecificationQuery struct {
		SpecificationID string
		UserID          string
	}

	SpecificPerformanceQuery struct {
		PerformanceID string
		UserID        string
	}
)

type (
	SpecificPerformance struct {
		ID              string
		SpecificationID string
		StartedAt       time.Time
		Flows           []Flow
	}

	Flow struct {
		StartedAt    time.Time
		OverallState string
		Statuses     []Status
	}

	Status struct {
		Slug           ScenarioSlug
		State          string
		ThesisStatuses []ThesisStatus
	}

	ScenarioSlug struct {
		Story    string
		Scenario string
	}

	ThesisStatus struct {
		ThesisSlug   string
		State        string
		OccurredErrs []string
	}

	// SpecificTestCampaign is most detailed application layer
	// representation of testcampaign.TestCampaign.
	SpecificTestCampaign struct {
		ID        string
		ViewName  string
		Summary   string
		CreatedAt time.Time
	}

	// SpecificSpecification is most detailed application layer
	// representation of specification.Specification.
	SpecificSpecification struct {
		ID             string
		TestCampaignID string
		LoadedAt       time.Time
		Author         string
		Title          string
		Description    string
		Stories        []Story
	}

	Story struct {
		Slug        string
		Description string
		AsA         string
		InOrderTo   string
		WantTo      string
		Scenarios   []Scenario
	}

	Scenario struct {
		Slug        string
		Description string
		Theses      []Thesis
	}

	Thesis struct {
		Slug      string
		After     []string
		Statement Statement
		HTTP      HTTP
		Assertion Assertion
	}

	Statement struct {
		Stage    string
		Behavior string
	}

	HTTP struct {
		Request  HTTPRequest
		Response HTTPResponse
	}

	HTTPRequest struct {
		Method      string
		URL         string
		ContentType string
		Body        map[string]interface{}
	}

	HTTPResponse struct {
		AllowedCodes       []int
		AllowedContentType string
	}

	Assertion struct {
		Method  string
		Asserts []Assert
	}

	Assert struct {
		Actual   string
		Expected interface{}
	}
)

func (h HTTP) IsZero() bool {
	return h.Request.IsZero() && h.Response.IsZero()
}

func (r HTTPRequest) IsZero() bool {
	return r.Method == "" &&
		r.URL == "" &&
		r.ContentType == "" &&
		len(r.Body) == 0
}

func (r HTTPResponse) IsZero() bool {
	return r.AllowedContentType == "" && len(r.AllowedCodes) == 0
}

func (a Assertion) IsZero() bool {
	return a.Method == "" && len(a.Asserts) == 0
}
