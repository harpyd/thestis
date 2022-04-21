package mongodb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/infrastructure/persistence/mongodb"
)

type PerformanceRepositoryTestSuite struct {
	MongoSuite

	repo *mongodb.PerformanceRepository
}

func (s *PerformanceRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewPerformanceRepository(s.db)
}

func (s *PerformanceRepositoryTestSuite) TearDownTest() {
	_, err := s.db.
		Collection("performances").
		DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err)
}

func TestPerformanceRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &PerformanceRepositoryTestSuite{})
}

func (s *PerformanceRepositoryTestSuite) TestAddPerformance() {
	var b specification.Builder

	availableSpec := b.
		WithID("1a63b6ea-df5f-4a68-bf04-c2e30044f2ef").
		ErrlessBuild()

	testCases := []struct {
		Name                      string
		InsertedBeforePerformance interface{}
		GivenPerformance          *performance.Performance
		ShouldBeErr               bool
		IsErr                     func(err error) bool
	}{
		{
			Name: "failed_adding_one_performance_twice",
			InsertedBeforePerformance: bson.M{
				"_id": "a4a2906d-4df5-42f1-8832-77a33cba4d7f",
			},
			GivenPerformance: performance.Unmarshal(performance.Params{
				ID:            "a4a2906d-4df5-42f1-8832-77a33cba4d7f",
				OwnerID:       "05bf69e9-d7b5-4e7a-8fab-24227dca033a",
				Specification: availableSpec,
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *app.DatabaseError

				return errors.As(err, &target) &&
					mongo.IsDuplicateKeyError(err)
			},
		},
		{
			Name: "successful_adding",
			GivenPerformance: performance.Unmarshal(performance.Params{
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
			if c.InsertedBeforePerformance != nil {
				s.insertPerformances(c.InsertedBeforePerformance)
			}

			err := s.repo.AddPerformance(context.Background(), c.GivenPerformance)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedPerf, err := s.repo.GetPerformance(
				context.Background(),
				c.GivenPerformance.ID(),
				app.AvailableSpecification(availableSpec),
			)
			s.Require().NoError(err)

			s.Require().Equal(c.GivenPerformance, persistedPerf)
		})
	}
}
