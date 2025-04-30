package server

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"task-tracker/internal/config"
	"task-tracker/internal/service"
)

func TestHandler_CreateTask(t *testing.T) {
	tests := map[string]struct {
		requestBody    string
		mockSetup      *service.TaskServiceMock
		expectedStatus int
	}{
		"success": {
			requestBody: `{"title":"title", "description":"description", "status":"todo"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusCreated,
		},

		"bad request on invalid json": {
			requestBody: `{"title":"title"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusBadRequest,
		},

		"conflict on task already exists": {
			requestBody: `{"title":"existing", "description":"description", "status":"todo"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusConflict,
		},

		"internal server error on service failure": {
			requestBody: `{"title":"title", "description":"description", "status":"todo"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: true,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := &HTTPServer{
				config:      *config.LoadConfig(),
				logger:      log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags),
				taskService: test.mockSetup,
			}

			body := bytes.NewBufferString(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/tasks", body)
			w := httptest.NewRecorder()

			server.handleCreateTask(w, req)

			if test.expectedStatus != w.Code {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, w.Code, test.expectedStatus)
			}
		})
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	tests := map[string]struct {
		taskID         string
		mockSetup      *service.TaskServiceMock
		expectedStatus int
	}{
		"success": {
			taskID: "task1",
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusNoContent,
		},

		"not found": {
			taskID: service.NotFound,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusNotFound,
		},

		"internal server error on service failure": {
			taskID: "task1",
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: true,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := &HTTPServer{
				config:      *config.LoadConfig(),
				logger:      log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags),
				taskService: test.mockSetup,
			}

			req := httptest.NewRequest(http.MethodDelete, "/tasks/{id}", http.NoBody)
			req.SetPathValue("id", test.taskID)

			w := httptest.NewRecorder()

			server.handleDeleteTask(w, req)

			if test.expectedStatus != w.Code {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, w.Code, test.expectedStatus)
			}
		})
	}
}

func TestHandler_GetTask(t *testing.T) {
	tests := map[string]struct {
		taskID         string
		mockSetup      *service.TaskServiceMock
		expectedStatus int
	}{
		"success": {
			taskID: "task1",
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusOK,
		},

		"not found": {
			taskID: service.NotFound,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusNotFound,
		},

		"internal server error on service failure": {
			taskID: "task1",
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: true,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := &HTTPServer{
				config:      *config.LoadConfig(),
				logger:      log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags),
				taskService: test.mockSetup,
			}

			req := httptest.NewRequest(http.MethodGet, "/tasks/{id}", http.NoBody)
			req.SetPathValue("id", test.taskID)

			w := httptest.NewRecorder()

			server.handleGetTask(w, req)

			if test.expectedStatus != w.Code {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, w.Code, test.expectedStatus)
			}
		})
	}
}

func TestHandler_GetAllTasks(t *testing.T) {
	tests := map[string]struct {
		mockSetup      *service.TaskServiceMock
		expectedStatus int
	}{
		"success": {
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusOK,
		},

		"internal server error on service failure": {
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: true,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := &HTTPServer{
				config:      *config.LoadConfig(),
				logger:      log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags),
				taskService: test.mockSetup,
			}

			req := httptest.NewRequest(http.MethodGet, "/tasks", http.NoBody)
			w := httptest.NewRecorder()

			server.handleGetTask(w, req)

			if test.expectedStatus != w.Code {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, w.Code, test.expectedStatus)
			}
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) {
	tests := map[string]struct {
		taskID         string
		requestBody    string
		mockSetup      *service.TaskServiceMock
		expectedStatus int
	}{
		"success": {
			taskID:      "task1",
			requestBody: `{"title":"title", "description":"description", "status":"todo"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusNoContent,
		},

		"bad request on invalid json": {
			taskID:      "task1",
			requestBody: `{"title":"title"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusBadRequest,
		},

		"not found": {
			taskID:      service.NotFound,
			requestBody: `{"title":"title", "description":"description", "status":"todo"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: false,
			},
			expectedStatus: http.StatusNotFound,
		},

		"interval server error": {
			taskID:      "task1",
			requestBody: `{"title":"title", "description":"description", "status":"todo"}`,
			mockSetup: &service.TaskServiceMock{
				ForceInternalError: true,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := &HTTPServer{
				config:      *config.LoadConfig(),
				logger:      log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags),
				taskService: test.mockSetup,
			}

			body := bytes.NewBufferString(test.requestBody)
			req := httptest.NewRequest(http.MethodGet, "/tasks/{id}", body)
			req.SetPathValue("id", test.taskID)

			w := httptest.NewRecorder()

			server.handleUpdateTask(w, req)

			if test.expectedStatus != w.Code {
				t.Fatalf("test-case: (%q); returned %v; expected %v", name, w.Code, test.expectedStatus)
			}
		})
	}
}
