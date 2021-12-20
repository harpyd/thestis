package runner

import (
	"context"
	"fmt"
	stdhttp "net/http"

	firebase "firebase.google.com/go"
	fireauth "firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	mongorepo "github.com/harpyd/thestis/internal/adapter/repository/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	"github.com/harpyd/thestis/internal/app/query"
	"github.com/harpyd/thestis/internal/config"
	"github.com/harpyd/thestis/internal/port/http"
	"github.com/harpyd/thestis/internal/port/http/auth"
	v1 "github.com/harpyd/thestis/internal/port/http/v1"
	"github.com/harpyd/thestis/internal/server"
	"github.com/harpyd/thestis/pkg/database/mongodb"
)

func Start(configsPath string) {
	newRunner().start(configsPath)
}

type runnerContext struct {
	logger       *zap.Logger
	config       *config.Config
	mongoDB      *mongo.Database
	firebaseAuth *fireauth.Client
	app          *app.Application
	middlewares  []func(stdhttp.Handler) stdhttp.Handler
	server       *server.Server
}

func newRunner() *runnerContext {
	return &runnerContext{}
}

func (c *runnerContext) start(configsPath string) {
	sync := c.initLogger()
	defer sync()

	c.initConfig(configsPath)
	c.initMongoDatabase()
	c.initApplication()
	c.initMiddlewares()
	c.initServer()

	c.serve()
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
}

func (c *runnerContext) initMongoDatabase() {
	client, err := mongodb.NewClient(c.config.Mongo.URI, c.config.Mongo.Username, c.config.Mongo.Password)
	if err != nil {
		c.logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	c.mongoDB = client.Database(c.config.Mongo.DatabaseName)
}

func (c *runnerContext) initFirebaseClient() {
	opt := option.WithCredentialsFile(c.config.Firebase.ServiceAccountFile)

	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		c.logger.Fatal("Failed to create Firebase app", zap.Error(err))
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		c.logger.Fatal("Failed to create Firebase Auth client", zap.Error(err))
	}

	c.firebaseAuth = authClient
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
}

func (c *runnerContext) initMiddlewares() {
	c.addAuthMiddleware()
}

func (c *runnerContext) addAuthMiddleware() {
	switch c.config.Auth.With {
	case config.FakeAuth:
		c.middlewares = append(c.middlewares, auth.FakeMiddleware)
	case config.FirebaseAuth:
		c.initFirebaseClient()
		c.middlewares = append(c.middlewares, auth.FirebaseMiddleware(c.firebaseAuth))
	default:
	}
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
}

func (c *runnerContext) serve() {
	c.logger.Info(
		"HTTP server started",
		zap.String("port", fmt.Sprintf(":%s", c.config.HTTP.Port)),
	)

	err := c.server.Start()

	c.logger.Fatal("HTTP server stopped", zap.Error(err))
}
