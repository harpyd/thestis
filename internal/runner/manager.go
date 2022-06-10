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
	"github.com/harpyd/thestis/pkg/correlationid"
	"github.com/harpyd/thestis/pkg/database/mongodb"
)

type Manager struct {
	zapSingleton
	mongoSingleton
	natsSingleton
	firebaseSingleton
	pipelineWPSingleton

	logger       service.Logger
	config       *config.Config
	persistent   persistentContext
	specParser   service.SpecificationParser
	metrics      metricsContext
	pipeline     pipelineContext
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

type pipelineWPSingleton struct {
	once sync.Once
	wp   *workerpool.WorkerPool
}

type pipelineContext struct {
	guard      service.PipelineGuard
	policy     service.PipelinePolicy
	maintainer service.PipelineMaintainer
	enqueuer   service.Enqueuer
}

type persistentContext struct {
	testCampaignRepo service.TestCampaignRepository
	specRepo         service.SpecificationRepository
	pipeRepo         service.PipelineRepository
	flowRepo         service.FlowRepository
	testCampaignRM   query.TestCampaignReadModel
	specificationRM  query.SpecificationReadModel
}

type signalBusContext struct {
	publisher  service.PipelineCancelPublisher
	subscriber service.PipelineCancelSubscriber
}

type metricsContext struct {
	httpMetric rest.MetricCollector
}

func New(configsPath string) *Manager {
	c := &Manager{}

	c.initLogger()
	c.initConfig(configsPath)
	c.initPersistent()
	c.initSpecificationParser()
	c.initMetrics()
	c.initSignalBus()
	c.initPipeline()
	c.initApplication()
	c.initAuthenticationProvider()
	c.initServer()

	return c
}

func (c *Manager) Start() {
	c.logger.Info("Runner started")

	c.logger.Info(
		"HTTP server started",
		"port", fmt.Sprintf(":%s", c.config.HTTP.Port),
	)

	if err := c.server.Start(); !errors.Is(err, http.ErrServerClosed) {
		c.logger.Fatal("HTTP server stopped unexpectedly", err)
	}
}

func (c *Manager) Stop() {
	defer c.syncZap()

	c.stopPipelineWorkerPool()
	c.logger.Info("Pipeline worker pool stopped")

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

func (c *Manager) shutdownServer() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		c.config.HTTP.ShutdownTimeout,
	)
	defer cancel()

	return c.server.Shutdown(ctx)
}

func (c *Manager) zap() *zap.Logger {
	c.zapSingleton.once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to create zap logger: %v", err)
		}

		c.zapSingleton.logger = logger
	})

	return c.zapSingleton.logger
}

func (c *Manager) syncZap() {
	if err := c.zap().Sync(); err != nil {
		log.Fatalf("Failed to sync zap logger: %v", err)
	}
}

func (c *Manager) mongo() *mongo.Database {
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

func (c *Manager) disconnectMongo() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		c.config.Mongo.DisconnectTimeout,
	)
	defer cancel()

	return c.mongo().Client().Disconnect(ctx)
}

func (c *Manager) nats() *nats.Conn {
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

func (c *Manager) disconnectNATS() {
	c.nats().Close()
}

func (c *Manager) pipelineWorkerPool() *workerpool.WorkerPool {
	c.pipelineWPSingleton.once.Do(func() {
		c.pipelineWPSingleton.wp = workerpool.New(10)

		c.logger.Info(
			"Pipeline worker pool initialization completed",
			"workers", c.config.Pipeline.Workers,
		)
	})

	return c.pipelineWPSingleton.wp
}

func (c *Manager) stopPipelineWorkerPool() {
	c.pipelineWorkerPool().StopWait()
}

func (c *Manager) firebaseAuth() *fireauth.Client {
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

func (c *Manager) initLogger() {
	c.logger = zapAdapter.NewLogger(c.zap())
}

func (c *Manager) initConfig(configsPath string) {
	cfg, err := config.FromPath(configsPath)
	if err != nil {
		c.logger.Fatal("Failed to parse config", err)
	}

	c.config = cfg

	c.logger.Info("Config parsing completed")
}

func (c *Manager) initPersistent() {
	db := c.mongo()
	args := []interface{}{"db", "mongo"}

	var (
		testCampaignRepo = mongoAdapter.NewTestCampaignRepository(db)
		specRepo         = mongoAdapter.NewSpecificationRepository(db)
		pipeRepo         = mongoAdapter.NewPipelineRepository(db)
		flowRepo         = mongoAdapter.NewFlowRepository(db)
	)

	c.persistent.testCampaignRepo = testCampaignRepo
	c.logger.Info("Test campaign repository initialization completed", args...)

	c.persistent.specRepo = specRepo
	c.logger.Info("Specification repository initialization completed", args...)

	c.persistent.pipeRepo = pipeRepo
	c.logger.Info("Pipeline repository initialization completed", args...)

	c.persistent.flowRepo = flowRepo
	c.logger.Info("Flow repository initialization completed", args...)

	c.persistent.testCampaignRM = testCampaignRepo
	c.logger.Info("Test campaign read model initialization completed", args...)

	c.persistent.specificationRM = specRepo
	c.logger.Info("Specification read model initialization completed", args...)
}

func (c *Manager) initSpecificationParser() {
	c.specParser = yaml.NewSpecificationParser()
	c.logger.Info("Specification parser service initialization completed", "type", "yaml")
}

func (c *Manager) initApplication() {
	c.app = &app.Application{
		Commands: app.Commands{
			CreateTestCampaign: command.NewCreateTestCampaignHandler(c.persistent.testCampaignRepo),
			LoadSpecification: command.NewLoadSpecificationHandler(
				c.persistent.specRepo,
				c.persistent.testCampaignRepo,
				c.specParser,
			),
			StartPipeline: command.NewStartPipelineHandler(
				c.persistent.specRepo,
				c.persistent.pipeRepo,
				c.pipeline.maintainer,
			),
			RestartPipeline: command.NewRestartPipelineHandler(
				c.persistent.pipeRepo,
				c.persistent.specRepo,
				c.pipeline.maintainer,
			),
			CancelPipeline: command.NewCancelPipelineHandler(c.persistent.pipeRepo, c.signalBus.publisher),
		},
		Queries: app.Queries{
			TestCampaign:  query.NewTestCampaignHandler(c.persistent.testCampaignRM),
			Specification: query.NewSpecificationHandler(c.persistent.specificationRM),
		},
	}

	c.logger.Info("Application context initialization completed")
}

func (c *Manager) initMetrics() {
	mrs, err := prometheus.NewMetricCollector()
	if err != nil {
		c.logger.Fatal("Failed to register metrics", err)
	}

	c.metrics.httpMetric = mrs

	c.logger.Info("Metrics registration completed", "db", "prometheus")
}

func (c *Manager) initSignalBus() {
	if c.config.Pipeline.SignalBus == config.Nats {
		bus := natsio.NewPipelineCancelSignalBus(c.nats())

		c.signalBus.publisher = bus
		c.signalBus.subscriber = bus
	} else {
		c.logger.Fatal(
			"Invalid pipeline signal bus",
			errors.Errorf("%s is not valid signal bus", c.config.Pipeline.SignalBus),
			"allowed", config.Nats,
		)
	}

	c.logger.Info(
		"Signal bus initialization completed",
		"signalBus", c.config.Pipeline.SignalBus,
	)
}

func (c *Manager) initPipeline() {
	c.initPipelineGuard()
	c.initPipelinePolicy()
	c.initEnqueuer()

	c.pipeline.maintainer = service.NewPipelineMaintainer(
		c.pipeline.guard,
		c.signalBus.subscriber,
		c.pipeline.policy,
		c.pipeline.enqueuer,
		c.logger.Named("PipelineMaintainer"),
		c.config.Pipeline.FlowTimeout,
	)

	c.logger.Info(
		"Pipeline maintainer initialized",
		"policy", c.config.Pipeline.Policy,
	)
}

func (c *Manager) initPipelineGuard() {
	c.pipeline.guard = mongoAdapter.NewPipelineGuard(c.mongo())
}

func (c *Manager) initPipelinePolicy() {
	if c.config.Pipeline.Policy == config.SavePerStepPolicy {
		c.pipeline.policy = service.NewSavePerStepPolicy(
			c.persistent.flowRepo,
			c.logger.Named("SavePerStepPolicy"),
			c.config.SavePerStep.SaveTimeout,
		)

		return
	}

	c.logger.Fatal(
		"Invalid pipeline steps policy",
		errors.Errorf("%s is not valid steps policy", c.config.Pipeline.Policy),
		"allowed", config.SavePerStepPolicy,
	)
}

func (c *Manager) initEnqueuer() {
	c.pipeline.enqueuer = service.EnqueueFunc(c.pipelineWorkerPool().Submit)
}

func (c *Manager) initAuthenticationProvider() {
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
			"allowed", strings.Join([]string{config.FakeAuth, config.FirebaseAuth}, ", "),
		)
	}

	c.logger.Info("Authentication provider initialization completed", "auth", authType)
}

func (c *Manager) initServer() {
	c.server = server.New(c.config.HTTP, rest.NewHandler(rest.Params{
		Middlewares: []rest.Middleware{
			correlationid.Middleware("X-Correlation-Id"),
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
