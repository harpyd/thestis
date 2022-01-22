package mongodb

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	lock "github.com/square/mongo-lock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/multierr"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type PerformancesRepository struct {
	performances *mongo.Collection
	locks        *lock.Client
}

const performancesCollection = "performances"

func NewPerformancesRepository(db *mongo.Database) *PerformancesRepository {
	performances := db.Collection(performancesCollection)

	r := &PerformancesRepository{
		performances: performances,
		locks:        lock.NewClient(performances),
	}

	if err := r.locks.CreateIndexes(context.Background()); err != nil {
		panic(err)
	}

	return r
}

func (r *PerformancesRepository) GetPerformance(ctx context.Context, perfID string) (*performance.Performance, error) {
	document, err := r.getPerformanceDocument(ctx, perfID, "")
	if err != nil {
		return nil, err
	}

	return document.unmarshalToPerformance(), nil
}

func (r *PerformancesRepository) getPerformanceDocument(
	ctx context.Context,
	perfID, userID string,
) (performanceDocument, error) {
	filter := makePerformanceFilter(perfID, userID)

	var document performanceDocument
	if err := r.performances.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return performanceDocument{}, app.NewPerformanceNotFoundError(err)
		}

		return performanceDocument{}, app.NewDatabaseError(err)
	}

	return document, nil
}

func makePerformanceFilter(specID, userID string) bson.M {
	filter := bson.M{"_id": specID}
	if userID != "" {
		filter["ownerId"] = userID
	}

	return filter
}

func (r *PerformancesRepository) AddPerformance(ctx context.Context, perf *performance.Performance) error {
	_, err := r.performances.InsertOne(ctx, marshalToPerformanceDocument(perf))

	return app.NewDatabaseError(err)
}

func (r *PerformancesRepository) DoWithPerformance(
	ctx context.Context,
	perfID string,
	action app.PerformanceAction,
) (err error) {
	if err := r.acquireLock(ctx, perfID); err != nil {
		return err
	}

	defer func() {
		err = multierr.Append(err, r.releaseLock(ctx, perfID))
	}()

	perf, err := r.GetPerformance(ctx, perfID)
	if err != nil {
		return err
	}

	action(ctx, perf)

	return
}

func (r *PerformancesRepository) RemoveAllPerformances(ctx context.Context) error {
	_, err := r.performances.DeleteMany(ctx, bson.D{})

	return app.NewDatabaseError(err)
}

func (r *PerformancesRepository) acquireLock(ctx context.Context, perfID string) error {
	lockName := performanceLock(perfID)

	err := r.locks.XLock(ctx, lockName, lockName, lock.LockDetails{})
	if errors.Is(err, lock.ErrAlreadyLocked) {
		return performance.NewAlreadyStartedError()
	}

	return app.NewDatabaseError(err)
}

func (r *PerformancesRepository) releaseLock(ctx context.Context, perfID string) error {
	_, err := r.locks.Unlock(ctx, performanceLock(perfID))

	return app.NewDatabaseError(err)
}

func performanceLock(perfID string) string {
	return fmt.Sprintf("performance#%s", perfID)
}
