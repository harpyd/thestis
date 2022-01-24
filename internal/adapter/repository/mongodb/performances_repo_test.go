package mongodb_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/harpyd/thestis/internal/adapter/repository/mongodb"
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
		Name               string
		PerformanceFactory func() *performance.Performance
		ShouldBeErr        bool
		IsErr              func(err error) bool
	}{
		{
			Name: "successful_adding",
			PerformanceFactory: func() *performance.Performance {
				return performance.UnmarshalFromDatabase(performance.Params{
					OwnerID:         "3614a95c-c278-4687-84e2-97b95b11d399",
					SpecificationID: "4e4465b0-a312-4f86-9051-a3ae72965215",
					Actions:         []performance.Action{s.action()},
				}, performance.WithID("3ce098e1-81ae-4610-8372-2f635b1b6a0c"))
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			perf := c.PerformanceFactory()

			err := s.repo.AddPerformance(context.Background(), perf)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedPerf := s.getPerformance(perf.ID())
			s.requirePerformancesEqual(perf, persistedPerf)
		})
	}
}

func (s *PerformancesRepositoryTestSuite) TestExclusivelyDoWithPerformance_concurrent_actions() {
	perf := performance.UnmarshalFromDatabase(performance.Params{
		OwnerID:         "e6cd6e6d-f58f-4a3e-a4d3-6b23dce29750",
		SpecificationID: "d91da0ce-1caa-43d6-95c0-1a03a9d3cd52",
		Actions:         []performance.Action{s.action()},
	}, performance.WithID("9a07bd86-3b6a-4202-88ec-633c1b5a1e91"))

	s.addPerformances(perf)

	const actionsNumber = 100

	var (
		finish = make(chan bool)
		errors = make(chan error)
	)

	go func() {
		defer close(errors)

		ctx := context.Background()

		var wg sync.WaitGroup

		wg.Add(actionsNumber)

		for i := 1; i <= actionsNumber; i++ {
			go func() {
				defer wg.Done()

				perfCopy := s.getPerformance(perf.ID())

				if err := s.repo.ExclusivelyDoWithPerformance(ctx, perfCopy, func(perf *performance.Performance) {
					finish <- true
				}); err != nil {
					errors <- err
				}
			}()
		}

		wg.Wait()
	}()

	alreadyStartedErrsCount := 0

	for err := range errors {
		if performance.IsAlreadyStartedError(err) {
			alreadyStartedErrsCount++
		}
	}

	s.Require().Equal(actionsNumber-1, alreadyStartedErrsCount)

	<-finish
}

func (s *PerformancesRepositoryTestSuite) action() performance.Action {
	s.T().Helper()

	thesis, err := specification.NewThesisBuilder().
		WithStatement("given", "test").
		WithAssertion(func(b *specification.AssertionBuilder) {
			b.WithMethod("jsonpath")
			b.WithAssert("some", "value")
		}).
		Build("to")
	s.Require().NoError(err)

	return performance.NewAction("stage.given", "story.scenario.to", thesis, performance.AssertionPerformer)
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
