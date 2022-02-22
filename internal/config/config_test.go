package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/config"
)

const fixturesPath = "./fixtures"

func TestFromPath(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	type env struct {
		HTTPPort           string
		MongoURI           string
		MongoDatabase      string
		ServiceAccountFile string
	}

	setEnv := func(env env) {
		_ = os.Setenv("HTTP_PORT", env.HTTPPort)
		_ = os.Setenv("MONGO_URI", env.MongoURI)
		_ = os.Setenv("MONGO_DATABASE", env.MongoDatabase)
		_ = os.Setenv("SERVICE_ACCOUNT_FILE", env.ServiceAccountFile)
	}

	testCases := []struct {
		Name           string
		ConfigsPath    string
		Env            env
		ExpectedConfig *config.Config
		ShouldBeErr    bool
		IsErr          func(err error) bool
	}{
		{
			Name:        "valid_test_config_with_env",
			ConfigsPath: fixturesPath,
			Env: env{
				HTTPPort:           "8080",
				MongoURI:           "some://uri",
				MongoDatabase:      "someName",
				ServiceAccountFile: "path/to/serviceAccount.json",
			},
			ExpectedConfig: &config.Config{
				Environment: "local",
				HTTP: config.HTTP{
					Port:           "8080",
					ReadTimeout:    8 * time.Second,
					WriteTimeout:   10 * time.Second,
					AllowedOrigins: []string{"https://some-a.com", "https://some-b.com"},
				},
				Mongo: config.Mongo{
					URI:          "some://uri",
					DatabaseName: "someName",
					Username:     "admin",
					Password:     "0000",
				},
				Firebase: config.Firebase{
					ServiceAccountFile: "path/to/serviceAccount.json",
				},
				Auth: config.Auth{
					With: "fake",
				},
				Performance: config.Performance{
					FlowTimeout: 24 * time.Hour,
					Policy:      config.EveryStepSavingPolicy,
					SignalBus:   config.Nats,
				},
				EveryStepSaving: config.EveryStepSaving{
					SaveTimeout: 30 * time.Second,
				},
				Nats: config.NatsServer{
					URL: "nats://nats:4222",
				},
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			setEnv(c.Env)

			cfg, err := config.FromPath(c.ConfigsPath)
			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.ExpectedConfig, cfg)
		})
	}
}
