package testcampaign

type TestCampaign struct {
	id                    string
	viewName              string
	activeSpecificationID string
}

func New(id string, viewName string) *TestCampaign {
	return &TestCampaign{
		id:       id,
		viewName: viewName,
	}
}

func UnmarshalFromDatabase(
	id string,
	viewName string,
	activeSpecificationID string,
) *TestCampaign {
	return &TestCampaign{
		id:                    id,
		viewName:              viewName,
		activeSpecificationID: activeSpecificationID,
	}
}

func (tc *TestCampaign) ID() string {
	return tc.id
}

func (tc *TestCampaign) ViewName() string {
	return tc.viewName
}

func (tc *TestCampaign) ActiveSpecificationID() string {
	return tc.activeSpecificationID
}

func (tc *TestCampaign) SetActiveSpecificationID(specificationID string) {
	tc.activeSpecificationID = specificationID
}
