package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/harpyd/thestis/internal/core/app/query"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/testcampaign"
)

type TestCampaignRepository struct {
	testCampaigns *mongo.Collection
}

const testCampaignCollection = "testCampaigns"

func NewTestCampaignRepository(db *mongo.Database) *TestCampaignRepository {
	return &TestCampaignRepository{
		testCampaigns: db.Collection(testCampaignCollection),
	}
}

func (r *TestCampaignRepository) GetTestCampaign(
	ctx context.Context,
	tcID string,
) (*testcampaign.TestCampaign, error) {
	document, err := r.getTestCampaignDocument(ctx, bson.M{"_id": tcID})
	if err != nil {
		return nil, err
	}

	return newTestCampaign(document), nil
}

func (r *TestCampaignRepository) FindTestCampaign(
	ctx context.Context,
	qry query.TestCampaign,
) (query.TestCampaignModel, error) {
	document, err := r.getTestCampaignDocument(ctx, bson.M{
		"_id":     qry.TestCampaignID,
		"ownerId": qry.UserID,
	})
	if err != nil {
		return query.TestCampaignModel{}, err
	}

	return newSpecificTestCampaignView(document), nil
}

func (r *TestCampaignRepository) getTestCampaignDocument(
	ctx context.Context,
	filter bson.M,
) (testCampaignDocument, error) {
	var document testCampaignDocument
	if err := r.testCampaigns.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return testCampaignDocument{}, service.ErrTestCampaignNotFound
		}

		return testCampaignDocument{}, service.WrapWithDatabaseError(err)
	}

	return document, nil
}

func (r *TestCampaignRepository) AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error {
	_, err := r.testCampaigns.InsertOne(ctx, newTestCampaignDocument(tc))

	return service.WrapWithDatabaseError(err)
}

func (r *TestCampaignRepository) UpdateTestCampaign(
	ctx context.Context,
	tcID string,
	updater service.TestCampaignUpdater,
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
				return nil, service.ErrTestCampaignNotFound
			}

			return nil, service.WrapWithDatabaseError(err)
		}

		tc := newTestCampaign(document)
		updatedTestCampaign, err := updater(ctx, tc)
		if err != nil {
			return nil, err
		}

		updatedDocument := newTestCampaignDocument(updatedTestCampaign)

		replaceOpt := options.Replace().SetUpsert(true)
		filter := bson.M{"_id": updatedDocument.ID}
		if _, err := r.testCampaigns.ReplaceOne(ctx, filter, updatedDocument, replaceOpt); err != nil {
			return nil, service.WrapWithDatabaseError(err)
		}

		return nil, nil
	})

	return err
}
