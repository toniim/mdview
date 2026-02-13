package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Host        string
	Port        int
	AllowedDirs []string
}

type Server struct {
	httpServer *http.Server
	config     Config
}

func New(cfg Config) *Server {
	s := &Server{config: cfg}
	mux := http.NewServeMux()
	s.registerRoutes(mux)
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}
