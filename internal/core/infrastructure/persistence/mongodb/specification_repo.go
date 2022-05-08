package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type SpecificationRepository struct {
	specifications *mongo.Collection
}

const specificationsCollection = "specifications"

func NewSpecificationRepository(db *mongo.Database) *SpecificationRepository {
	col := db.Collection(specificationsCollection)

	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}

	return &SpecificationRepository{
		specifications: col,
	}
}

func (r *SpecificationRepository) GetSpecification(
	ctx context.Context,
	specID string,
) (*specification.Specification, error) {
	document, err := r.getSpecificationDocument(ctx, bson.M{"id": specID}, nil)
	if err != nil {
		return nil, err
	}

	return newSpecification(document), nil
}

func (r *SpecificationRepository) GetActiveSpecificationByTestCampaignID(
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

	return newSpecification(document), nil
}

func (r *SpecificationRepository) FindSpecification(
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

	return newAppSpecificSpecification(document), nil
}

func (r *SpecificationRepository) getSpecificationDocument(
	ctx context.Context,
	filter bson.M,
	sort bson.M,
) (specificationDocument, error) {
	opt := options.FindOne().SetSort(sort)

	var document specificationDocument
	if err := r.specifications.FindOne(ctx, filter, opt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return specificationDocument{}, service.ErrSpecificationNotFound
		}

		return specificationDocument{}, service.WrapWithDatabaseError(err)
	}

	return document, nil
}

func (r *SpecificationRepository) AddSpecification(ctx context.Context, spec *specification.Specification) error {
	_, err := r.specifications.InsertOne(ctx, newSpecificationDocument(spec))

	return service.WrapWithDatabaseError(err)
}
