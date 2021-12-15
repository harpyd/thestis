package runner

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/harpyd/thestis/internal/adapter/parser/yaml"
	mongorepo "github.com/harpyd/thestis/internal/adapter/repository/mongodb"
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	"github.com/harpyd/thestis/internal/config"
	"github.com/harpyd/thestis/internal/port/http"
	"github.com/harpyd/thestis/internal/server"
	"github.com/harpyd/thestis/pkg/database/mongodb"
)

func Start(configsPath string) {
	logger, sync := newLogger()
	defer sync()

	cfg := newConfig(configsPath, logger)
	db := newMongoDatabase(cfg, logger)
	application := newApplication(db)

	startServer(cfg, application, logger)
}

func newLogger() (*zap.Logger, func()) {
	logger, _ := zap.NewProduction()
	sync := func() {
		_ = logger.Sync()
	}

	return logger, sync
}

func newConfig(configsPath string, logger *zap.Logger) *config.Config {
	cfg, err := config.FromPath(configsPath)
	if err != nil {
		logger.Fatal("Failed to parse config", zap.Error(err))
	}

	return cfg
}

func newMongoDatabase(cfg *config.Config, logger *zap.Logger) *mongo.Database {
	client, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.Username, cfg.Mongo.Password)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	return client.Database(cfg.Mongo.DatabaseName)
}

func newApplication(db *mongo.Database) app.Application {
	tcRepo := mongorepo.NewTestCampaignsRepository(db)
	specRepo := mongorepo.NewSpecificationsRepository(db)
	parserService := yaml.NewSpecificationParserService()

	return app.Application{
		Commands: app.Commands{
			CreateTestCampaign: command.NewCreateTestCampaignHandler(tcRepo),
			LoadSpecification:  command.NewLoadSpecificationHandler(specRepo, tcRepo, parserService),
		},
	}
}

func startServer(cfg *config.Config, application app.Application, logger *zap.Logger) {
	logger.Info("HTTP server started", zap.String("port", fmt.Sprintf(":%s", cfg.HTTP.Port)))

	serv := server.New(cfg, http.NewHandler(application, logger))
	err := serv.Start()

	logger.Fatal("HTTP server stopped", zap.Error(err))
}
