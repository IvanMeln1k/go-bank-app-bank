package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ServerConfig struct {
	Host string
	Port string
}

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg ServerConfig, handler http.Handler) *Server {
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}

	return &Server{
		httpServer: server,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Close()
}
