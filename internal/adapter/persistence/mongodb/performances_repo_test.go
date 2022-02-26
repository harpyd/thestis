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
		Name        string
		Before      func()
		Performance *performance.Performance
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "failed_adding_one_performance_twice",
			Before: func() {
				perf := performance.Unmarshal(performance.Params{
					SpecificationID: "1a63b6ea-df5f-4a68-bf04-c2e30044f2ef",
				}, performance.WithID("a4a2906d-4df5-42f1-8832-77a33cba4d7f"))

				s.addPerformances(perf)
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID:         "05bf69e9-d7b5-4e7a-8fab-24227dca033a",
				SpecificationID: "1a63b6ea-df5f-4a68-bf04-c2e30044f2ef",
				Actions:         s.actions(),
			}, performance.WithID("a4a2906d-4df5-42f1-8832-77a33cba4d7f")),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return app.IsAlreadyExistsError(err) && mongo.IsDuplicateKeyError(err)
			},
		},
		{
			Name: "successful_adding",
			Performance: performance.Unmarshal(performance.Params{
				OwnerID:         "3614a95c-c278-4687-84e2-97b95b11d399",
				SpecificationID: "4e4465b0-a312-4f86-9051-a3ae72965215",
				Actions:         s.actions(),
			}, performance.WithID("3ce098e1-81ae-4610-8372-2f635b1b6a0c")),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			if c.Before != nil {
				c.Before()
			}

			err := s.repo.AddPerformance(context.Background(), c.Performance)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedPerf := s.getPerformance(c.Performance.ID())
			s.requirePerformancesEqual(c.Performance, persistedPerf)
		})
	}
}

func (s *PerformancesRepositoryTestSuite) actions() []performance.Action {
	s.T().Helper()

	thesis, err := specification.NewThesisBuilder().
		WithStatement("given", "test").
		WithAssertion(func(b *specification.AssertionBuilder) {
			b.WithMethod("jsonpath")
			b.WithAssert("some", "value")
		}).
		Build(specification.NewThesisSlug("story", "scenario", "to"))
	s.Require().NoError(err)

	return []performance.Action{
		performance.NewAction(
			"stage.given",
			"story.scenario.to",
			thesis,
			performance.AssertionPerformer,
		),
	}
}

func (s *PerformancesRepositoryTestSuite) getPerformance(perfID string) *performance.Performance {
	s.T().Helper()

	perf, err := s.repo.GetPerformance(context.Background(), perfID)
	s.Require().NoError(err)

	return perf
}

func (s *PerformancesRepositoryTestSuite) addPerformances(perfs ...*performance.Performance) {
	s.T().Helper()

	ctx := context.Background()

	for _, perf := range perfs {
		err := s.repo.AddPerformance(ctx, perf)
		s.Require().NoError(err)
	}
}

func (s *PerformancesRepositoryTestSuite) requirePerformancesEqual(expected, actual *performance.Performance) {
	s.T().Helper()

	s.Require().Equal(expected.ID(), actual.ID())
	s.Require().Equal(expected.OwnerID(), actual.OwnerID())
	s.Require().Equal(expected.SpecificationID(), actual.SpecificationID())
	s.Require().Equal(expected.Actions(), actual.Actions())
}
