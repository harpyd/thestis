package mongodb

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	lock "github.com/square/mongo-lock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type PerformancesRepository struct {
	performances *mongo.Collection
	locks        *lock.Client
}

const (
	performancesCollection     = "performances"
	performanceLocksCollection = "performanceLocks"
)

func NewPerformancesRepository(db *mongo.Database) *PerformancesRepository {
	r := &PerformancesRepository{
		performances: db.Collection(performancesCollection),
		locks:        lock.NewClient(db.Collection(performanceLocksCollection)),
	}

	if err := r.locks.CreateIndexes(context.Background()); err != nil {
		panic(err)
	}

	return r
}

func (r *PerformancesRepository) GetPerformance(ctx context.Context, perfID string) (*performance.Performance, error) {
	document, err := r.getPerformanceDocument(ctx, bson.M{"_id": perfID})
	if err != nil {
		return nil, err
	}

	return document.unmarshalToPerformance(), nil
}

func (r *PerformancesRepository) AcquirePerformance(_ context.Context, _ string) (*performance.Performance, error) {
	return nil, nil
}

func (r *PerformancesRepository) ReleasePerformance(_ context.Context, _ string) error {
	return nil
}

func (r *PerformancesRepository) getPerformanceDocument(
	ctx context.Context,
	filter bson.M,
) (performanceDocument, error) {
	var document performanceDocument
	if err := r.performances.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return performanceDocument{}, app.NewPerformanceNotFoundError(err)
		}

		return performanceDocument{}, app.NewDatabaseError(err)
	}

	return document, nil
}

func (r *PerformancesRepository) AddPerformance(ctx context.Context, perf *performance.Performance) error {
	_, err := r.performances.InsertOne(ctx, marshalToPerformanceDocument(perf))
	if mongo.IsDuplicateKeyError(err) {
		return app.NewAlreadyExistsError(err)
	}

	return app.NewDatabaseError(err)
}

func (r *PerformancesRepository) ExclusivelyDoWithPerformance(
	ctx context.Context,
	perf *performance.Performance,
	action app.PerformanceAction,
) error {
	if err := r.acquireLock(ctx, perf.ID()); err != nil {
		return err
	}

	go func() {
		defer r.releaseLock(context.Background(), perf.ID())

		action(ctx, perf)
	}()

	return nil
}

func (r *PerformancesRepository) acquireLock(ctx context.Context, perfID string) error {
	lockName := performanceLock(perfID)

	err := r.locks.XLock(ctx, lockName, lockName, lock.LockDetails{})
	if errors.Is(err, lock.ErrAlreadyLocked) {
		return performance.NewAlreadyStartedError()
	}

	return app.NewDatabaseError(err)
}

func (r *PerformancesRepository) releaseLock(ctx context.Context, perfID string) {
	if _, err := r.locks.Unlock(ctx, performanceLock(perfID)); err != nil {
		panic(app.NewDatabaseError(err))
	}
}

func performanceLock(perfID string) string {
	return fmt.Sprintf("performance#%s", perfID)
}

func (r *PerformancesRepository) RemoveAllPerformances(ctx context.Context) error {
	_, err := r.performances.DeleteMany(ctx, bson.D{})

	return app.NewDatabaseError(err)
}
