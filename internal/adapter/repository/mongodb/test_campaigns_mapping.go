package mongodb

import "github.com/harpyd/thestis/internal/domain/testcampaign"

type (
	testCampaignDocument struct {
		ID                    string `bson:"_id,omitempty"`
		ViewName              string `bson:"viewName"`
		ActiveSpecificationID string `bson:"activeSpecificationId"`
	}
)

func marshalToTestCampaignDocument(tc *testcampaign.TestCampaign) testCampaignDocument {
	return testCampaignDocument{
		ID:                    tc.ID(),
		ViewName:              tc.ViewName(),
		ActiveSpecificationID: tc.ActiveSpecificationID(),
	}
}

func (d testCampaignDocument) unmarshalToTestCampaign() *testcampaign.TestCampaign {
	return testcampaign.UnmarshalFromDatabase(
		d.ID,
		d.ViewName,
		d.ActiveSpecificationID,
	)
}
