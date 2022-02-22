package runner

import (
	"fmt"
	stdhttp "net/http"
	"strings"
	"sync"

	fireauth "firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	fakeAdapter "github.com/harpyd/thestis/internal/adapter/auth/fake"
	firebaseAdapter "github.com/harpyd/thestis/internal/adapter/auth/firebase"
	zapAdapter "github.com/harpyd/thestis/internal/adapter/logger/zap"
	"github.com/harpyd/thestis/internal/adapter/metrics/prometheus"
	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	mongoAdapter "github.com/harpyd/thestis/internal/adapter/persistence/mongodb"
	"github.com/harpyd/thestis/internal/adapter/pubsub/natsio"
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
	mongoSingletone
	natsSingletone
	firebaseSingletone

	logger       app.LoggingService
	config       *config.Config
	persistent   persistentContext
	specParser   app.SpecificationParserService
	metrics      app.MetricsService
	performance  performanceContext
	signalBus    signalBusContext
	app          *app.Application
	authProvider auth.Provider
	server       *server.Server

	cancel func()
}

type mongoSingletone struct {
	once sync.Once
	db   *mongo.Database
}

type natsSingletone struct {
	once sync.Once
	conn *nats.Conn
}

type firebaseSingletone struct {
	once   sync.Once
	client *fireauth.Client
}

type performanceContext struct {
	guard       app.PerformanceGuard
	stepsPolicy app.StepsPolicy
	maintainer  app.PerformanceMaintainer
}

type persistentContext struct {
	testCampaignsRepo      app.TestCampaignsRepository
	specsRepo              app.SpecificationsRepository
	perfsRepo              app.PerformancesRepository
	flowsRepo              app.FlowsRepository
	specificTestCampaignRM app.SpecificTestCampaignReadModel
	specificSpecRM         app.SpecificSpecificationReadModel
}

type signalBusContext struct {
	publisher  app.PerformanceCancelPublisher
	subscriber app.PerformanceCancelSubscriber
}

func newRunner(configsPath string) *runnerContext {
	c := &runnerContext{}

	c.cancel = c.initLogger()
	c.initConfig(configsPath)
	c.initPersistent()
	c.initSpecificationParser()
	c.initMetrics()
	c.initSignalBus()
	c.initPerformance()
	c.initApplication()
	c.initAuthenticationProvider()
	c.initServer()

	return c
}

func (c *runnerContext) mongoDatabase() *mongo.Database {
	c.mongoSingletone.once.Do(func() {
		client, err := mongodb.NewClient(
			c.config.Mongo.URI,
			c.config.Mongo.Username,
			c.config.Mongo.Password,
		)
		if err != nil {
			c.logger.Fatal("Failed to connect to MongoDB", err)
		}

		c.mongoSingletone.db = client.Database(c.config.Mongo.DatabaseName)

		c.logger.Info("Connected to MongoDB")
	})

	return c.mongoSingletone.db
}

func (c *runnerContext) natsConnection() *nats.Conn {
	c.natsSingletone.once.Do(func() {
		conn, err := nats.Connect(c.config.Nats.URL)
		if err != nil {
			c.logger.Fatal("Failed to connect to Nats server", err)
		}

		c.natsSingletone.conn = conn

		c.logger.Info("Connected to Nats")
	})

	return c.natsSingletone.conn
}

func (c *runnerContext) firebaseClient() *fireauth.Client {
	c.firebaseSingletone.once.Do(func() {
		client, err := firebase.NewClient(c.config.Firebase.ServiceAccountFile)
		if err != nil {
			c.logger.Fatal("Failed to create Firebase Auth client", err)
		}

		c.firebaseSingletone.client = client

		c.logger.Info("Firebase Auth client created")
	})

	return c.firebaseSingletone.client
}

func (c *runnerContext) start() {
	defer c.cancel()

	c.logger.Info(
		"HTTP server started",
		app.StringLogField("port", fmt.Sprintf(":%s", c.config.HTTP.Port)),
	)

	err := c.server.Start()

	c.logger.Fatal("HTTP server stopped", err)
}

func (c *runnerContext) initLogger() func() {
	logger, _ := zap.NewProduction()
	sync := func() {
		_ = logger.Sync()
	}

	c.logger = zapAdapter.NewLoggingService(logger)

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
	logField := app.StringLogField("db", "mongo")

	var (
		testCampaignsRepo = mongoAdapter.NewTestCampaignsRepository(db)
		specsRepo         = mongoAdapter.NewSpecificationsRepository(db)
		perfsRepo         = mongoAdapter.NewPerformancesRepository(db)
		flowsRepo         = mongoAdapter.NewFlowsRepository(db)
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
	c.logger.Info("Specification parser service initialization completed", app.StringLogField("type", "yaml"))
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
			StartPerformance: command.NewStartPerformanceHandler(
				c.persistent.specsRepo,
				c.persistent.perfsRepo,
				c.performance.maintainer,
			),
			RestartPerformance: command.NewRestartPerformanceHandler(c.persistent.perfsRepo, c.performance.maintainer),
			CancelPerformance:  command.NewCancelPerformanceHandler(c.persistent.perfsRepo, c.signalBus.publisher),
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

	c.logger.Info("Metrics registration completed", app.StringLogField("db", "prometheus"))
}

func (c *runnerContext) initSignalBus() {
	if c.config.Performance.SignalBus == config.Nats {
		bus := natsio.NewPerformanceCancelSignalBus(c.natsConnection())

		c.signalBus.publisher = bus
		c.signalBus.subscriber = bus
	} else {
		c.logger.Fatal(
			"Invalid performance signal bus",
			errors.Errorf("%s is not valid signal bus", c.config.Performance.SignalBus),
			app.StringLogField("allowed", config.Nats),
		)
	}

	c.logger.Info(
		"Signal bus initialization completed",
		app.StringLogField("signalBus", c.config.Performance.SignalBus),
	)
}

func (c *runnerContext) initPerformance() {
	c.initPerformanceGuard()
	c.initStepsPolicy()

	c.performance.maintainer = app.NewPerformanceMaintainer(
		c.performance.guard,
		c.signalBus.subscriber,
		c.performance.stepsPolicy,
		c.config.Performance.FlowTimeout,
	)

	c.logger.Info(
		"Performance maintainer initialized",
		app.StringLogField("stepsPolicy", c.config.Performance.Policy),
	)
}

func (c *runnerContext) initPerformanceGuard() {
	c.performance.guard = mongoAdapter.NewPerformanceGuard(c.mongoDatabase())
}

func (c *runnerContext) initStepsPolicy() {
	if c.config.Performance.Policy == config.EveryStepSavingPolicy {
		c.performance.stepsPolicy = app.NewEveryStepSavingPolicy(
			c.persistent.flowsRepo,
			c.config.EveryStepSaving.SaveTimeout,
		)

		return
	}

	c.logger.Fatal(
		"Invalid performance steps policy",
		errors.Errorf("%s is not valid steps policy", c.config.Performance.Policy),
		app.StringLogField("allowed", config.EveryStepSavingPolicy),
	)
}

func (c *runnerContext) initAuthenticationProvider() {
	authType := c.config.Auth.With

	switch authType {
	case config.FakeAuth:
		c.authProvider = fakeAdapter.NewProvider()
	case config.FirebaseAuth:
		c.authProvider = firebaseAdapter.NewProvider(c.firebaseClient())
	default:
		c.logger.Fatal(
			"Invalid auth type",
			errors.Errorf("%s is not valid auth type", authType),
			app.StringLogField("allowed", strings.Join([]string{config.FakeAuth, config.FirebaseAuth}, ", ")),
		)
	}

	c.logger.Info("Authentication provider initialization completed", app.StringLogField("auth", authType))
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
