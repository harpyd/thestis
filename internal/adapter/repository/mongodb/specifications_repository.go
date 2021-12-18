package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationsRepository struct {
	specifications *mongo.Collection
}

const specificationsCollection = "specifications"

func NewSpecificationsRepository(db *mongo.Database) *SpecificationsRepository {
	return &SpecificationsRepository{
		specifications: db.Collection(specificationsCollection),
	}
}

func (r *SpecificationsRepository) GetSpecification(
	ctx context.Context,
	specID string,
) (*specification.Specification, error) {
	document, err := r.getSpecificationDocument(ctx, specID)
	if err != nil {
		return nil, err
	}

	return document.unmarshalToSpecification(), nil
}

func (r *SpecificationsRepository) FindSpecification(
	ctx context.Context,
	qry app.SpecificSpecificationQuery,
) (app.SpecificSpecification, error) {
	document, err := r.getSpecificationDocument(ctx, qry.SpecificationID)
	if err != nil {
		return app.SpecificSpecification{}, err
	}

	return document.unmarshalToSpecificSpecification(), nil
}

func (r *SpecificationsRepository) getSpecificationDocument(
	ctx context.Context,
	specID string,
) (specificationDocument, error) {
	var document specificationDocument
	if err := r.specifications.FindOne(ctx, bson.M{"_id": specID}).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return specificationDocument{}, app.NewSpecificationNotFoundError(err)
		}

		return specificationDocument{}, app.NewDatabaseError(err)
	}

	return document, nil
}

func (r *SpecificationsRepository) AddSpecification(ctx context.Context, spec *specification.Specification) error {
	_, err := r.specifications.InsertOne(ctx, marshalToSpecificationDocument(spec))

	return app.NewDatabaseError(err)
}

func (r *SpecificationsRepository) RemoveAllSpecifications(ctx context.Context) error {
	_, err := r.specifications.DeleteMany(ctx, bson.D{})

	return app.NewDatabaseError(err)
}
