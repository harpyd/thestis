package runner

import (
	"net/http"

	"go.uber.org/zap"

	httpserver "github.com/harpyd/thestis/internal/port/http"
)

func Start() {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("HTTP server started", zap.String("port", ":8080"))
	err := http.ListenAndServe(":8080", httpserver.NewHandler(logger))
	logger.Fatal("HTTP server stopped", zap.Error(err))
}
