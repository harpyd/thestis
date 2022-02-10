package app

type (
	CreateTestCampaignCommand struct {
		OwnerID  string
		ViewName string
		Summary  string
	}

	LoadSpecificationCommand struct {
		TestCampaignID string
		LoadedByID     string
		Content        []byte
	}

	StartPerformanceCommand struct {
		TestCampaignID string
		StartedByID    string
	}

	RestartPerformanceCommand struct {
		PerformanceID string
		StartedByID   string
	}

	CancelPerformanceCommand struct {
		PerformanceID string
		CanceledByID  string
	}
)
