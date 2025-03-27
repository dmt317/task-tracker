package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"task-tracker/internal/models"
)

func (s *HttpServer) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetAllTasks(w, r)
	case http.MethodPost:
		s.handleCreateTask(w, r)
	default:
		s.handleError(w, http.StatusMethodNotAllowed, models.ErrMethodNotAllowed)
	}
}

func (s *HttpServer) handleTaskById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetTask(w, r)
	case http.MethodDelete:
		s.handleDeleteTask(w, r)
	case http.MethodPatch:
		s.handleUpdateTask(w, r)
	default:
		s.handleError(w, http.StatusMethodNotAllowed, models.ErrMethodNotAllowed)
	}

}

func (s *HttpServer) handleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.storage.GetAll()
	if err != nil {
		s.processStorageError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		s.handleError(w, http.StatusInternalServerError, models.ErrEncodeJSON)
		return
	}
}

func (s *HttpServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		s.handleError(w, http.StatusBadRequest, models.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	task.Id = uuid.New().String()

	if err := s.storage.Add(&task); err != nil {
		s.processStorageError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		s.handleError(w, http.StatusInternalServerError, models.ErrEncodeJSON)
		return
	}
}

func (s *HttpServer) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskId := strings.TrimPrefix(r.URL.Path, "/task/")

	task, err := s.storage.Get(taskId)

	if err != nil {
		s.processStorageError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		s.handleError(w, http.StatusInternalServerError, models.ErrEncodeJSON)
		return
	}
}

func (s *HttpServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := strings.TrimPrefix(r.URL.Path, "/task/")

	if err := s.storage.Delete(taskId); err != nil {
		s.processStorageError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HttpServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	taskId := strings.TrimPrefix(r.URL.Path, "/task/")

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		s.handleError(w, http.StatusBadRequest, models.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	task.Id = taskId

	if err := s.storage.Update(task); err != nil {
		s.processStorageError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (s *HttpServer) handleError(w http.ResponseWriter, statusCode int, err error) {
	s.logger.Printf("HTTP error (%d): %s", statusCode, err)
	http.Error(w, err.Error(), statusCode)
}

func (s *HttpServer) processStorageError(w http.ResponseWriter, err error) {
	switch err {
	case models.ErrIdIsEmpty:
		s.handleError(w, http.StatusBadRequest, err)
	case models.ErrTaskNotFound:
		s.handleError(w, http.StatusNotFound, err)
	case models.ErrTaskExists:
		s.handleError(w, http.StatusConflict, err)
	default:
		s.handleError(w, http.StatusInternalServerError, err)
	}
}
