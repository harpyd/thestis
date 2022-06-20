package mongodb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/core/adapter/driven/persistence/mongodb"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type PipelineRepositoryTestSuite struct {
	MongoSuite

	repo *mongodb.PipelineRepository
}

func (s *PipelineRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewPipelineRepository(s.db)
}

func (s *PipelineRepositoryTestSuite) TearDownTest() {
	_, err := s.db.
		Collection("pipelines").
		DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err)
}

func TestPipelineRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &PipelineRepositoryTestSuite{})
}

func (s *PipelineRepositoryTestSuite) TestAddPipeline() {
	var b specification.Builder

	availableSpec := b.
		WithID("1a63b6ea-df5f-4a68-bf04-c2e30044f2ef").
		ErrlessBuild()

	testCases := []struct {
		Name                   string
		InsertedBeforePipeline interface{}
		GivenPipeline          *pipeline.Pipeline
		ShouldBeErr            bool
		IsErr                  func(err error) bool
	}{
		{
			Name: "failed_adding_one_pipeline_twice",
			InsertedBeforePipeline: bson.M{
				"_id": "a4a2906d-4df5-42f1-8832-77a33cba4d7f",
			},
			GivenPipeline: pipeline.Unmarshal(pipeline.Params{
				ID:            "a4a2906d-4df5-42f1-8832-77a33cba4d7f",
				OwnerID:       "05bf69e9-d7b5-4e7a-8fab-24227dca033a",
				Specification: availableSpec,
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *service.DatabaseError

				return errors.As(err, &target) &&
					mongo.IsDuplicateKeyError(err)
			},
		},
		{
			Name: "successful_adding",
			GivenPipeline: pipeline.Unmarshal(pipeline.Params{
				ID:            "3ce098e1-81ae-4610-8372-2f635b1b6a0c",
				OwnerID:       "3614a95c-c278-4687-84e2-97b95b11d399",
				Specification: availableSpec,
				Started:       true,
			}),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			if c.InsertedBeforePipeline != nil {
				s.insertPipelines(c.InsertedBeforePipeline)
			}

			err := s.repo.AddPipeline(context.Background(), c.GivenPipeline)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedPipe, err := s.repo.GetPipeline(
				context.Background(),
				c.GivenPipeline.ID(),
				service.AvailableSpecification(availableSpec),
			)
			s.Require().NoError(err)

			s.Require().Equal(c.GivenPipeline, persistedPipe)
		})
	}
}
