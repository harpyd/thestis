package app

type (
	SpecificTestCampaignQuery struct {
		TestCampaignID string
	}

	SpecificSpecificationQuery struct {
		SpecificationID string
	}
)

type (
	// SpecificTestCampaign is application layer
	// representation of testcampaign.TestCampaign.
	SpecificTestCampaign struct {
		ID                    string
		ViewName              string
		Summary               string
		ActiveSpecificationID string
	}

	// SpecificSpecification is application layer
	// representation of specification.Specification.
	SpecificSpecification struct {
		ID          string
		Author      string
		Title       string
		Description string
		Stories     []Story
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
		Keyword  string
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
