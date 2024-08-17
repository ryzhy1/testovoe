package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	log          *slog.Logger
	port         string
	handler      *gin.Engine
	httpServer   *http.Server
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewServer(log *slog.Logger, port string, handler *gin.Engine) *Server {
	return &Server{
		log:          log,
		port:         port,
		handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	const op = "HTTPServer.Run"

	log := s.log.With(
		slog.String("op", op),
		slog.String("port", s.port),
	)

	log.Info("HTTP server started", "port", s.port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.handler); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Stop() error {
	const op = "HTTPServer.Stop"

	s.log.With(slog.String("op", op)).
		Info("HTTP server stopped")

	return s.httpServer.Shutdown(context.Background())
}
