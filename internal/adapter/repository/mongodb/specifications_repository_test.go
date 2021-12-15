package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/harpyd/thestis/internal/adapter/repository/mongodb"
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

func (s *SpecificationsRepositoryTestSuite) TestAddSpecification() {
	testCases := []struct {
		Name                 string
		SpecificationFactory func() *specification.Specification
		ShouldBeErr          bool
		IsErr                func(err error) bool
	}{
		{
			Name: "with_id_author_and_title",
			SpecificationFactory: func() *specification.Specification {
				return specification.NewBuilder().
					WithID("f517f320-7d07-44a5-9fbf-7e1eb6889e87").
					WithAuthor("Djerys").
					WithTitle("Test title").
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
	spec, err := s.repo.GetSpecification(context.Background(), specID)
	s.Require().NoError(err)

	return spec
}
