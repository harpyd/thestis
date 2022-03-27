package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	document, err := r.getTestCampaignDocument(ctx, bson.M{"_id": tcID})
	if err != nil {
		return nil, err
	}

	return document.unmarshalToTestCampaign(), nil
}

func (r *TestCampaignsRepository) FindTestCampaign(
	ctx context.Context,
	qry app.SpecificTestCampaignQuery,
) (app.SpecificTestCampaign, error) {
	document, err := r.getTestCampaignDocument(ctx, bson.M{
		"_id":     qry.TestCampaignID,
		"ownerId": qry.UserID,
	})
	if err != nil {
		return app.SpecificTestCampaign{}, err
	}

	return document.unmarshalToSpecificTestCampaign(), nil
}

func (r *TestCampaignsRepository) getTestCampaignDocument(
	ctx context.Context,
	filter bson.M,
) (testCampaignDocument, error) {
	var document testCampaignDocument
	if err := r.testCampaigns.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return testCampaignDocument{}, app.ErrTestCampaignNotFound
		}

		return testCampaignDocument{}, app.WrapWithDatabaseError(err)
	}

	return document, nil
}

func (r *TestCampaignsRepository) AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error {
	_, err := r.testCampaigns.InsertOne(ctx, marshalToTestCampaignDocument(tc))

	return app.WrapWithDatabaseError(err)
}

func (r *TestCampaignsRepository) UpdateTestCampaign(
	ctx context.Context,
	tcID string,
	updateFn app.TestCampaignUpdater,
) error {
	session, err := r.testCampaigns.Database().Client().StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(_ mongo.SessionContext) (interface{}, error) {
		var document testCampaignDocument
		if err := r.testCampaigns.FindOne(ctx, bson.M{"_id": tcID}).Decode(&document); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, app.ErrTestCampaignNotFound
			}

			return nil, app.WrapWithDatabaseError(err)
		}

		tc := document.unmarshalToTestCampaign()
		updatedTestCampaign, err := updateFn(ctx, tc)
		if err != nil {
			return nil, err
		}

		updatedDocument := marshalToTestCampaignDocument(updatedTestCampaign)

		replaceOpt := options.Replace().SetUpsert(true)
		filter := bson.M{"_id": updatedDocument.ID}
		if _, err := r.testCampaigns.ReplaceOne(ctx, filter, updatedDocument, replaceOpt); err != nil {
			return nil, app.WrapWithDatabaseError(err)
		}

		return nil, nil
	})

	return err
}

func (r *TestCampaignsRepository) RemoveAllTestCampaigns(ctx context.Context) error {
	_, err := r.testCampaigns.DeleteMany(ctx, bson.D{})

	return app.WrapWithDatabaseError(err)
}
