package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-tracker/internal/config"
	"task-tracker/internal/storage"
)

type HttpServer struct {
	config  config.Config
	logger  *log.Logger
	storage *storage.Storage
}

func NewHttpServer(config config.Config) *HttpServer {
	return &HttpServer{
		config:  config,
		logger:  log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags),
		storage: storage.NewStorage(),
	}
}

func (s *HttpServer) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/task", s.handleTasks)
	mux.HandleFunc("/task/", s.handleTaskById)
}

func (s *HttpServer) StartHttpServer() {
	mux := http.NewServeMux()

	s.setupRoutes(mux)

	server := &http.Server{
		Addr:    ":" + s.config.ServerPort,
		Handler: mux,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		s.logger.Println("Starting HTTP server on port", s.config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigs
	s.logger.Println("Shutting down server...")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := server.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	s.logger.Println("Server gracefully shut down")
}
