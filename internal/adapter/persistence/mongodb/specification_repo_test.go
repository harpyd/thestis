package mongodb_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type SpecificationRepositoryTestSuite struct {
	MongoSuite

	repo *mongodb.SpecificationRepository
}

func (s *SpecificationRepositoryTestSuite) SetupTest() {
	s.repo = mongodb.NewSpecificationRepository(s.db)
}

func (s *SpecificationRepositoryTestSuite) TearDownTest() {
	_, err := s.db.
		Collection("specifications").
		DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err)
}

func TestSpecificationRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &SpecificationRepositoryTestSuite{})
}

func (s *SpecificationRepositoryTestSuite) TestFindSpecification() {
	insertedSpec := s.rawSpecification()

	s.insertSpecifications(insertedSpec)

	testCases := []struct {
		Name        string
		Query       app.SpecificSpecificationQuery
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "by_non_existing_specification_id_and_non_existing_owner_id",
			Query: app.SpecificSpecificationQuery{
				SpecificationID: "34eff819-c14b-4d89-98a9-8e21d9f3cf21",
				UserID:          "b5f0e13d-ca3a-41f9-b297-53a2440c6080",
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrSpecificationNotFound)
			},
		},
		{
			Name: "by_existing_specification_id_and_non_existing_owner_id",
			Query: app.SpecificSpecificationQuery{
				SpecificationID: "64825e35-7fa7-44a4-9ca2-81cfc7b0f0d8",
				UserID:          "4699c306-ba54-4a5d-916e-92c40646faca",
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrSpecificationNotFound)
			},
		},
		{
			Name: "by_non_existing_specification_id_and_existing_owner_id",
			Query: app.SpecificSpecificationQuery{
				SpecificationID: "75d48408-e86a-427b-9f22-a77a22f13348",
				UserID:          "52d9af60-26be-46ea-90a6-efec5fbb4ccd",
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrSpecificationNotFound)
			},
		},
		{
			Name: "by_existing_specification_id_and_existing_owner_id",
			Query: app.SpecificSpecificationQuery{
				SpecificationID: "64825e35-7fa7-44a4-9ca2-81cfc7b0f0d8",
				UserID:          "52d9af60-26be-46ea-90a6-efec5fbb4ccd",
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

			s.requireAppSpecificationEqualRaw(insertedSpec, spec)
		})
	}
}

func (s *SpecificationRepositoryTestSuite) rawSpecification() bson.M {
	s.T().Helper()

	return bson.M{
		"id":             "64825e35-7fa7-44a4-9ca2-81cfc7b0f0d8",
		"ownerId":        "52d9af60-26be-46ea-90a6-efec5fbb4ccd",
		"testCampaignId": "1896c290-1dde-42f5-8449-e197f49daad2",
		"loadedAt":       time.Now().UTC(),
		"author":         "Djerys",
		"title":          "test",
		"description":    "desc",
		"stories": []bson.M{
			{
				"slug":        "story",
				"description": "desc",
				"asA":         "some",
				"inOrderTo":   "some",
				"wantTo":      "want some",
				"scenarios": []bson.M{
					{
						"slug":        "scenario",
						"description": "desc",
						"theses": []bson.M{
							{
								"slug": "a",
								"statement": bson.M{
									"stage":    specification.Given,
									"behavior": "a",
								},
								"http": bson.M{
									"request": bson.M{
										"method":      specification.POST,
										"url":         "https://some-domain.com",
										"contentType": specification.ApplicationJSON,
										"body": map[string]interface{}{
											"bar": "baz",
											"bad": "foo",
										},
									},
									"response": bson.M{
										"allowedCodes":       []int{201},
										"allowedContentType": specification.ApplicationJSON,
									},
								},
							},
							{
								"slug":  "b",
								"after": []string{"a"},
								"statement": bson.M{
									"stage":    specification.Given,
									"behavior": "b",
								},
								"assertion": bson.M{
									"method": specification.JSONPath,
									"asserts": []bson.M{
										{
											"actual":   "some.field",
											"expected": "foo",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *SpecificationRepositoryTestSuite) TestGetActiveSpecificationByTestCampaignID() {
	testCampaignID := "d0832b59-6e8a-46f6-9b57-92e8bf656e93"

	var (
		firstSpec = bson.M{
			"id":             "358a938f-8191-4264-8070-4ac5914bc130",
			"testCampaignId": testCampaignID,
		}
		secondSpec = bson.M{
			"id":             "8c1058aa-295a-47a0-83e9-a128c2bd22af",
			"testCampaignId": testCampaignID,
		}
		lastSpec = bson.M{
			"id":             "aa056dc5-b0e7-4695-a209-1d46805373c6",
			"testCampaignId": testCampaignID,
		}
	)

	s.insertSpecifications(firstSpec, secondSpec, lastSpec)

	testCases := []struct {
		Name           string
		TestCampaignID string
		ShouldBeErr    bool
		IsErr          func(err error) bool
	}{
		{
			Name:           "non_existing_specification_with_test_campaign_id",
			TestCampaignID: "1ba42415-588b-4a11-ab06-76b4a298658e",
			ShouldBeErr:    true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrSpecificationNotFound)
			},
		},
		{
			Name:           "last_added_specification",
			TestCampaignID: testCampaignID,
			ShouldBeErr:    false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			spec, err := s.repo.GetActiveSpecificationByTestCampaignID(
				context.Background(),
				c.TestCampaignID,
			)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			var b specification.Builder

			lastBuiltSpec := b.
				WithID("aa056dc5-b0e7-4695-a209-1d46805373c6").
				WithTestCampaignID(testCampaignID).
				ErrlessBuild()

			s.Require().Equal(lastBuiltSpec, spec)
		})
	}
}

func (s *SpecificationRepositoryTestSuite) TestAddSpecification() {
	testCases := []struct {
		Name                        string
		InsertedBeforeSpecification interface{}
		Specification               *specification.Specification
		ShouldBeErr                 bool
		IsErr                       func(err error) bool
	}{
		{
			Name: "failed_adding_one_specification_twice",
			InsertedBeforeSpecification: bson.M{
				"id": "62a4e06b-c00f-49a5-a1c1-5906e5e2e1d5",
			},
			Specification: (&specification.Builder{}).
				WithID("62a4e06b-c00f-49a5-a1c1-5906e5e2e1d5").
				WithAuthor("Djerys").
				WithOwnerID("6c204a55-023f-49bf-8c3d-1e7915b64f3a").
				ErrlessBuild(),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *app.DatabaseError

				return errors.As(err, &target) &&
					mongo.IsDuplicateKeyError(err)
			},
		},
		{
			Name: "successful_adding",
			Specification: (&specification.Builder{}).
				WithID("f517f320-7d07-44a5-9fbf-7e1eb6889e87").
				WithAuthor("Djerys").
				WithTitle("Test title").
				WithDescription("Test description").
				WithOwnerID("393a989b-31a2-4c52-a6bd-abd83f5b2392").
				WithTestCampaignID("e0a9361a-3605-4116-bb9b-957d9e0460f8").
				ErrlessBuild(),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		s.Run(c.Name, func() {
			if c.InsertedBeforeSpecification != nil {
				s.insertSpecifications(c.InsertedBeforeSpecification)
			}

			err := s.repo.AddSpecification(context.Background(), c.Specification)

			if c.ShouldBeErr {
				s.Require().True(c.IsErr(err))

				return
			}

			s.Require().NoError(err)

			persistedSpec := s.getSpecification(c.Specification.ID())
			s.Require().Equal(c.Specification, persistedSpec)
		})
	}
}

func (s *SpecificationRepositoryTestSuite) getSpecification(specID string) *specification.Specification {
	s.T().Helper()

	spec, err := s.repo.GetSpecification(context.Background(), specID)
	s.Require().NoError(err)

	return spec
}

func (s *SpecificationRepositoryTestSuite) requireAppSpecificationEqualRaw(
	expected bson.M,
	actual app.SpecificSpecification,
) {
	s.T().Helper()

	s.Require().Equal(expected["id"], actual.ID)
	s.Require().Equal(expected["testCampaignId"], actual.TestCampaignID)
	expectedLoadedAt, ok := expected["loadedAt"].(time.Time)
	s.Require().True(ok)
	s.Require().WithinDuration(expectedLoadedAt, actual.LoadedAt, 1*time.Second)
	s.Require().Equal(expected["author"], actual.Author)
	s.Require().Equal(expected["title"], actual.Title)
	s.Require().Equal(expected["description"], actual.Description)
	expectedStories, ok := expected["stories"].([]bson.M)
	s.Require().True(ok)
	s.Require().Equal(len(expectedStories), len(actual.Stories))

	for i := range expectedStories {
		s.requireAppStoryEqualRaw(expectedStories[i], actual.Stories[i])
	}
}

func (s *SpecificationRepositoryTestSuite) requireAppStoryEqualRaw(expected bson.M, actual app.Story) {
	s.T().Helper()

	s.Require().Equal(expected["slug"], actual.Slug)
	s.Require().Equal(expected["description"], actual.Description)
	s.Require().Equal(expected["asA"], actual.AsA)
	s.Require().Equal(expected["inOrderTo"], actual.InOrderTo)
	s.Require().Equal(expected["wantTo"], actual.WantTo)
	expectedScenarios, ok := expected["scenarios"].([]bson.M)
	s.Require().True(ok)
	s.Require().Equal(len(expectedScenarios), len(actual.Scenarios))

	for i := range expectedScenarios {
		s.requireAppScenarioEqualRaw(expectedScenarios[i], actual.Scenarios[i])
	}
}

func (s *SpecificationRepositoryTestSuite) requireAppScenarioEqualRaw(expected bson.M, actual app.Scenario) {
	s.T().Helper()

	s.Require().Equal(expected["slug"], actual.Slug)
	s.Require().Equal(expected["description"], actual.Description)
	expectedTheses, ok := expected["theses"].([]bson.M)
	s.Require().True(ok)
	s.Require().Equal(len(expectedTheses), len(actual.Theses))

	for i := range expectedTheses {
		s.requireAppThesisEqualRaw(expectedTheses[i], actual.Theses[i])
	}
}

func (s *SpecificationRepositoryTestSuite) requireAppThesisEqualRaw(expected bson.M, actual app.Thesis) {
	s.T().Helper()

	s.Require().Equal(expected["slug"], actual.Slug)

	expectedAfter, ok := expected["after"].([]string)
	if ok {
		s.Require().ElementsMatch(expectedAfter, actual.After)
	} else {
		s.Require().Empty(actual.After)
	}

	expectedStatement, ok := expected["statement"].(bson.M)
	s.Require().True(ok)
	s.requireAppStatementEqualRaw(expectedStatement, actual.Statement)

	expectedHTTP, ok := expected["http"].(bson.M)
	actualHTTP := actual.HTTP

	s.Require().Equal(ok, !actualHTTP.IsZero())

	if ok {
		s.requireAppHTTPEqualRaw(expectedHTTP, actualHTTP)
	}

	expectedAssertion, ok := expected["assertion"].(bson.M)
	actualAssertion := actual.Assertion

	s.Require().Equal(ok, !actualAssertion.IsZero())

	if ok {
		s.requireAppAssertionEqualRaw(expectedAssertion, actualAssertion)
	}
}

func (s *SpecificationRepositoryTestSuite) requireAppStatementEqualRaw(expected bson.M, actual app.Statement) {
	s.T().Helper()

	s.Require().Equal(fmt.Sprintf("%s", expected["stage"]), actual.Stage)
	s.Require().Equal(expected["behavior"], actual.Behavior)
}

func (s *SpecificationRepositoryTestSuite) requireAppHTTPEqualRaw(expected bson.M, actual app.HTTP) {
	s.T().Helper()

	expectedHTTPReq, ok := expected["request"].(bson.M)
	actualHTTPReq := actual.Request

	s.Require().Equal(ok, !actualHTTPReq.IsZero())

	if ok {
		s.requireAppHTTPRequestEqualRaw(expectedHTTPReq, actualHTTPReq)
	}

	expectedHTTPRes, ok := expected["response"].(bson.M)
	actualHTTPRes := actual.Response

	s.Require().Equal(ok, !actualHTTPRes.IsZero())

	if ok {
		s.requireAppHTTPResponseEqualRaw(expectedHTTPRes, actualHTTPRes)
	}
}

func (s *SpecificationRepositoryTestSuite) requireAppHTTPRequestEqualRaw(
	expected bson.M,
	actual app.HTTPRequest,
) {
	s.T().Helper()

	s.Require().Equal(fmt.Sprintf("%s", expected["method"]), actual.Method)
	s.Require().Equal(expected["url"], actual.URL)
	s.Require().Equal(fmt.Sprintf("%s", expected["contentType"]), actual.ContentType)
	s.Require().Equal(expected["body"], actual.Body)
}

func (s *SpecificationRepositoryTestSuite) requireAppHTTPResponseEqualRaw(
	expected bson.M,
	actual app.HTTPResponse,
) {
	s.T().Helper()

	s.Require().Equal(expected["allowedCodes"], actual.AllowedCodes)
	s.Require().Equal(
		fmt.Sprintf("%s", expected["allowedContentType"]),
		actual.AllowedContentType,
	)
}

func (s *SpecificationRepositoryTestSuite) requireAppAssertionEqualRaw(expected bson.M, actual app.Assertion) {
	s.T().Helper()

	s.Require().Equal(
		fmt.Sprintf("%s", expected["method"]),
		actual.Method,
	)

	expectedAsserts, ok := expected["asserts"].([]bson.M)
	actualAsserts := actual.Asserts

	s.Require().True(ok)
	s.Require().Equal(len(expectedAsserts), len(actualAsserts))

	for i := range expectedAsserts {
		s.requireAppAssertEqualRaw(expectedAsserts[i], actualAsserts[i])
	}
}

func (s *SpecificationRepositoryTestSuite) requireAppAssertEqualRaw(expected bson.M, actual app.Assert) {
	s.T().Helper()

	s.Require().Equal(expected["actual"], actual.Actual)
	s.Require().Equal(expected["expected"], actual.Expected)
}
