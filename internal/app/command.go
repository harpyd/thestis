package app

type (
	CreateTestCampaignCommand struct {
		ViewName string
		Summary  string
	}

	LoadSpecificationCommand struct {
		TestCampaignID string
		Content        []byte
	}
)
