package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type FlowRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repo *mongodb.FlowRepository
}

func (s *FlowRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewFlowRepository(s.db)
}

func (s *FlowRepositoryTestSuite) TearDownTest() {
	_, err := s.db.
		Collection("flows").
		DeleteOne(context.Background(), bson.D{})
	s.Require().NoError(err)
}

func TestFlowRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &FlowRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *FlowRepositoryTestSuite) TestUpsertFlow() {
	testCases := []struct {
		Name               string
		InsertedBeforeFlow interface{}
		GivenFlow          *flow.Flow
		ShouldBeErr        bool
		IsErr              func(err error) bool
	}{
		{
			Name: "success_inserting_flow",
			GivenFlow: flow.Unmarshal(flow.Params{
				ID:            "60b7b0eb-eb62-49bc-bf76-d4fca3ad48b8",
				PerformanceID: "b1728f29-c897-4258-bad8-dd824b8f84cf",
				Statuses: []*flow.Status{
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "bar"),
						flow.Canceled,
						flow.NewThesisStatus("baz", flow.Canceled),
					),
				},
			}),
			ShouldBeErr: false,
		},
		{
			Name: "success_updating_flow",
			InsertedBeforeFlow: bson.M{
				"_id":           "07e3468b-a195-4b30-81df-8e3e8d389da9",
				"performanceId": "37a5f844-25db-4aad-a3e2-628674e7e1e5",
			},
			GivenFlow: flow.Unmarshal(flow.Params{
				ID:            "07e3468b-a195-4b30-81df-8e3e8d389da9",
				PerformanceID: "407b3e37-a4b2-4fa1-aa47-4d75e658e455",
				Statuses: []*flow.Status{
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "bar"),
						flow.Failed,
						flow.NewThesisStatus("baz", flow.Passed),
						flow.NewThesisStatus("bad", flow.Failed),
					),
				},
			}),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			if c.InsertedBeforeFlow != nil {
				s.addFlows(c.InsertedBeforeFlow)
			}

			err := s.repo.UpsertFlow(context.Background(), c.GivenFlow)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedFlow := s.getFlow(c.GivenFlow.ID())
			s.Require().Equal(c.GivenFlow, persistedFlow)
		})
	}
}

func (s *FlowRepositoryTestSuite) getFlow(flowID string) *flow.Flow {
	s.T().Helper()

	f, err := s.repo.GetFlow(context.Background(), flowID)
	s.Require().NoError(err)

	return f
}

func (s *FlowRepositoryTestSuite) addFlows(flows ...interface{}) {
	s.T().Helper()

	ctx := context.Background()

	for _, f := range flows {
		_, err := s.db.Collection("flows").InsertOne(ctx, f)
		s.Require().NoError(err)
	}
}
