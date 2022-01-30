package runner

import (
	"fmt"
	stdhttp "net/http"
	"strings"

	fireauth "firebase.google.com/go/auth"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5/middleware"
	fakeAuth "github.com/harpyd/thestis/internal/adapter/auth/fake"
	firebaseAuth "github.com/harpyd/thestis/internal/adapter/auth/firebase"
	zapadap "github.com/harpyd/thestis/internal/adapter/logger/zap"
	"github.com/harpyd/thestis/internal/adapter/metrics/prometheus"
	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	mongoadap "github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	"github.com/harpyd/thestis/internal/app/query"
	"github.com/harpyd/thestis/internal/config"
	"github.com/harpyd/thestis/internal/port/http"
	"github.com/harpyd/thestis/internal/port/http/auth"
	"github.com/harpyd/thestis/internal/port/http/logging"
	"github.com/harpyd/thestis/internal/port/http/metrics"
	v1 "github.com/harpyd/thestis/internal/port/http/v1"
	"github.com/harpyd/thestis/internal/server"
	"github.com/harpyd/thestis/pkg/auth/firebase"
	"github.com/harpyd/thestis/pkg/database/mongodb"
	"github.com/harpyd/thestis/pkg/http/cors"
)

func Start(configsPath string) {
	newRunner(configsPath).start()
}

type runnerContext struct {
	logger       app.LoggingService
	config       *config.Config
	persistent   persistentContext
	specParser   app.SpecificationParserService
	metrics      app.MetricsService
	flowManager  app.FlowManager
	app          *app.Application
	authProvider auth.Provider
	server       *server.Server

	cancel func()
}

type persistentContext struct {
	testCampaignsRepo      app.TestCampaignsRepository
	specsRepo              app.SpecificationsRepository
	perfsRepo              app.PerformancesRepository
	flowsRepo              app.FlowsRepository
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
	c.initFlowManager()
	c.initApplication()
	c.initAuthenticationProvider()
	c.initServer()

	return c
}

func (c *runnerContext) start() {
	defer c.cancel()

	c.logger.Info(
		"HTTP server started",
		app.LogField{Key: "port", Value: fmt.Sprintf(":%s", c.config.HTTP.Port)},
	)

	err := c.server.Start()

	c.logger.Fatal("HTTP server stopped", err)
}

func (c *runnerContext) initLogger() func() {
	logger, _ := zap.NewProduction()
	sync := func() {
		_ = logger.Sync()
	}

	c.logger = zapadap.NewLoggingService(logger)

	return sync
}

func (c *runnerContext) initConfig(configsPath string) {
	cfg, err := config.FromPath(configsPath)
	if err != nil {
		c.logger.Fatal("Failed to parse config", err)
	}

	c.config = cfg

	c.logger.Info("Config parsing completed")
}

func (c *runnerContext) initPersistent() {
	db := c.mongoDatabase()
	logField := app.LogField{Key: "db", Value: "mongo"}

	var (
		testCampaignsRepo = mongoadap.NewTestCampaignsRepository(db)
		specsRepo         = mongoadap.NewSpecificationsRepository(db)
		perfsRepo         = mongoadap.NewPerformancesRepository(db)
		flowsRepo         = mongoadap.NewFlowsRepository(db)
	)

	c.persistent.testCampaignsRepo = testCampaignsRepo
	c.logger.Info("Test campaigns repository initialization completed", logField)

	c.persistent.specsRepo = specsRepo
	c.logger.Info("Specifications repository initialization completed", logField)

	c.persistent.perfsRepo = perfsRepo
	c.logger.Info("Performances repository initialization completed", logField)

	c.persistent.flowsRepo = flowsRepo
	c.logger.Info("Flows repository initialization completed", logField)

	c.persistent.specificTestCampaignRM = testCampaignsRepo
	c.logger.Info("Specific test campaigns read model initialization completed", logField)

	c.persistent.specificSpecRM = specsRepo
	c.logger.Info("Specific specifications read model initialization completed", logField)
}

func (c *runnerContext) initSpecificationParser() {
	c.specParser = yaml.NewSpecificationParserService()
	c.logger.Info("Specification parser service initialization completed", app.LogField{
		Key: "type", Value: "yaml",
	})
}

func (c *runnerContext) mongoDatabase() *mongo.Database {
	client, err := mongodb.NewClient(c.config.Mongo.URI, c.config.Mongo.Username, c.config.Mongo.Password)
	if err != nil {
		c.logger.Fatal("Failed to connect to MongoDB", err)
	}

	return client.Database(c.config.Mongo.DatabaseName)
}

func (c *runnerContext) initApplication() {
	c.app = &app.Application{
		Commands: app.Commands{
			CreateTestCampaign: command.NewCreateTestCampaignHandler(c.persistent.testCampaignsRepo),
			LoadSpecification: command.NewLoadSpecificationHandler(
				c.persistent.specsRepo,
				c.persistent.testCampaignsRepo,
				c.specParser,
			),
			StartNewPerformance: command.NewStartPerformanceHandler(
				c.flowManager,
				c.persistent.specsRepo,
				c.persistent.perfsRepo,
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
		c.logger.Fatal("Failed to register metrics", err)
	}

	c.metrics = mrs

	c.logger.Info("Metrics registration completed", app.LogField{Key: "db", Value: "prometheus"})
}

func (c *runnerContext) initFlowManager() {
	c.flowManager = app.NewEveryStepSavingFlowManager(c.persistent.perfsRepo, c.persistent.flowsRepo)
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
			errors.Errorf("%s is not valid auth type", authType),
			app.LogField{
				Key: "allowed",
				Value: strings.Join([]string{
					config.FakeAuth,
					config.FirebaseAuth,
				}, ", "),
			},
		)
	}

	c.logger.Info("Authentication provider initialization completed",
		app.LogField{Key: "auth", Value: authType},
	)
}

func (c *runnerContext) firebaseClient() *fireauth.Client {
	client, err := firebase.NewClient(c.config.Firebase.ServiceAccountFile)
	if err != nil {
		c.logger.Fatal("Failed to create Firebase Auth client", err)
	}

	return client
}

func (c *runnerContext) initServer() {
	c.server = server.New(c.config, http.NewHandler(http.Params{
		Middlewares: []http.Middleware{
			middleware.RequestID,
			middleware.RealIP,
			logging.Middleware(c.logger),
			middleware.Recoverer,
			cors.Middleware(c.config.HTTP.AllowedOrigins),
			middleware.NoCache,
			metrics.Middleware(c.metrics),
		},
		Routes: []http.Route{
			{
				Pattern: "/v1",
				Handler: v1.NewHandler(c.app, c.logger, auth.Middleware(c.authProvider)),
			},
			{
				Pattern: "/swagger",
				Handler: stdhttp.StripPrefix("/swagger/", stdhttp.FileServer(stdhttp.Dir("./swagger"))),
			},
			{
				Pattern: "/metrics",
				Handler: promhttp.Handler(),
			},
		},
	}))

	c.logger.Info("Server initializing completed")
}
