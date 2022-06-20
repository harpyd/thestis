package mongodb_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/harpyd/thestis/internal/core/adapter/driven/persistence/mongodb"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

type PipelineGuardTestSuite struct {
	MongoSuite

	guard  *mongodb.PipelineGuard
	pipeID string
}

func (s *PipelineGuardTestSuite) SetupTest() {
	s.guard = mongodb.NewPipelineGuard(s.db)

	s.pipeID = "2db44433-7142-4080-bada-844afccfedbf"

	s.insertPipelines(bson.M{
		"_id":     s.pipeID,
		"started": false,
	})
}

func (s *PipelineGuardTestSuite) TearDownTest() {
	_, err := s.db.
		Collection("pipelines").
		DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err)
}

func TestPipelineGuard(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &PipelineGuardTestSuite{})
}

func (s *PipelineGuardTestSuite) TestAcquirePipeline() {
	ctx := context.Background()

	err := s.guard.AcquirePipeline(ctx, s.pipeID)
	s.Require().NoError(err)

	err = s.guard.AcquirePipeline(ctx, s.pipeID)
	s.Require().ErrorIs(err, pipeline.ErrAlreadyStarted)

	s.Require().True(s.getPipelineStarted())
}

func (s *PipelineGuardTestSuite) TestAcquirePipelineConcurrently() {
	const triesCount = 100

	var wg sync.WaitGroup

	var acquiredCount int32

	wg.Add(triesCount)

	for i := 1; i <= triesCount; i++ {
		go func() {
			defer wg.Done()

			if err := s.guard.AcquirePipeline(
				context.Background(),
				s.pipeID,
			); err == nil {
				atomic.AddInt32(&acquiredCount, 1)
			}
		}()
	}

	wg.Wait()

	s.Require().EqualValues(1, acquiredCount)
}

func (s *PipelineGuardTestSuite) TestReleasePipeline() {
	err := s.guard.AcquirePipeline(context.Background(), s.pipeID)
	s.Require().NoError(err)

	s.Require().True(s.getPipelineStarted())

	err = s.guard.ReleasePipeline(context.Background(), s.pipeID)
	s.Require().NoError(err)

	s.Require().False(s.getPipelineStarted())
}

func (s *PipelineGuardTestSuite) getPipelineStarted() bool {
	var document struct {
		Started bool `bson:"started"`
	}

	err := s.db.Collection("pipelines").
		FindOne(context.Background(), bson.M{"_id": s.pipeID}).
		Decode(&document)
	s.Require().NoError(err)

	return document.Started
}
