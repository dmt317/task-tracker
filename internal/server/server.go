package server

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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
	mux         *http.ServeMux
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

func (s *HTTPServer) Handle(method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req := httptest.NewRequest(method, path, body)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rr := httptest.NewRecorder()

	s.mux.ServeHTTP(rr, req)

	res := &http.Response{
		StatusCode: rr.Code,
		Body:       io.NopCloser(bytes.NewReader(rr.Body.Bytes())),
		Header:     rr.Header(),
	}

	return res, nil
}

func (s *HTTPServer) ConfigureServer(ctx context.Context) error {
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

	s.mux = http.NewServeMux()

	s.setupRoutes(s.mux)

	s.server = &http.Server{
		Addr:              ":" + s.config.ServerPort,
		Handler:           s.mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return nil
}

func (s *HTTPServer) Start(cancel context.CancelFunc) error {
	return s.startHTTPServer(cancel)
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) startHTTPServer(cancel context.CancelFunc) error {
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
