package runner

import (
	"fmt"
	"log"
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
	v1 "github.com/harpyd/thestis/internal/port/http/v1"
	"github.com/harpyd/thestis/internal/server"
	"github.com/harpyd/thestis/pkg/auth/firebase"
	"github.com/harpyd/thestis/pkg/database/mongodb"
)

type Context struct {
	zapSingleton
	mongoSingleton
	natsSingleton
	firebaseSingleton

	logger       app.Logger
	config       *config.Config
	persistent   persistentContext
	specParser   app.SpecificationParser
	metrics      metricsContext
	performance  performanceContext
	signalBus    signalBusContext
	app          *app.Application
	authProvider http.AuthProvider
	server       *server.Server
}

type mongoSingleton struct {
	once sync.Once
	db   *mongo.Database
}

type natsSingleton struct {
	once sync.Once
	conn *nats.Conn
}

type firebaseSingleton struct {
	once   sync.Once
	client *fireauth.Client
}

type zapSingleton struct {
	once   sync.Once
	logger *zap.Logger
}

type performanceContext struct {
	guard       app.PerformanceGuard
	stepsPolicy app.StepsPolicy
	maintainer  app.PerformanceMaintainer
}

type persistentContext struct {
	testCampaignRepo       app.TestCampaignRepository
	specRepo               app.SpecificationRepository
	perfRepo               app.PerformanceRepository
	flowRepo               app.FlowRepository
	specificTestCampaignRM app.SpecificTestCampaignReadModel
	specificSpecRM         app.SpecificSpecificationReadModel
}

type signalBusContext struct {
	publisher  app.PerformanceCancelPublisher
	subscriber app.PerformanceCancelSubscriber
}

type metricsContext struct {
	httpMetric http.MetricCollector
}

func New(configsPath string) *Context {
	c := &Context{}

	c.initLogger()
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

func (c *Context) Start() {
	c.logger.Info(
		"HTTP server started",
		app.StringLogField("port", fmt.Sprintf(":%s", c.config.HTTP.Port)),
	)

	if err := c.server.Start(); !errors.Is(err, stdhttp.ErrServerClosed) {
		c.logger.Fatal("HTTP server stopped unexpectedly", err)
	}

	c.logger.Info("HTTP server stopped")
}

func (c *Context) Stop() {
	if err := c.server.Shutdown(); err != nil {
		c.logger.Fatal("Server shutdown failed", err)
	}

	if err := c.zapSingleton.logger.Sync(); err != nil {
		log.Fatal("Failed to sync zap logger")
	}
}

func (c *Context) zapLogger() *zap.Logger {
	c.zapSingleton.once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatal("Failed to initialize zap logger")
		}

		c.zapSingleton.logger = logger
	})

	return c.zapSingleton.logger
}

func (c *Context) mongoDatabase() *mongo.Database {
	c.mongoSingleton.once.Do(func() {
		client, err := mongodb.NewClient(
			c.config.Mongo.URI,
			c.config.Mongo.Username,
			c.config.Mongo.Password,
		)
		if err != nil {
			c.logger.Fatal("Failed to connect to MongoDB", err)
		}

		c.mongoSingleton.db = client.Database(c.config.Mongo.DatabaseName)

		c.logger.Info("Connected to MongoDB")
	})

	return c.mongoSingleton.db
}

func (c *Context) natsConnection() *nats.Conn {
	c.natsSingleton.once.Do(func() {
		conn, err := nats.Connect(c.config.Nats.URL)
		if err != nil {
			c.logger.Fatal("Failed to connect to NATS server", err)
		}

		c.natsSingleton.conn = conn

		c.logger.Info("Connected to NATS server")
	})

	return c.natsSingleton.conn
}

func (c *Context) firebaseClient() *fireauth.Client {
	c.firebaseSingleton.once.Do(func() {
		client, err := firebase.NewClient(c.config.Firebase.ServiceAccountFile)
		if err != nil {
			c.logger.Fatal("Failed to create Firebase Auth client", err)
		}

		c.firebaseSingleton.client = client

		c.logger.Info("Firebase Auth client created")
	})

	return c.firebaseSingleton.client
}

func (c *Context) initLogger() {
	c.logger = zapAdapter.NewLogger(c.zapLogger())
}

func (c *Context) initConfig(configsPath string) {
	cfg, err := config.FromPath(configsPath)
	if err != nil {
		c.logger.Fatal("Failed to parse config", err)
	}

	c.config = cfg

	c.logger.Info("Config parsing completed")
}

func (c *Context) initPersistent() {
	db := c.mongoDatabase()
	logField := app.StringLogField("db", "mongo")

	var (
		testCampaignRepo = mongoAdapter.NewTestCampaignRepository(db)
		specRepo         = mongoAdapter.NewSpecificationRepository(db)
		perfRepo         = mongoAdapter.NewPerformanceRepository(db)
		flowRepo         = mongoAdapter.NewFlowRepository(db)
	)

	c.persistent.testCampaignRepo = testCampaignRepo
	c.logger.Info("Test campaign repository initialization completed", logField)

	c.persistent.specRepo = specRepo
	c.logger.Info("Specification repository initialization completed", logField)

	c.persistent.perfRepo = perfRepo
	c.logger.Info("Performance repository initialization completed", logField)

	c.persistent.flowRepo = flowRepo
	c.logger.Info("Flow repository initialization completed", logField)

	c.persistent.specificTestCampaignRM = testCampaignRepo
	c.logger.Info("Specific test campaign read model initialization completed", logField)

	c.persistent.specificSpecRM = specRepo
	c.logger.Info("Specific specification read model initialization completed", logField)
}

func (c *Context) initSpecificationParser() {
	c.specParser = yaml.NewSpecificationParser()
	c.logger.Info("Specification parser service initialization completed", app.StringLogField("type", "yaml"))
}

func (c *Context) initApplication() {
	c.app = &app.Application{
		Commands: app.Commands{
			CreateTestCampaign: command.NewCreateTestCampaignHandler(c.persistent.testCampaignRepo),
			LoadSpecification: command.NewLoadSpecificationHandler(
				c.persistent.specRepo,
				c.persistent.testCampaignRepo,
				c.specParser,
			),
			StartPerformance: command.NewStartPerformanceHandler(
				c.persistent.specRepo,
				c.persistent.perfRepo,
				c.performance.maintainer,
			),
			RestartPerformance: command.NewRestartPerformanceHandler(
				c.persistent.perfRepo,
				c.persistent.specRepo,
				c.performance.maintainer,
			),
			CancelPerformance: command.NewCancelPerformanceHandler(c.persistent.perfRepo, c.signalBus.publisher),
		},
		Queries: app.Queries{
			SpecificTestCampaign:  query.NewSpecificTestCampaignHandler(c.persistent.specificTestCampaignRM),
			SpecificSpecification: query.NewSpecificSpecificationHandler(c.persistent.specificSpecRM),
		},
	}

	c.logger.Info("Application context initialization completed")
}

func (c *Context) initMetrics() {
	mrs, err := prometheus.NewMetricCollector()
	if err != nil {
		c.logger.Fatal("Failed to register metrics", err)
	}

	c.metrics.httpMetric = mrs

	c.logger.Info("Metrics registration completed", app.StringLogField("db", "prometheus"))
}

func (c *Context) initSignalBus() {
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

func (c *Context) initPerformance() {
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

func (c *Context) initPerformanceGuard() {
	c.performance.guard = mongoAdapter.NewPerformanceGuard(c.mongoDatabase())
}

func (c *Context) initStepsPolicy() {
	if c.config.Performance.Policy == config.EveryStepSavingPolicy {
		c.performance.stepsPolicy = app.NewEveryStepSavingPolicy(
			c.persistent.flowRepo,
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

func (c *Context) initAuthenticationProvider() {
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

func (c *Context) initServer() {
	c.server = server.New(c.config.HTTP, http.NewHandler(http.Params{
		Middlewares: []http.Middleware{
			middleware.RequestID,
			middleware.RealIP,
			http.LoggingMiddleware(c.logger),
			middleware.Recoverer,
			http.CORSMiddleware(c.config.HTTP.AllowedOrigins),
			middleware.NoCache,
			http.MetricsMiddleware(c.metrics.httpMetric),
		},
		Routes: []http.Route{
			{
				Pattern: "/v1",
				Handler: v1.NewHandler(c.app, c.logger, http.AuthMiddleware(c.authProvider)),
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
