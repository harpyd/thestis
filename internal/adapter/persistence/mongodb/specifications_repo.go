package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationsRepository struct {
	specifications *mongo.Collection
}

const specificationsCollection = "specifications"

func NewSpecificationsRepository(db *mongo.Database) *SpecificationsRepository {
	col := db.Collection(specificationsCollection)

	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}

	return &SpecificationsRepository{
		specifications: col,
	}
}

func (r *SpecificationsRepository) GetSpecification(
	ctx context.Context,
	specID string,
) (*specification.Specification, error) {
	document, err := r.getSpecificationDocument(ctx, bson.M{"id": specID}, nil)
	if err != nil {
		return nil, err
	}

	return document.unmarshalToSpecification(), nil
}

func (r *SpecificationsRepository) GetActiveSpecificationByTestCampaignID(
	ctx context.Context,
	testCampaignID string,
) (*specification.Specification, error) {
	var (
		filter = bson.M{"testCampaignId": testCampaignID}
		sort   = bson.M{"_id": -1}
	)

	document, err := r.getSpecificationDocument(ctx, filter, sort)
	if err != nil {
		return nil, err
	}

	return document.unmarshalToSpecification(), nil
}

func (r *SpecificationsRepository) FindSpecification(
	ctx context.Context,
	qry app.SpecificSpecificationQuery,
) (app.SpecificSpecification, error) {
	filter := bson.M{
		"id":      qry.SpecificationID,
		"ownerId": qry.UserID,
	}

	document, err := r.getSpecificationDocument(ctx, filter, nil)
	if err != nil {
		return app.SpecificSpecification{}, err
	}

	return document.unmarshalToSpecificSpecification(), nil
}

func (r *SpecificationsRepository) getSpecificationDocument(
	ctx context.Context,
	filter bson.M,
	sort bson.M,
) (specificationDocument, error) {
	opt := options.FindOne().SetSort(sort)

	var document specificationDocument
	if err := r.specifications.FindOne(ctx, filter, opt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return specificationDocument{}, app.NewSpecificationNotFoundError(err)
		}

		return specificationDocument{}, app.NewDatabaseError(err)
	}

	return document, nil
}

func (r *SpecificationsRepository) AddSpecification(ctx context.Context, spec *specification.Specification) error {
	_, err := r.specifications.InsertOne(ctx, marshalToSpecificationDocument(spec))
	if mongo.IsDuplicateKeyError(err) {
		return app.NewAlreadyExistsError(err)
	}

	return app.NewDatabaseError(err)
}

func (r *SpecificationsRepository) RemoveAllSpecifications(ctx context.Context) error {
	_, err := r.specifications.DeleteMany(ctx, bson.D{})

	return app.NewDatabaseError(err)
}
