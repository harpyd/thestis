package mongodb_test

import (
	"context"
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
				thesis, err := specification.NewThesisBuilder().
					WithStatement("given", "test").
					WithAssertion(func(b *specification.AssertionBuilder) {
						b.WithMethod("jsonpath")
						b.WithAssert("some", "value")
					}).
					Build("to")
				s.Require().NoError(err)

				return performance.UnmarshalFromDatabase(performance.Params{
					OwnerID:         "3614a95c-c278-4687-84e2-97b95b11d399",
					SpecificationID: "4e4465b0-a312-4f86-9051-a3ae72965215",
					Actions: []performance.ActionParam{
						{
							From:   "stage.given",
							To:     "story.scenario.to",
							Thesis: thesis,
						},
					},
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

func (s *PerformancesRepositoryTestSuite) getPerformance(perfID string) *performance.Performance {
	s.T().Helper()

	perf, err := s.repo.GetPerformance(context.Background(), perfID)
	s.Require().NoError(err)

	return perf
}

func (s *PerformancesRepositoryTestSuite) requirePerformancesEqual(expected, actual *performance.Performance) {
	s.T().Helper()

	s.Require().Equal(expected.ID(), actual.ID())
	s.Require().Equal(expected.OwnerID(), actual.OwnerID())
	s.Require().Equal(expected.SpecificationID(), actual.SpecificationID())
	s.Require().Equal(expected.Actions(), actual.Actions())
}
