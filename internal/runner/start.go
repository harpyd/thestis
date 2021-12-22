package runner

import (
	"fmt"
	stdhttp "net/http"
	"strings"

	fireauth "firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

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
	logger       *zap.Logger
	config       *config.Config
	mongoDB      *mongo.Database
	firebaseAuth *fireauth.Client
	app          *app.Application
	middlewares  []func(stdhttp.Handler) stdhttp.Handler
	server       *server.Server

	cancel func()
}

func newRunner(configsPath string) *runnerContext {
	c := &runnerContext{}

	c.cancel = c.initLogger()
	c.initConfig(configsPath)
	c.initMongoDatabase()
	c.initApplication()
	c.initMiddlewares()
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

func (c *runnerContext) initMongoDatabase() {
	client, err := mongodb.NewClient(c.config.Mongo.URI, c.config.Mongo.Username, c.config.Mongo.Password)
	if err != nil {
		c.logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	c.mongoDB = client.Database(c.config.Mongo.DatabaseName)

	c.logger.Info("MongoDB connection completed")
}

func (c *runnerContext) initFirebaseClient() {
	firebaseAuth, err := firebase.NewClient(c.config.Firebase.ServiceAccountFile)
	if err != nil {
		c.logger.Fatal("Failed to create Firebase Auth client", zap.Error(err))
	}

	c.firebaseAuth = firebaseAuth

	c.logger.Info("Firebase Auth client creation completed")
}

func (c *runnerContext) initApplication() {
	tcRepo := mongorepo.NewTestCampaignsRepository(c.mongoDB)
	specRepo := mongorepo.NewSpecificationsRepository(c.mongoDB)
	parserService := yaml.NewSpecificationParserService()

	c.app = &app.Application{
		Commands: app.Commands{
			CreateTestCampaign: command.NewCreateTestCampaignHandler(tcRepo),
			LoadSpecification:  command.NewLoadSpecificationHandler(specRepo, tcRepo, parserService),
		},
		Queries: app.Queries{
			SpecificTestCampaign:  query.NewSpecificTestCampaignHandler(tcRepo),
			SpecificSpecification: query.NewSpecificSpecificationHandler(specRepo),
		},
	}

	c.logger.Info("Application context creation completed")
}

func (c *runnerContext) initMiddlewares() {
	c.addAuthMiddleware()

	c.logger.Info("Middlewares adding completed")
}

func (c *runnerContext) addAuthMiddleware() {
	authType := c.config.Auth.With

	switch authType {
	case config.FakeAuth:
		c.middlewares = append(c.middlewares, auth.FakeMiddleware)
	case config.FirebaseAuth:
		c.initFirebaseClient()
		c.middlewares = append(c.middlewares, auth.FirebaseMiddleware(c.firebaseAuth))
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

	c.logger.Info("Auth middleware creation completed", zap.String("authType", authType))
}

func (c *runnerContext) initServer() {
	v1Router := chi.NewRouter()
	v1Router.Use(c.middlewares...)

	c.server = server.New(c.config, http.NewHandler(
		c.logger,
		http.Route{
			Pattern: "/v1",
			Handler: v1.NewHandler(c.app, v1Router),
		},
		http.Route{
			Pattern: "/swagger",
			Handler: stdhttp.StripPrefix("/swagger/", stdhttp.FileServer(stdhttp.Dir("./swagger"))),
		},
	))

	c.logger.Info("Server initializing completed")
}
