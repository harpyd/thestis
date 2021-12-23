package runner

import (
	"fmt"
	stdhttp "net/http"
	"strings"

	fireauth "firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	fakeAuth "github.com/harpyd/thestis/internal/adapter/auth/fake"
	firebaseAuth "github.com/harpyd/thestis/internal/adapter/auth/firebase"
	"github.com/harpyd/thestis/internal/adapter/metrics/prometheus"
	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	mongorepo "github.com/harpyd/thestis/internal/adapter/repository/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	"github.com/harpyd/thestis/internal/app/query"
	"github.com/harpyd/thestis/internal/config"
	"github.com/harpyd/thestis/internal/port/http"
	v1 "github.com/harpyd/thestis/internal/port/http/v1"
	"github.com/harpyd/thestis/internal/server"
	"github.com/harpyd/thestis/pkg/auth/firebase"
	"github.com/harpyd/thestis/pkg/database/mongodb"
	"github.com/harpyd/thestis/pkg/http/auth"
)

func Start(configsPath string) {
	newRunner(configsPath).start()
}

type runnerContext struct {
	logger              *zap.Logger
	config              *config.Config
	persistent          persistentContext
	specificationParser app.SpecificationParserService
	metrics             app.MetricsService
	app                 *app.Application
	authProvider        auth.Provider
	server              *server.Server

	cancel func()
}

type persistentContext struct {
	testCampaignRepo       app.TestCampaignsRepository
	specRepo               app.SpecificationsRepository
	specificTestCampaignRM app.SpecificTestCampaignReadModel
	specificSpecRM         app.SpecificSpecificationReadModel
}

func newRunner(configsPath string) *runnerContext {
	c := &runnerContext{}

	c.cancel = c.initLogger()
	c.initConfig(configsPath)
	c.initPersistent()
	c.initSpecificationParser()
	c.initMetrics()
	c.initApplication()
	c.initAuthenticationProvider()
	c.initServer()

	return c
}

func (c *runnerContext) start() {
	defer c.cancel()

	c.logger.Info(
		"HTTP server started",
		zap.String("port", fmt.Sprintf(":%s", c.config.HTTP.Port)),
	)

	err := c.server.Start()

	c.logger.Fatal("HTTP server stopped", zap.Error(err))
}

func (c *runnerContext) initLogger() func() {
	logger, _ := zap.NewProduction()
	sync := func() {
		_ = logger.Sync()
	}

	c.logger = logger

	return sync
}

func (c *runnerContext) initConfig(configsPath string) {
	cfg, err := config.FromPath(configsPath)
	if err != nil {
		c.logger.Fatal("Failed to parse config", zap.Error(err))
	}

	c.config = cfg

	c.logger.Info("Config parsing completed")
}

func (c *runnerContext) initPersistent() {
	db := c.mongoDatabase()
	logField := zap.String("db", "mongo")

	var (
		testCampaignRepo = mongorepo.NewTestCampaignsRepository(db)
		specRepo         = mongorepo.NewSpecificationsRepository(db)
	)

	c.persistent.testCampaignRepo = testCampaignRepo
	c.logger.Info("Test campaigns repository initialization completed", logField)

	c.persistent.specRepo = specRepo
	c.logger.Info("Specifications repository initialization completed", logField)

	c.persistent.specificTestCampaignRM = testCampaignRepo
	c.logger.Info("Specific test campaigns read model initialization completed", logField)

	c.persistent.specificSpecRM = specRepo
	c.logger.Info("Specific specifications read model initialization completed", logField)
}

func (c *runnerContext) initSpecificationParser() {
	c.specificationParser = yaml.NewSpecificationParserService()
	c.logger.Info("Specification parser service initialization completed", zap.String("type", "yaml"))
}

func (c *runnerContext) mongoDatabase() *mongo.Database {
	client, err := mongodb.NewClient(c.config.Mongo.URI, c.config.Mongo.Username, c.config.Mongo.Password)
	if err != nil {
		c.logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	return client.Database(c.config.Mongo.DatabaseName)
}

func (c *runnerContext) initApplication() {
	c.app = &app.Application{
		Commands: app.Commands{
			CreateTestCampaign: command.NewCreateTestCampaignHandler(c.persistent.testCampaignRepo),
			LoadSpecification: command.NewLoadSpecificationHandler(
				c.persistent.specRepo,
				c.persistent.testCampaignRepo,
				c.specificationParser,
			),
		},
		Queries: app.Queries{
			SpecificTestCampaign:  query.NewSpecificTestCampaignHandler(c.persistent.specificTestCampaignRM),
			SpecificSpecification: query.NewSpecificSpecificationHandler(c.persistent.specificSpecRM),
		},
	}

	c.logger.Info("Application context initialization completed")
}

func (c *runnerContext) initMetrics() {
	mrs, err := prometheus.NewMetricsService()
	if err != nil {
		c.logger.Fatal("Failed to register metrics", zap.Error(err))
	}

	c.metrics = mrs

	c.logger.Info("Metrics registration completed", zap.String("db", "prometheus"))
}

func (c *runnerContext) initAuthenticationProvider() {
	authType := c.config.Auth.With

	switch authType {
	case config.FakeAuth:
		c.authProvider = fakeAuth.NewProvider()
	case config.FirebaseAuth:
		c.authProvider = firebaseAuth.NewProvider(c.firebaseClient())
	default:
		c.logger.Fatal(
			"Invalid auth type",
			zap.String("actual", authType),
			zap.String("allowed", strings.Join([]string{
				config.FakeAuth,
				config.FirebaseAuth,
			}, ", ")),
		)
	}

	c.logger.Info("Authentication provider initialization completed", zap.String("auth", authType))
}

func (c *runnerContext) firebaseClient() *fireauth.Client {
	client, err := firebase.NewClient(c.config.Firebase.ServiceAccountFile)
	if err != nil {
		c.logger.Fatal("Failed to create Firebase Auth client", zap.Error(err))
	}

	return client
}

func (c *runnerContext) initServer() {
	c.server = server.New(c.config, http.NewHandler(
		c.logger,
		http.Route{
			Pattern: "/v1",
			Handler: v1.NewHandler(c.app, c.v1Router()),
		},
		http.Route{
			Pattern: "/swagger",
			Handler: stdhttp.StripPrefix("/swagger/", stdhttp.FileServer(stdhttp.Dir("./swagger"))),
		},
		http.Route{
			Pattern: "/metrics",
			Handler: promhttp.Handler(),
		},
	))

	c.logger.Info("Server initializing completed")
}

func (c *runnerContext) v1Router() chi.Router {
	r := chi.NewRouter()
	r.Use(
		auth.Middleware(c.authProvider),
	)

	return r
}
