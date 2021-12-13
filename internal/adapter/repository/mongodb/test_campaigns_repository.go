package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type TestCampaignsRepository struct {
	testCampaigns *mongo.Collection
}

const testCampaignsCollection = "testCampaigns"

func NewTestCampaignsRepository(db *mongo.Database) *TestCampaignsRepository {
	return &TestCampaignsRepository{
		testCampaigns: db.Collection(testCampaignsCollection),
	}
}

func (r *TestCampaignsRepository) GetTestCampaign(
	ctx context.Context,
	tcID string,
) (*testcampaign.TestCampaign, error) {
	document, err := r.getTestCampaignDocument(ctx, tcID)
	if err != nil {
		return nil, err
	}

	return document.unmarshalToTestCampaign(), nil
}

func (r *TestCampaignsRepository) getTestCampaignDocument(
	ctx context.Context,
	tcID string,
) (testCampaignDocument, error) {
	var document testCampaignDocument
	if err := r.testCampaigns.FindOne(ctx, bson.M{"_id": tcID}).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return testCampaignDocument{}, app.NewTestCampaignNotFoundError(err)
		}

		return testCampaignDocument{}, app.NewDatabaseError(err)
	}

	return document, nil
}

func (r *TestCampaignsRepository) AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error {
	_, err := r.testCampaigns.InsertOne(ctx, marshalToTestCampaignDocument(tc))

	return app.NewDatabaseError(err)
}

func (r *TestCampaignsRepository) RemoveAllTestCampaigns(ctx context.Context) error {
	_, err := r.testCampaigns.DeleteMany(ctx, bson.D{})

	return app.NewDatabaseError(err)
}
