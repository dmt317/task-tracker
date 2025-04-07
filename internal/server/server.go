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

type HTTPServer struct {
	config  config.Config
	logger  *log.Logger
	storage *storage.Storage
}

func NewHTTPServer(config config.Config) *HTTPServer {
	return &HTTPServer{
		config:  config,
		logger:  log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags),
		storage: storage.NewStorage(),
	}
}

func (s *HTTPServer) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", s.handleTasks)
	mux.HandleFunc("/tasks/{id}", s.handleTaskByID)
}

func (s *HTTPServer) StartHTTPServer() {
	mux := http.NewServeMux()

	s.setupRoutes(mux)

	server := &http.Server{
		Addr:              ":" + s.config.ServerPort,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
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
