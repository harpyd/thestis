package mongodb

import (
	"time"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type (
	testCampaignDocument struct {
		ID        string    `bson:"_id,omitempty"`
		ViewName  string    `bson:"viewName"`
		Summary   string    `bson:"summary"`
		OwnerID   string    `bson:"ownerId"`
		CreatedAt time.Time `bson:"createdAt"`
	}
)

func marshalToTestCampaignDocument(tc *testcampaign.TestCampaign) testCampaignDocument {
	return testCampaignDocument{
		ID:        tc.ID(),
		ViewName:  tc.ViewName(),
		Summary:   tc.Summary(),
		OwnerID:   tc.OwnerID(),
		CreatedAt: tc.CreatedAt(),
	}
}

func (d testCampaignDocument) unmarshalToTestCampaign() *testcampaign.TestCampaign {
	tc, _ := testcampaign.New(testcampaign.Params{
		ID:        d.ID,
		ViewName:  d.ViewName,
		Summary:   d.Summary,
		OwnerID:   d.OwnerID,
		CreatedAt: d.CreatedAt,
	})

	return tc
}

func (d testCampaignDocument) unmarshalToSpecificTestCampaign() app.SpecificTestCampaign {
	return app.SpecificTestCampaign{
		ID:        d.ID,
		ViewName:  d.ViewName,
		Summary:   d.Summary,
		CreatedAt: d.CreatedAt,
	}
}
