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

func (s *TestCampaignsRepositoryTestSuite) TestFindTestCampaign() {
	testCampaignToFind, err := testcampaign.New(testcampaign.Params{
		ID:       "c0b28d44-d603-4756-bd25-8b3034e1dc77",
		ViewName: "some name",
		Summary:  "info",
		UserID:   "54112816-3a55-4a28-82df-3c8e082fa0f8",
	})
	s.Require().NoError(err)

	s.addTestCampaigns(testCampaignToFind)

	testCases := []struct {
		Name        string
		Query       app.SpecificTestCampaignQuery
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "non_existing_test_campaign",
			Query: app.SpecificTestCampaignQuery{
				TestCampaignID: "53a7f280-247a-410f-b6ee-c336fe9c643f",
			},
			ShouldBeErr: true,
			IsErr:       app.IsTestCampaignNotFoundError,
		},
		{
			Name: "by_existing_test_campaign_id",
			Query: app.SpecificTestCampaignQuery{
				TestCampaignID: "c0b28d44-d603-4756-bd25-8b3034e1dc77",
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			tc, err := s.repo.FindTestCampaign(context.Background(), c.Query)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			s.Require().Equal(testCampaignToFind.ID(), tc.ID)
		})
	}
}

func (s *TestCampaignsRepositoryTestSuite) TestAddTestCampaign() {
	testCases := []struct {
		Name                 string
		TestCampaignsFactory func() *testcampaign.TestCampaign
		ShouldBeErr          bool
		IsErr                func(err error) bool
	}{
		{
			Name: "with_all_fields",
			TestCampaignsFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New(testcampaign.Params{
					ID:       "e75690c2-e659-409d-a528-ffd40d17c4bc",
					ViewName: "some campaign",
					Summary:  "summary",
					UserID:   "6c11693f-3376-4873-a8ef-a77a327ccb46",
				})
				s.Require().NoError(err)

				return tc
			},
			ShouldBeErr: false,
		},
		{
			Name: "without_view_name",
			TestCampaignsFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New(testcampaign.Params{
					ID:       "1153796c-58d4-4b26-8c2f-f32a1a875dac",
					ViewName: "",
					Summary:  "summary",
					UserID:   "9c845592-5e9e-4160-8e2f-0309a6949f04",
				})
				s.Require().NoError(err)

				return tc
			},
			ShouldBeErr: false,
		},
		{
			Name: "without_summary",
			TestCampaignsFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New(testcampaign.Params{
					ID:       "1153796c-58d4-4b26-8c2f-f32a1a875dac",
					ViewName: "view name name name",
					Summary:  "",
					UserID:   "54112816-3a55-4a28-82df-3c8e082fa0f8",
				})
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
	testCampaignToUpdate, err := testcampaign.New(testcampaign.Params{
		ID:       "0b723635-4691-4eae-aca8-79b230989f9d",
		ViewName: "some name",
		Summary:  "summary",
	})
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
