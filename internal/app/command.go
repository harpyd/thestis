package app

type (
	CreateTestCampaignCommand struct {
		UserID   string
		ViewName string
		Summary  string
	}

	LoadSpecificationCommand struct {
		TestCampaignID string
		Content        []byte
	}
)
