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
	"task-tracker/internal/repository"
	"task-tracker/internal/service"
)

type HTTPServer struct {
	config      config.Config
	logger      *log.Logger
	taskService service.TaskService
	server      *http.Server
}

func NewHTTPServer(config config.Config) *HTTPServer {
	return &HTTPServer{
		config: config,
	}
}

func (s *HTTPServer) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", s.handleTasks)
	mux.HandleFunc("/tasks/{id}", s.handleTaskByID)
	mux.HandleFunc("/swagger", s.handleSwagger)

	mux.Handle("/swagger/static/", http.StripPrefix("/swagger/static/", http.FileServer(http.Dir("docs/static"))))
	mux.Handle("/swagger/swagger.yaml", http.StripPrefix("/swagger/", http.FileServer(http.Dir("docs"))))
}

func (s *HTTPServer) configureServer(ctx context.Context) error {
	pool, err := repository.CreateDBPool(ctx, s.config.DBConn)
	if err != nil {
		return err
	}

	var repo repository.TaskRepository

	if s.config.InMemory == "True" {
		repo = repository.NewMemoryTaskRepository()
	} else {
		repo = repository.NewPostgresTaskRepository(pool)
	}

	s.taskService = service.NewDefaultTaskService(repo)
	s.logger = log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags)

	return nil
}

func (s *HTTPServer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	return s.startHTTPServer(ctx, cancel)
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) startHTTPServer(ctx context.Context, cancel context.CancelFunc) error {
	err := s.configureServer(ctx)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	s.setupRoutes(mux)

	s.server = &http.Server{
		Addr:              ":" + s.config.ServerPort,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		s.logger.Println("Starting HTTP server on port", s.config.ServerPort)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigs
	s.logger.Println("Shutting down server...")

	cancel()

	shutdownCtx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		s.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	s.logger.Println("Server gracefully shut down")

	return s.server.Shutdown(shutdownCtx)
}
