package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationsRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repo *mongodb.SpecificationsRepository
}

func (s *SpecificationsRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewSpecificationsRepository(s.db)
}

func (s *SpecificationsRepositoryTestSuite) TearDownTest() {
	err := s.repo.RemoveAllSpecifications(context.Background())
	s.Require().NoError(err)
}

func TestSpecificationsRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &SpecificationsRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *SpecificationsRepositoryTestSuite) TestFindSpecification() {
	specificationToFind := specification.NewBuilder().
		WithID("64825e35-7fa7-44a4-9ca2-81cfc7b0f0d8").
		WithAuthor("Djerys").
		WithTitle("test").
		ErrlessBuild()

	s.addSpecifications(specificationToFind)

	testCases := []struct {
		Name        string
		Query       app.SpecificSpecificationQuery
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "non_existing_specification",
			Query: app.SpecificSpecificationQuery{
				SpecificationID: "34eff819-c14b-4d89-98a9-8e21d9f3cf21",
			},
			ShouldBeErr: true,
			IsErr:       app.IsSpecificationNotFoundError,
		},
		{
			Name: "by_existing_specification_id",
			Query: app.SpecificSpecificationQuery{
				SpecificationID: "64825e35-7fa7-44a4-9ca2-81cfc7b0f0d8",
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			spec, err := s.repo.FindSpecification(context.Background(), c.Query)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			s.Require().Equal(specificationToFind.ID(), spec.ID)
		})
	}
}

func (s *SpecificationsRepositoryTestSuite) TestAddSpecification() {
	testCases := []struct {
		Name                 string
		SpecificationFactory func() *specification.Specification
		ShouldBeErr          bool
		IsErr                func(err error) bool
	}{
		{
			Name: "successful_adding",
			SpecificationFactory: func() *specification.Specification {
				return specification.NewBuilder().
					WithID("f517f320-7d07-44a5-9fbf-7e1eb6889e87").
					WithAuthor("Djerys").
					WithTitle("Test title").
					WithDescription("Test description").
					WithOwnerID("393a989b-31a2-4c52-a6bd-abd83f5b2392").
					WithTestCampaignID("e0a9361a-3605-4116-bb9b-957d9e0460f8").
					ErrlessBuild()
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			spec := c.SpecificationFactory()

			err := s.repo.AddSpecification(context.Background(), spec)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedSpec := s.getSpecification(spec.ID())
			s.Require().Equal(spec, persistedSpec)
		})
	}
}

func (s *SpecificationsRepositoryTestSuite) getSpecification(specID string) *specification.Specification {
	s.T().Helper()

	spec, err := s.repo.GetSpecification(context.Background(), specID)
	s.Require().NoError(err)

	return spec
}

func (s *SpecificationsRepositoryTestSuite) addSpecifications(specs ...*specification.Specification) {
	s.T().Helper()

	ctx := context.Background()

	for _, spec := range specs {
		err := s.repo.AddSpecification(ctx, spec)
		s.Require().NoError(err)
	}
}
