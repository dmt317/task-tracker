package server

import (
	"encoding/json"
	"net/http"

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
		s.handleError(w, http.StatusMethodNotAllowed, r.RemoteAddr, models.ErrMethodNotAllowed)
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
		s.handleError(w, http.StatusMethodNotAllowed, r.RemoteAddr, models.ErrMethodNotAllowed)
	}
}

func (s *HttpServer) handleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.storage.GetAll()
	if err != nil {
		s.processStorageError(w, r.RemoteAddr, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		s.handleError(w, http.StatusInternalServerError, r.RemoteAddr, models.ErrEncodeJSON)
		return
	}
}

func (s *HttpServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		s.handleError(w, http.StatusBadRequest, r.RemoteAddr, models.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	task.ID = uuid.New().String()

	if err := s.storage.Add(&task); err != nil {
		s.processStorageError(w, r.RemoteAddr, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		s.handleError(w, http.StatusInternalServerError, r.RemoteAddr, models.ErrEncodeJSON)
		return
	}
}

func (s *HttpServer) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	task, err := s.storage.Get(taskID)

	if err != nil {
		s.processStorageError(w, r.RemoteAddr, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		s.handleError(w, http.StatusInternalServerError, r.RemoteAddr, models.ErrEncodeJSON)
		return
	}
}

func (s *HttpServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	if err := s.storage.Delete(taskID); err != nil {
		s.processStorageError(w, r.RemoteAddr, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HttpServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		s.handleError(w, http.StatusBadRequest, r.RemoteAddr, models.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	task.ID = taskID

	if err := s.storage.Update(task); err != nil {
		s.processStorageError(w, r.RemoteAddr, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HttpServer) handleError(w http.ResponseWriter, statusCode int, ip string, err error) {
	s.logger.Printf("HTTP error (%d) from %s: %s", statusCode, ip, err)
	http.Error(w, err.Error(), statusCode)
}

func (s *HttpServer) processStorageError(w http.ResponseWriter, ip string, err error) {
	switch err {
	case models.ErrIdIsEmpty:
		s.handleError(w, http.StatusBadRequest, ip, err)
	case models.ErrTaskNotFound:
		s.handleError(w, http.StatusNotFound, ip, err)
	case models.ErrTaskExists:
		s.handleError(w, http.StatusConflict, ip, err)
	default:
		s.handleError(w, http.StatusInternalServerError, ip, err)
	}
}
