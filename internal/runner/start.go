package runner

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/harpyd/thestis/internal/config"
	"github.com/harpyd/thestis/internal/port/http"
	"github.com/harpyd/thestis/internal/server"
)

func Start(configsPath string) {
	logger, sync := newLogger()
	defer sync()

	cfg := newConfig(configsPath, logger)

	startServer(cfg, logger)
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

func startServer(cfg *config.Config, logger *zap.Logger) {
	logger.Info("HTTP server started", zap.String("port", fmt.Sprintf(":%s", cfg.HTTP.Port)))

	serv := server.New(cfg, http.NewHandler(logger))
	err := serv.Start()

	logger.Fatal("HTTP server stopped", zap.Error(err))
}
