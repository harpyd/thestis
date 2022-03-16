package mongodb_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type PerformanceGuardTestSuite struct {
	suite.Suite
	MongoTestFixtures

	guard *mongodb.PerformanceGuard
	repo  *mongodb.PerformancesRepository

	perfID string
}

func (s *PerformanceGuardTestSuite) SetupTest() {
	s.guard = mongodb.NewPerformanceGuard(s.db)
	s.repo = mongodb.NewPerformancesRepository(s.db)

	s.perfID = "2db44433-7142-4080-bada-844afccfedbf"

	err := s.repo.AddPerformance(context.Background(), performance.Unmarshal(performance.Params{
		ID:      s.perfID,
		OwnerID: "b4a232b3-d853-4022-9b4d-f2aefb24d82c",
	}))
	s.Require().NoError(err)
}

func (s *PerformanceGuardTestSuite) TearDownTest() {
	err := s.repo.RemoveAllPerformances(context.Background())
	s.Require().NoError(err)
}

func TestPerformanceGuard(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &PerformanceGuardTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *PerformanceGuardTestSuite) TestAcquirePerformance() {
	err := s.guard.AcquirePerformance(context.Background(), s.perfID)
	s.Require().NoError(err)

	err = s.guard.AcquirePerformance(context.Background(), s.perfID)
	s.Require().True(performance.IsAlreadyStartedError(err))

	persistedPerf, err := s.repo.GetPerformance(
		context.Background(),
		s.perfID,
		app.DontGetSpecification(),
	)
	s.Require().NoError(err)

	s.Require().True(persistedPerf.Started())
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

	persistedPerf, err := s.repo.GetPerformance(
		context.Background(),
		s.perfID,
		app.DontGetSpecification(),
	)
	s.Require().NoError(err)

	s.Require().True(persistedPerf.Started())

	err = s.guard.ReleasePerformance(context.Background(), s.perfID)
	s.Require().NoError(err)

	persistedPerf, err = s.repo.GetPerformance(
		context.Background(),
		s.perfID,
		app.DontGetSpecification(),
	)
	s.Require().NoError(err)

	s.Require().False(persistedPerf.Started())
}
