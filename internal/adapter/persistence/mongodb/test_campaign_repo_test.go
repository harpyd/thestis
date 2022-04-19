package mongodb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type TestCampaignRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repo *mongodb.TestCampaignRepository
}

func (s *TestCampaignRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewTestCampaignRepository(s.db)
}

func (s *TestCampaignRepositoryTestSuite) TearDownTest() {
	_, err := s.db.
		Collection("testCampaigns").
		DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err)
}

func TestCampaignRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &TestCampaignRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *TestCampaignRepositoryTestSuite) TestFindTestCampaign() {
	testCampaignToFind, err := testcampaign.New(testcampaign.Params{
		ID:        "c0b28d44-d603-4756-bd25-8b3034e1dc77",
		ViewName:  "some name",
		Summary:   "info",
		OwnerID:   "54112816-3a55-4a28-82df-3c8e082fa0f8",
		CreatedAt: time.Now().UTC(),
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
			Name: "by_non_existing_test_campaign_id_and_non_existing_owner_id",
			Query: app.SpecificTestCampaignQuery{
				TestCampaignID: "145337c9-ef1b-4dfd-a227-24805d25b52e",
				UserID:         "52d308aa-5a8a-4931-b558-73962e55d443",
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrTestCampaignNotFound)
			},
		},
		{
			Name: "by_non_existing_test_campaign_id_and_existing_owner_id",
			Query: app.SpecificTestCampaignQuery{
				TestCampaignID: "53a7f280-247a-410f-b6ee-c336fe9c643f",
				UserID:         "54112816-3a55-4a28-82df-3c8e082fa0f8",
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrTestCampaignNotFound)
			},
		},
		{
			Name: "by_existing_test_campaign_id_and_non_existing_owner_id",
			Query: app.SpecificTestCampaignQuery{
				TestCampaignID: "c0b28d44-d603-4756-bd25-8b3034e1dc77",
				UserID:         "0c972872-b033-4fb1-a29b-2e14a8eb56f4",
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrTestCampaignNotFound)
			},
		},
		{
			Name: "by_existing_test_campaign_id_and_existing_owner_id",
			Query: app.SpecificTestCampaignQuery{
				TestCampaignID: "c0b28d44-d603-4756-bd25-8b3034e1dc77",
				UserID:         "54112816-3a55-4a28-82df-3c8e082fa0f8",
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
			s.Require().Equal(testCampaignToFind.Summary(), tc.Summary)
			s.Require().Equal(testCampaignToFind.ViewName(), tc.ViewName)
			s.requireTimesEqualWithRounding(testCampaignToFind.CreatedAt(), tc.CreatedAt)
		})
	}
}

func (s *TestCampaignRepositoryTestSuite) TestAddTestCampaign() {
	testCases := []struct {
		Name         string
		Before       func()
		TestCampaign *testcampaign.TestCampaign
		ShouldBeErr  bool
		IsErr        func(err error) bool
	}{
		{
			Name: "failed_adding_one_test_campaign_twice",
			Before: func() {
				tc, err := testcampaign.New(testcampaign.Params{
					ID:      "d7822ba3-7bec-48c8-8b32-491108a75390",
					OwnerID: "a40b27f2-1206-4309-be31-4100a9d7c0c8",
				})
				s.Require().NoError(err)

				s.addTestCampaigns(tc)
			},
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:        "d7822ba3-7bec-48c8-8b32-491108a75390",
				OwnerID:   "42154bef-3d63-4341-a989-71932ecb4220",
				ViewName:  "bbbbbo",
				Summary:   "ssssso",
				CreatedAt: time.Now().UTC(),
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *app.DatabaseError

				return errors.As(err, &target) &&
					mongo.IsDuplicateKeyError(err)
			},
		},
		{
			Name: "successful_adding_with_all_fields",
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:        "e75690c2-e659-409d-a528-ffd40d17c4bc",
				ViewName:  "some campaign",
				Summary:   "summary",
				OwnerID:   "6c11693f-3376-4873-a8ef-a77a327ccb46",
				CreatedAt: time.Now().UTC(),
			}),
			ShouldBeErr: false,
		},
		{
			Name: "successful_adding_without_view_name",
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:        "1153796c-58d4-4b26-8c2f-f32a1a875dac",
				ViewName:  "",
				Summary:   "summary",
				OwnerID:   "9c845592-5e9e-4160-8e2f-0309a6949f04",
				CreatedAt: time.Now().UTC(),
			}),
			ShouldBeErr: false,
		},
		{
			Name: "successful_adding_without_summary",
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:        "9ed07209-7a4f-4dd9-bf0d-6f8b70280f85",
				ViewName:  "view name name name",
				Summary:   "",
				OwnerID:   "54112816-3a55-4a28-82df-3c8e082fa0f8",
				CreatedAt: time.Now().UTC(),
			}),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			if c.Before != nil {
				c.Before()
			}

			err := s.repo.AddTestCampaign(context.Background(), c.TestCampaign)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedTestCampaign := s.getTestCampaign(c.TestCampaign.ID())
			s.requireTestCampaignsEqual(c.TestCampaign, persistedTestCampaign)
		})
	}
}

func (s *TestCampaignRepositoryTestSuite) TestUpdateTestCampaign() {
	testCampaignToUpdate, err := testcampaign.New(testcampaign.Params{
		ID:        "0b723635-4691-4eae-aca8-79b230989f9d",
		OwnerID:   "3dd1ee11-2520-4de1-859a-b8d6fbb003e9",
		ViewName:  "some name",
		Summary:   "summary",
		CreatedAt: time.Now().UTC(),
	})
	s.Require().NoError(err)

	s.addTestCampaigns(testCampaignToUpdate)

	testCases := []struct {
		Name                   string
		TestCampaignIDToUpdate string
		Update                 app.TestCampaignUpdater
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
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrTestCampaignNotFound)
			},
		},
		{
			Name:                   "update_summary",
			TestCampaignIDToUpdate: "0b723635-4691-4eae-aca8-79b230989f9d",
			Update: func(_ context.Context, tc *testcampaign.TestCampaign) (*testcampaign.TestCampaign, error) {
				tc.SetSummary("new summary")

				return tc, nil
			},
			TestCampaignUpdated: func(tc *testcampaign.TestCampaign) bool {
				return tc.Summary() == "new summary"
			},
			ShouldBeErr: false,
		},
		{
			Name:                   "update_view_name",
			TestCampaignIDToUpdate: "0b723635-4691-4eae-aca8-79b230989f9d",
			Update: func(_ context.Context, tc *testcampaign.TestCampaign) (*testcampaign.TestCampaign, error) {
				tc.SetViewName("new view name")

				return tc, nil
			},
			TestCampaignUpdated: func(tc *testcampaign.TestCampaign) bool {
				return tc.ViewName() == "new view name"
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

func (s *TestCampaignRepositoryTestSuite) addTestCampaigns(tcs ...*testcampaign.TestCampaign) {
	s.T().Helper()

	ctx := context.Background()

	for _, tc := range tcs {
		err := s.repo.AddTestCampaign(ctx, tc)
		s.Require().NoError(err)
	}
}

func (s *TestCampaignRepositoryTestSuite) getTestCampaign(tcID string) *testcampaign.TestCampaign {
	s.T().Helper()

	tc, err := s.repo.GetTestCampaign(context.Background(), tcID)
	s.Require().NoError(err)

	return tc
}

func (s *TestCampaignRepositoryTestSuite) requireTestCampaignsEqual(expected, actual *testcampaign.TestCampaign) {
	s.Require().Equal(expected.ID(), actual.ID())
	s.Require().Equal(expected.OwnerID(), actual.OwnerID())
	s.Require().Equal(expected.Summary(), actual.Summary())
	s.Require().Equal(expected.ViewName(), actual.ViewName())
	s.requireTimesEqualWithRounding(expected.CreatedAt(), actual.CreatedAt())
}

func (s *TestCampaignRepositoryTestSuite) requireTimesEqualWithRounding(expected, actual time.Time) {
	s.T().Helper()

	s.Require().True(round(expected).Equal(round(actual)))
}

func round(t time.Time) time.Time {
	return t.Round(1 * time.Second)
}
