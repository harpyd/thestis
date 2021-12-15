package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/harpyd/thestis/internal/adapter/repository/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
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
				tc, err := testcampaign.New(
					"e75690c2-e659-409d-a528-ffd40d17c4bc",
					"some campaign",
					"summary",
				)
				s.Require().NoError(err)

				return tc
			},
			ShouldBeErr: false,
		},
		{
			Name: "test_campaign_without_view_name",
			TestCampaignsFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New(
					"1153796c-58d4-4b26-8c2f-f32a1a875dac",
					"",
					"summary",
				)
				s.Require().NoError(err)

				return tc
			},
			ShouldBeErr: false,
		},
		{
			Name: "test_campaign_without_summary",
			TestCampaignsFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New(
					"290360f8-8d28-437b-9fe6-9bbd23198c76",
					"view",
					"",
				)
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

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedTestCampaign := s.getTestCampaign(testCampaign.ID())
			s.Require().Equal(testCampaign, persistedTestCampaign)
		})
	}
}

func (s *TestCampaignsRepositoryTestSuite) TestUpdateTestCampaign() {
	testCampaignToUpdate, err := testcampaign.New("0b723635-4691-4eae-aca8-79b230989f9d", "some name", "summary")
	s.Require().NoError(err)

	s.addTestCampaigns(testCampaignToUpdate)

	testCases := []struct {
		Name                   string
		TestCampaignIDToUpdate string
		Update                 command.TestCampaignUpdater
		TestCampaignUpdated    func(tc *testcampaign.TestCampaign) bool
		ShouldBeErr            bool
		IsErr                  func(err error) bool
	}{
		{
			Name:                   "non_existing_test_campaign",
			TestCampaignIDToUpdate: "24206e05-630d-4903-837b-1b3615a5d802",
			Update: func(_ context.Context, tc *testcampaign.TestCampaign) (*testcampaign.TestCampaign, error) {
				return tc, nil
			},
			ShouldBeErr: true,
			IsErr:       app.IsTestCampaignNotFoundError,
		},
		{
			Name:                   "update_active_specification_id",
			TestCampaignIDToUpdate: "0b723635-4691-4eae-aca8-79b230989f9d",
			Update: func(_ context.Context, tc *testcampaign.TestCampaign) (*testcampaign.TestCampaign, error) {
				tc.SetActiveSpecificationID("4d5526c4-c588-4b77-8cf8-2a6663be9147")

				return tc, nil
			},
			TestCampaignUpdated: func(tc *testcampaign.TestCampaign) bool {
				return tc.ActiveSpecificationID() == "4d5526c4-c588-4b77-8cf8-2a6663be9147"
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			ctx := context.Background()

			err := s.repo.UpdateTestCampaign(ctx, c.TestCampaignIDToUpdate, c.Update)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			tc := s.getTestCampaign(c.TestCampaignIDToUpdate)
			s.Require().True(c.TestCampaignUpdated(tc))
		})
	}
}

func (s *TestCampaignsRepositoryTestSuite) addTestCampaigns(tcs ...*testcampaign.TestCampaign) {
	ctx := context.Background()

	for _, tc := range tcs {
		err := s.repo.AddTestCampaign(ctx, tc)
		s.Require().NoError(err)
	}
}

func (s *TestCampaignsRepositoryTestSuite) getTestCampaign(tcID string) *testcampaign.TestCampaign {
	tc, err := s.repo.GetTestCampaign(context.Background(), tcID)
	s.Require().NoError(err)

	return tc
}
