package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

type ServerConfig struct {
	Host string
	Port string
}

type Server struct {
	server *http.Server
	logger *slog.Logger
}

func NewServer(config *ServerConfig, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		server: &http.Server{
			Handler: handler,
			Addr:    net.JoinHostPort(config.Host, config.Port),
		},
		logger: logger,
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic("http server error")
		}
	}()
	s.logger.Info(fmt.Sprintf("listening on %s", s.server.Addr))
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.server.Shutdown(ctx)
}
