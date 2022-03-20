package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type PerformanceGuard struct {
	performances *mongo.Collection
}

func NewPerformanceGuard(db *mongo.Database) *PerformanceGuard {
	return &PerformanceGuard{
		performances: db.Collection(performancesCollection),
	}
}

func (g *PerformanceGuard) AcquirePerformance(ctx context.Context, perfID string) error {
	var document performanceDocument

	var (
		filter = bson.M{"_id": perfID}
		update = bson.M{"$set": bson.M{"started": true}}
		opt    = options.FindOneAndUpdate().SetProjection(bson.M{"_id": 0, "started": 1})
	)

	if err := g.performances.FindOneAndUpdate(ctx, filter, update, opt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app.NewPerformanceNotFoundError(err)
		}

		return app.NewDatabaseError(err)
	}

	if document.Started {
		return performance.ErrAlreadyStarted
	}

	return nil
}

func (g *PerformanceGuard) ReleasePerformance(ctx context.Context, perfID string) error {
	update := bson.M{"$set": bson.M{"started": false}}

	_, err := g.performances.UpdateByID(ctx, perfID, update)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return app.NewPerformanceNotFoundError(err)
	}

	return app.NewDatabaseError(err)
}
