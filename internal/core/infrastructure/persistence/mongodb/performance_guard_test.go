package mongodb_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/infrastructure/persistence/mongodb"
)

type PerformanceGuardTestSuite struct {
	MongoSuite

	guard  *mongodb.PerformanceGuard
	perfID string
}

func (s *PerformanceGuardTestSuite) SetupTest() {
	s.guard = mongodb.NewPerformanceGuard(s.db)

	s.perfID = "2db44433-7142-4080-bada-844afccfedbf"

	s.insertPerformances(bson.M{
		"_id":     s.perfID,
		"started": false,
	})
}

func (s *PerformanceGuardTestSuite) TearDownTest() {
	_, err := s.db.
		Collection("performances").
		DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err)
}

func TestPerformanceGuard(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &PerformanceGuardTestSuite{})
}

func (s *PerformanceGuardTestSuite) TestAcquirePerformance() {
	ctx := context.Background()

	err := s.guard.AcquirePerformance(ctx, s.perfID)
	s.Require().NoError(err)

	err = s.guard.AcquirePerformance(ctx, s.perfID)
	s.Require().ErrorIs(err, performance.ErrAlreadyStarted)

	s.Require().True(s.getPerformanceStarted())
}

func (s *PerformanceGuardTestSuite) TestAcquirePerformanceConcurrently() {
	const triesCount = 100

	var wg sync.WaitGroup

	var acquiredCount int32

	wg.Add(triesCount)

	for i := 1; i <= triesCount; i++ {
		go func() {
			defer wg.Done()

			if err := s.guard.AcquirePerformance(
				context.Background(),
				s.perfID,
			); err == nil {
				atomic.AddInt32(&acquiredCount, 1)
			}
		}()
	}

	wg.Wait()

	s.Require().EqualValues(1, acquiredCount)
}

func (s *PerformanceGuardTestSuite) TestReleasePerformance() {
	err := s.guard.AcquirePerformance(context.Background(), s.perfID)
	s.Require().NoError(err)

	s.Require().True(s.getPerformanceStarted())

	err = s.guard.ReleasePerformance(context.Background(), s.perfID)
	s.Require().NoError(err)

	s.Require().False(s.getPerformanceStarted())
}

func (s *PerformanceGuardTestSuite) getPerformanceStarted() bool {
	var document struct {
		Started bool `bson:"started"`
	}

	err := s.db.Collection("performances").
		FindOne(context.Background(), bson.M{"_id": s.perfID}).
		Decode(&document)
	s.Require().NoError(err)

	return document.Started
}
