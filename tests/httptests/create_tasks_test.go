package httptests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"task-tracker/internal/models"
	"task-tracker/tests/testutils"
)

func TestCreateTask(t *testing.T) {
	t.Run("happy path - create valid task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		taskReq := models.CreateTaskRequest{
			Title:       "Test Task",
			Description: "Some task to test",
			Status:      "Todo",
		}

		body, err := json.Marshal(taskReq)

		if err != nil {
			t.Fatalf("failed to marshal task request: %v", err)
		}

		resp, err := http.Post(env.ServerURL+"/tasks", "application/json", bytes.NewReader(body))

		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		var createdTask models.Task

		if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if createdTask.Title != taskReq.Title {
			t.Errorf("expected title %s, got %s", taskReq.Title, createdTask.Title)
		}
	})

	t.Run("unhappy path - empty title", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		taskReq := models.CreateTaskRequest{
			Title:       "",
			Description: "Missing title field",
			Status:      "Todo",
		}

		body, _ := json.Marshal(taskReq)

		resp, err := http.Post(env.ServerURL+"/tasks", "application/json", bytes.NewReader(body))

		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}
