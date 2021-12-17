package app

type (
	SpecificTestCampaignQuery struct {
		TestCampaignID string
	}
)

type (
	SpecificTestCampaign struct {
		ID                    string
		ViewName              string
		Summary               string
		ActiveSpecificationID string
	}
)
