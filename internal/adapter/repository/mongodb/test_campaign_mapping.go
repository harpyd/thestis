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
		UserID                string `bson:"userId"`
	}
)

func marshalToTestCampaignDocument(tc *testcampaign.TestCampaign) testCampaignDocument {
	return testCampaignDocument{
		ID:                    tc.ID(),
		ViewName:              tc.ViewName(),
		Summary:               tc.Summary(),
		ActiveSpecificationID: tc.ActiveSpecificationID(),
		UserID:                tc.UserID(),
	}
}

func (d testCampaignDocument) unmarshalToTestCampaign() *testcampaign.TestCampaign {
	tc, _ := testcampaign.New(testcampaign.Params{
		ID:                    d.ID,
		ViewName:              d.ViewName,
		Summary:               d.Summary,
		ActiveSpecificationID: d.ActiveSpecificationID,
		UserID:                d.UserID,
	})

	return tc
}

func (d testCampaignDocument) unmarshalToSpecificTestCampaign() app.SpecificTestCampaign {
	return app.SpecificTestCampaign{
		ID:                    d.ID,
		ViewName:              d.ViewName,
		Summary:               d.Summary,
		ActiveSpecificationID: d.ActiveSpecificationID,
		UserID:                d.UserID,
	}
}
