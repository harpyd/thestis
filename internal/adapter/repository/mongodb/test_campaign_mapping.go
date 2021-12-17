package mongodb

import (
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type (
	testCampaignDocument struct {
		ID                    string `bson:"_id,omitempty"`
		ViewName              string `bson:"viewName"`
		Summary               string `bson:"summary"`
		ActiveSpecificationID string `bson:"activeSpecificationId"`
	}
)

func marshalToTestCampaignDocument(tc *testcampaign.TestCampaign) testCampaignDocument {
	return testCampaignDocument{
		ID:                    tc.ID(),
		ViewName:              tc.ViewName(),
		Summary:               tc.Summary(),
		ActiveSpecificationID: tc.ActiveSpecificationID(),
	}
}

func (d testCampaignDocument) unmarshalToTestCampaign() *testcampaign.TestCampaign {
	return testcampaign.UnmarshalFromDatabase(
		d.ID,
		d.ViewName,
		d.Summary,
		d.ActiveSpecificationID,
	)
}

func (d testCampaignDocument) unmarshalToSpecificTestCampaign() app.SpecificTestCampaign {
	return app.SpecificTestCampaign{
		ID:                    d.ID,
		ViewName:              d.ViewName,
		Summary:               d.Summary,
		ActiveSpecificationID: d.ActiveSpecificationID,
	}
}
