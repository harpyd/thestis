package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
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
	return nil
}

func (g *PerformanceGuard) ReleasePerformance(ctx context.Context, perfID string) error {
	return nil
}
