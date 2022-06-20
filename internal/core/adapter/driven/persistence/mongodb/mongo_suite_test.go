package mongodb_test

import (
	"context"
	"os"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/pkg/database/mongodb"
)

type MongoSuite struct {
	suite.Suite

	db *mongo.Database
}

func (s *MongoSuite) SetupSuite() {
	var (
		dbURI  = os.Getenv("TEST_DB_URI")
		dbName = os.Getenv("TEST_DB_NAME")
	)

	client, err := mongodb.NewClient(dbURI, "", "")
	s.Require().NoError(err)

	s.db = client.Database(dbName)
}

func (s *MongoSuite) TearDownSuite() {
	ctx := context.Background()

	err := s.db.Drop(ctx)
	s.Require().NoError(err)

	err = s.db.Client().Disconnect(ctx)
	s.Require().NoError(err)
}

func (s *MongoSuite) insertFlows(flows ...interface{}) {
	s.T().Helper()

	ctx := context.Background()

	for _, flow := range flows {
		_, err := s.db.Collection("flows").InsertOne(ctx, flow)
		s.Require().NoError(err)
	}
}

func (s *MongoSuite) insertPipelines(pipes ...interface{}) {
	s.T().Helper()

	ctx := context.Background()

	for _, pipe := range pipes {
		_, err := s.db.Collection("pipelines").InsertOne(ctx, pipe)
		s.Require().NoError(err)
	}
}

func (s *MongoSuite) insertSpecifications(specs ...interface{}) {
	s.T().Helper()

	ctx := context.Background()

	for _, spec := range specs {
		_, err := s.db.Collection("specifications").InsertOne(ctx, spec)
		s.Require().NoError(err)
	}
}

func (s *MongoSuite) insertTestCampaigns(tcs ...interface{}) {
	s.T().Helper()

	ctx := context.Background()

	for _, tc := range tcs {
		_, err := s.db.Collection("testCampaigns").InsertOne(ctx, tc)
		s.Require().NoError(err)
	}
}
