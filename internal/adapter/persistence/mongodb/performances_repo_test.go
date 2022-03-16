package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type PerformancesRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repo *mongodb.PerformancesRepository
}

func (s *PerformancesRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewPerformancesRepository(s.db)
}

func (s *PerformancesRepositoryTestSuite) TearDownTest() {
	err := s.repo.RemoveAllPerformances(context.Background())
	s.Require().NoError(err)
}

func TestPerformancesRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &PerformancesRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *PerformancesRepositoryTestSuite) TestAddPerformance() {
	testCases := []struct {
		Name                   string
		Before                 func()
		GivenPerformanceParams performance.Params
		ShouldBeErr            bool
		IsErr                  func(err error) bool
	}{
		{
			Name: "failed_adding_one_performance_twice",
			Before: func() {
				perf := performance.Unmarshal(performance.Params{
					ID: "a4a2906d-4df5-42f1-8832-77a33cba4d7f",
				})

				s.addPerformances(perf)
			},
			GivenPerformanceParams: performance.Params{
				ID:      "a4a2906d-4df5-42f1-8832-77a33cba4d7f",
				OwnerID: "05bf69e9-d7b5-4e7a-8fab-24227dca033a",
				Specification: specification.NewBuilder().
					WithID("1a63b6ea-df5f-4a68-bf04-c2e30044f2ef").
					ErrlessBuild(),
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return app.IsAlreadyExistsError(err) && mongo.IsDuplicateKeyError(err)
			},
		},
		{
			Name: "successful_adding",
			GivenPerformanceParams: performance.Params{
				ID:      "3ce098e1-81ae-4610-8372-2f635b1b6a0c",
				OwnerID: "3614a95c-c278-4687-84e2-97b95b11d399",
				Specification: specification.NewBuilder().
					WithID("spec").
					ErrlessBuild(),
				Started: true,
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			if c.Before != nil {
				c.Before()
			}

			perf := performance.Unmarshal(c.GivenPerformanceParams)

			err := s.repo.AddPerformance(context.Background(), perf)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedPerf, err := s.repo.GetPerformance(
				context.Background(),
				perf.ID(),
				app.AvailableSpecification(c.GivenPerformanceParams.Specification),
			)
			s.Require().NoError(err)

			s.Require().Equal(perf, persistedPerf)
		})
	}
}

func (s *PerformancesRepositoryTestSuite) addPerformances(perfs ...*performance.Performance) {
	s.T().Helper()

	ctx := context.Background()

	for _, perf := range perfs {
		err := s.repo.AddPerformance(ctx, perf)
		s.Require().NoError(err)
	}
}
