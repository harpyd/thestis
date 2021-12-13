package mongodb_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/pkg/database/mongodb"
)

type MongoTestFixtures struct {
	t *testing.T

	db *mongo.Database
}

func (f *MongoTestFixtures) SetupSuite() {
	var (
		dbURI  = os.Getenv("TEST_DB_URI")
		dbName = os.Getenv("TEST_DB_NAME")
	)

	client, err := mongodb.NewClient(dbURI, "", "")
	require.NoError(f.t, err)

	f.db = client.Database(dbName)
}

func (f *MongoTestFixtures) TearDownSuite() {
	ctx := context.Background()

	err := f.db.Drop(ctx)
	require.NoError(f.t, err)

	err = f.db.Client().Disconnect(ctx)
	require.NoError(f.t, err)
}
