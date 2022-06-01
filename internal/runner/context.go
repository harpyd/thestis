package runner

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	fireauth "firebase.google.com/go/auth"
	"github.com/gammazero/workerpool"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/harpyd/thestis/internal/config"
	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/query"
	"github.com/harpyd/thestis/internal/core/app/service"
	fakeAdapter "github.com/harpyd/thestis/internal/core/infrastructure/auth/fake"
	firebaseAdapter "github.com/harpyd/thestis/internal/core/infrastructure/auth/firebase"
	zapAdapter "github.com/harpyd/thestis/internal/core/infrastructure/logger/zap"
	"github.com/harpyd/thestis/internal/core/infrastructure/metrics/prometheus"
	"github.com/harpyd/thestis/internal/core/infrastructure/parser/yaml"
	mongoAdapter "github.com/harpyd/thestis/internal/core/infrastructure/persistence/mongodb"
	"github.com/harpyd/thestis/internal/core/infrastructure/pubsub/natsio"
	"github.com/harpyd/thestis/internal/core/interface/rest"
	v1 "github.com/harpyd/thestis/internal/core/interface/rest/v1"
	"github.com/harpyd/thestis/internal/server"
	"github.com/harpyd/thestis/pkg/auth/firebase"
	"github.com/harpyd/thestis/pkg/database/mongodb"
)

type Context struct {
	zapSingleton
	mongoSingleton
	natsSingleton
	firebaseSingleton
	performanceWPSingleton

	logger       service.Logger
	config       *config.Config
	persistent   persistentContext
	specParser   service.SpecificationParser
	metrics      metricsContext
	performance  performanceContext
	signalBus    signalBusContext
	app          *app.Application
	authProvider rest.AuthProvider
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

type performanceWPSingleton struct {
	once sync.Once
	wp   *workerpool.WorkerPool
}

type performanceContext struct {
	guard      service.PerformanceGuard
	policy     service.PerformancePolicy
	maintainer service.PerformanceMaintainer
	enqueuer   service.Enqueuer
}

type persistentContext struct {
	testCampaignRepo       service.TestCampaignRepository
	specRepo               service.SpecificationRepository
	perfRepo               service.PerformanceRepository
	flowRepo               service.FlowRepository
	specificTestCampaignRM query.SpecificTestCampaignReadModel
	specificSpecRM         query.SpecificSpecificationReadModel
}

type signalBusContext struct {
	publisher  service.PerformanceCancelPublisher
	subscriber service.PerformanceCancelSubscriber
}

type metricsContext struct {
	httpMetric rest.MetricCollector
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
	c.logger.Info("Runner started")

	c.logger.Info(
		"HTTP server started",
		service.StringLogField("port", fmt.Sprintf(":%s", c.config.HTTP.Port)),
	)

	if err := c.server.Start(); !errors.Is(err, http.ErrServerClosed) {
		c.logger.Fatal("HTTP server stopped unexpectedly", err)
	}
}

func (c *Context) Stop() {
	defer c.syncZap()

	c.stopPerformanceWorkerPool()
	c.logger.Info("Performance worker pool stopped")

	err := multierr.Append(
		c.shutdownServer(),
		c.disconnectMongo(),
	)

	c.logger.Info("Server shutdown succeeded")
	c.logger.Info("Mongo disconnected")

	c.disconnectNATS()
	c.logger.Info("NATS disconnected")

	if err != nil {
		c.logger.Fatal("Runner stopped incorrectly", err)
	}

	c.logger.Info("Runner stopped")
}

func (c *Context) shutdownServer() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		c.config.HTTP.ShutdownTimeout,
	)
	defer cancel()

	return c.server.Shutdown(ctx)
}

func (c *Context) zap() *zap.Logger {
	c.zapSingleton.once.Do(func() {
		logger, err := zap.NewProduction(zap.AddCallerSkip(1))
		if err != nil {
			log.Fatalf("Failed to create zap logger: %v", err)
		}

		c.zapSingleton.logger = logger
	})

	return c.zapSingleton.logger
}

func (c *Context) syncZap() {
	if err := c.zap().Sync(); err != nil {
		log.Fatalf("Failed to sync zap logger: %v", err)
	}
}

func (c *Context) mongo() *mongo.Database {
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

func (c *Context) disconnectMongo() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		c.config.Mongo.DisconnectTimeout,
	)
	defer cancel()

	return c.mongo().Client().Disconnect(ctx)
}

func (c *Context) nats() *nats.Conn {
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

func (c *Context) disconnectNATS() {
	c.nats().Close()
}

func (c *Context) performanceWorkerPool() *workerpool.WorkerPool {
	c.performanceWPSingleton.once.Do(func() {
		c.performanceWPSingleton.wp = workerpool.New(10)

		c.logger.Info(
			"Performance worker pool initialization completed",
			service.IntLogField("workers", c.config.Performance.Workers),
		)
	})

	return c.performanceWPSingleton.wp
}

func (c *Context) stopPerformanceWorkerPool() {
	c.performanceWorkerPool().StopWait()
}

func (c *Context) firebaseAuth() *fireauth.Client {
	c.firebaseSingleton.once.Do(func() {
		client, err := firebase.NewClient(c.config.Firebase.ServiceAccountFile)
		if err != nil {
			c.logger.Fatal("Failed to connect to Firebase Auth", err)
		}

		c.firebaseSingleton.client = client

		c.logger.Info("Connected to Firebase Auth")
	})

	return c.firebaseSingleton.client
}

func (c *Context) initLogger() {
	c.logger = zapAdapter.NewLogger(c.zap())
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
	db := c.mongo()
	logField := service.StringLogField("db", "mongo")

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
	c.logger.Info("Specification parser service initialization completed", service.StringLogField("type", "yaml"))
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

	c.logger.Info("Metrics registration completed", service.StringLogField("db", "prometheus"))
}

func (c *Context) initSignalBus() {
	if c.config.Performance.SignalBus == config.Nats {
		bus := natsio.NewPerformanceCancelSignalBus(c.nats())

		c.signalBus.publisher = bus
		c.signalBus.subscriber = bus
	} else {
		c.logger.Fatal(
			"Invalid performance signal bus",
			errors.Errorf("%s is not valid signal bus", c.config.Performance.SignalBus),
			service.StringLogField("allowed", config.Nats),
		)
	}

	c.logger.Info(
		"Signal bus initialization completed",
		service.StringLogField("signalBus", c.config.Performance.SignalBus),
	)
}

func (c *Context) initPerformance() {
	c.initPerformanceGuard()
	c.initPerformancePolicy()
	c.initEnqueuer()

	c.performance.maintainer = service.NewPerformanceMaintainer(
		c.performance.guard,
		c.signalBus.subscriber,
		c.performance.policy,
		c.performance.enqueuer,
		c.config.Performance.FlowTimeout,
	)

	c.logger.Info(
		"Performance maintainer initialized",
		service.StringLogField("policy", c.config.Performance.Policy),
	)
}

func (c *Context) initPerformanceGuard() {
	c.performance.guard = mongoAdapter.NewPerformanceGuard(c.mongo())
}

func (c *Context) initPerformancePolicy() {
	if c.config.Performance.Policy == config.SavePerStepPolicy {
		c.performance.policy = service.NewSavePerStepPolicy(
			c.persistent.flowRepo,
			c.config.SavePerStep.SaveTimeout,
		)

		return
	}

	c.logger.Fatal(
		"Invalid performance steps policy",
		errors.Errorf("%s is not valid steps policy", c.config.Performance.Policy),
		service.StringLogField("allowed", config.SavePerStepPolicy),
	)
}

func (c *Context) initEnqueuer() {
	c.performance.enqueuer = service.EnqueueFunc(c.performanceWorkerPool().Submit)
}

func (c *Context) initAuthenticationProvider() {
	authType := c.config.Auth.With

	switch authType {
	case config.FakeAuth:
		c.authProvider = fakeAdapter.NewProvider()
	case config.FirebaseAuth:
		c.authProvider = firebaseAdapter.NewProvider(c.firebaseAuth())
	default:
		c.logger.Fatal(
			"Invalid auth type",
			errors.Errorf("%s is not valid auth type", authType),
			service.StringLogField("allowed", strings.Join([]string{config.FakeAuth, config.FirebaseAuth}, ", ")),
		)
	}

	c.logger.Info("Authentication provider initialization completed", service.StringLogField("auth", authType))
}

func (c *Context) initServer() {
	c.server = server.New(c.config.HTTP, rest.NewHandler(rest.Params{
		Middlewares: []rest.Middleware{
			middleware.RequestID,
			middleware.RealIP,
			rest.LoggingMiddleware(c.logger),
			middleware.Recoverer,
			rest.CORSMiddleware(c.config.HTTP.AllowedOrigins),
			middleware.NoCache,
			rest.MetricsMiddleware(c.metrics.httpMetric),
		},
		Routes: []rest.Route{
			{
				Pattern: "/v1",
				Handler: v1.NewHandler(c.app, c.logger, rest.AuthMiddleware(c.authProvider)),
			},
			{
				Pattern: "/swagger",
				Handler: http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))),
			},
			{
				Pattern: "/metrics",
				Handler: promhttp.Handler(),
			},
		},
	}))

	c.logger.Info("Server initializing completed")
}
