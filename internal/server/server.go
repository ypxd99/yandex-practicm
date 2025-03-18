package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ypxd99/yandex-practicm/util"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	util.GetLogger().Infof("shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

func NewServer(handler http.Handler) *Server {
	cfg := util.GetConfig().Server
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.ServerAddress,
			ReadTimeout:  time.Duration(cfg.RTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WTimeout) * time.Second,
			Handler:      handler,
		},
	}
}
