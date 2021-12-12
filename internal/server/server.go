package server

import (
	"net/http"

	"github.com/harpyd/thestis/internal/config"
)

type Server struct {
	serv *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		serv: &http.Server{
			Addr:         ":" + cfg.HTTP.Port,
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.serv.ListenAndServe()
}
