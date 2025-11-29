package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"zentinel/pkg/config"
	"zentinel/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler *gin.Engine) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           fmt.Sprintf(":%s", cfg.ServerPort),
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *Server) Run() error {
	logger.Info("Server starting", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Server shutting down")
	return s.httpServer.Shutdown(ctx)
}
