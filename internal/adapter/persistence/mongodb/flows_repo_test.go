package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type FlowsRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repo *mongodb.FlowsRepository
}

func (s *FlowsRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewFlowsRepository(s.db)
}

func (s *FlowsRepositoryTestSuite) TearDownTest() {
	err := s.repo.RemoveAllFlows(context.Background())
	s.Require().NoError(err)
}

func TestFlowsRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &FlowsRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *FlowsRepositoryTestSuite) TestUpsertFlow() {
	testCases := []struct {
		Name        string
		Before      func()
		Flow        performance.Flow
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "success_inserting_flow",
			Before: func() {
				// do not insert before
			},
			Flow: performance.UnmarshalFlow(performance.FlowParams{
				ID:            "60b7b0eb-eb62-49bc-bf76-d4fca3ad48b8",
				PerformanceID: "b1728f29-c897-4258-bad8-dd824b8f84cf",
				State:         performance.Performing,
				Transitions:   s.transitions(),
			}),
			ShouldBeErr: false,
		},
		{
			Name: "success_updating_flow",
			Before: func() {
				flow := performance.UnmarshalFlow(performance.FlowParams{
					ID:            "07e3468b-a195-4b30-81df-8e3e8d389da9",
					PerformanceID: "37a5f844-25db-4aad-a3e2-628674e7e1e5",
					State:         performance.Performing,
					Transitions:   s.transitions(),
				})

				s.addFlows(flow)
			},
			Flow: performance.UnmarshalFlow(performance.FlowParams{
				ID:            "07e3468b-a195-4b30-81df-8e3e8d389da9",
				PerformanceID: "407b3e37-a4b2-4fa1-aa47-4d75e658e455",
				State:         performance.Passed,
				Transitions:   s.transitions(),
			}),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			c.Before()

			err := s.repo.UpsertFlow(context.Background(), c.Flow)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedFlow := s.getFlow(c.Flow.ID())
			s.Require().Equal(c.Flow, persistedFlow)
		})
	}
}

func (s *FlowsRepositoryTestSuite) transitions() []performance.Transition {
	s.T().Helper()

	return []performance.Transition{
		performance.NewTransition(
			performance.Passed,
			"stage.give",
			"story.scenario.to",
			"",
		),
	}
}

func (s *FlowsRepositoryTestSuite) getFlow(flowID string) performance.Flow {
	s.T().Helper()

	flow, err := s.repo.GetFlow(context.Background(), flowID)
	s.Require().NoError(err)

	return flow
}

func (s *FlowsRepositoryTestSuite) addFlows(flows ...performance.Flow) {
	s.T().Helper()

	ctx := context.Background()

	for _, f := range flows {
		err := s.repo.UpsertFlow(ctx, f)
		s.Require().NoError(err)
	}
}
