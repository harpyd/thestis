package server

import (
	"context"
	"net/http"

	"github.com/harpyd/thestis/internal/config"
)

type Server struct {
	serv *http.Server
}

func New(cfg config.HTTP, handler http.Handler) *Server {
	return &Server{
		serv: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.serv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.serv.Shutdown(ctx)
}
