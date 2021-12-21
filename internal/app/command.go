package app

type (
	CreateTestCampaignCommand struct {
		OwnerID  string
		ViewName string
		Summary  string
	}

	LoadSpecificationCommand struct {
		TestCampaignID string
		Content        []byte
	}
)
