package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"task-tracker/internal/models"
)

func (s *HTTPServer) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetAllTasks(w, r)
	case http.MethodPost:
		s.handleCreateTask(w, r)
	default:
		s.handleError(w, r.RemoteAddr, models.ErrMethodNotAllowed)
	}
}

func (s *HTTPServer) handleSwagger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.handleError(w, r.RemoteAddr, models.ErrMethodNotAllowed)
		return
	}

	if _, err := os.Stat("docs/static/index.html"); os.IsNotExist(err) {
		s.handleError(w, r.RemoteAddr, models.ErrSwaggerUINotFound)
		return
	}

	http.ServeFile(w, r, "docs/static/index.html")
}

func (s *HTTPServer) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetTask(w, r)
	case http.MethodDelete:
		s.handleDeleteTask(w, r)
	case http.MethodPatch:
		s.handleUpdateTask(w, r)
	default:
		s.handleError(w, r.RemoteAddr, models.ErrMethodNotAllowed)
	}
}

func (s *HTTPServer) handleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.taskService.GetAll(r.Context())
	if err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}
}

func (s *HTTPServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var request models.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.handleError(w, r.RemoteAddr, models.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	if err := request.Validate(); err != nil {
		s.handleError(w, r.RemoteAddr, fmt.Errorf("request validation: %w", err))
		return
	}

	task := request.ConvertToTask()

	if err := s.taskService.Add(r.Context(), task); err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}
}

func (s *HTTPServer) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	task, err := s.taskService.Get(r.Context(), taskID)

	if err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}
}

func (s *HTTPServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	if err := s.taskService.Delete(r.Context(), taskID); err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	var request models.UpdateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.handleError(w, r.RemoteAddr, models.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	if err := request.Validate(); err != nil {
		s.handleError(w, r.RemoteAddr, fmt.Errorf("request validation: %w", err))
		return
	}

	task := request.ConvertToTask(taskID)

	if err := s.taskService.Update(r.Context(), task); err != nil {
		s.handleError(w, r.RemoteAddr, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) handleError(w http.ResponseWriter, ip string, err error) {
	var modelError models.Error

	if !errors.As(err, &modelError) {
		s.processError(w, http.StatusInternalServerError, ip, err)
	}

	s.processError(w, modelError.StatusCode, ip, modelError.Err)
}

func (s *HTTPServer) processError(w http.ResponseWriter, statusCode int, ip string, err error) {
	s.logger.Printf("HTTP error (%d) from %s: %s", statusCode, ip, err)
	http.Error(w, err.Error(), statusCode)
}
