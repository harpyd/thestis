package server

import (
	"context"
	"net/http"

	"github.com/harpyd/thestis/internal/config"
)

type Server struct {
	serv *http.Server
	cfg  config.HTTP
}

func New(cfg config.HTTP, handler http.Handler) *Server {
	return &Server{
		serv: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	return s.serv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	return s.serv.Shutdown(ctx)
}
