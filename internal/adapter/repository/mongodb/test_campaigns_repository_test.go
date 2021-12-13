package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/harpyd/thestis/internal/adapter/repository/mongodb"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type TestCampaignsRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repo *mongodb.TestCampaignsRepository
}

func (s *TestCampaignsRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewTestCampaignsRepository(s.db)
}

func (s *TestCampaignsRepositoryTestSuite) TearDownTest() {
	err := s.repo.RemoveAllTestCampaigns(context.Background())
	s.Require().NoError(err)
}

func TestCampaignsRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &TestCampaignsRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *TestCampaignsRepositoryTestSuite) TestAddTestCampaign() {
	testCases := []struct {
		Name                 string
		TestCampaignsFactory func() *testcampaign.TestCampaign
		ShouldBeErr          bool
		IsErr                func(err error) bool
	}{
		{
			Name: "test_campaign",
			TestCampaignsFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New("some-id", "some campaign")
				s.Require().NoError(err)

				return tc
			},
			ShouldBeErr: false,
		},
		{
			Name: "test_campaign_without_view_name",
			TestCampaignsFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New("some-id", "")
				s.Require().NoError(err)

				return tc
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			testCampaign := c.TestCampaignsFactory()

			err := s.repo.AddTestCampaign(context.Background(), testCampaign)
			s.Require().NoError(err)

			persistedTestCampaign := s.getTestCampaign(testCampaign.ID())
			s.Require().Equal(testCampaign, persistedTestCampaign)
		})
	}
}

func (s *TestCampaignsRepositoryTestSuite) getTestCampaign(tcID string) *testcampaign.TestCampaign {
	tc, err := s.repo.GetTestCampaign(context.Background(), tcID)
	s.Require().NoError(err)

	return tc
}
