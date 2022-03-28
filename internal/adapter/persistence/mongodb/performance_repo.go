package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type PerformanceRepository struct {
	performances *mongo.Collection
}

const performancesCollection = "performances"

func NewPerformanceRepository(db *mongo.Database) *PerformanceRepository {
	r := &PerformanceRepository{
		performances: db.Collection(performancesCollection),
	}

	return r
}

func (r *PerformanceRepository) GetPerformance(
	ctx context.Context,
	perfID string,
	specGetter app.SpecificationGetter,
	opts ...performance.Option,
) (*performance.Performance, error) {
	document, err := r.getPerformanceDocument(ctx, bson.M{"_id": perfID})
	if err != nil {
		return nil, err
	}

	spec, err := specGetter.GetSpecification(ctx, document.SpecificationID)
	if err != nil {
		return nil, err
	}

	return document.unmarshalToPerformance(spec, opts), nil
}

func (r *PerformanceRepository) getPerformanceDocument(
	ctx context.Context,
	filter bson.M,
) (performanceDocument, error) {
	var document performanceDocument
	if err := r.performances.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return performanceDocument{}, app.ErrPerformanceNotFound
		}

		return performanceDocument{}, app.WrapWithDatabaseError(err)
	}

	return document, nil
}

func (r *PerformanceRepository) AddPerformance(ctx context.Context, perf *performance.Performance) error {
	_, err := r.performances.InsertOne(ctx, marshalToPerformanceDocument(perf))

	return app.WrapWithDatabaseError(err)
}

func (r *PerformanceRepository) RemoveAllPerformances(ctx context.Context) error {
	_, err := r.performances.DeleteMany(ctx, bson.D{})

	return app.WrapWithDatabaseError(err)
}