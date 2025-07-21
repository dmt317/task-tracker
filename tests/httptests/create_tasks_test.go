package httptests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"task-tracker/internal/models"
	"task-tracker/tests/testutils"
)

func TestCreateTask(t *testing.T) {
	t.Run("happy path - create valid task", func(t *testing.T) {
		t.Parallel()

		env := testutils.SetupIntegrationTest(t)

		task := models.CreateTaskRequest{
			Title:       "Test Task",
			Description: "Some task to test",
			Status:      "Todo",
		}

		body, err := json.Marshal(task)
		require.NoErrorf(t, err, "failed to marshal task request: %v", err)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		resp, err := env.Server.Handle(http.MethodPost, "/tasks", bytes.NewReader(body), headers)
		require.NoErrorf(t, err, "failed to send post request: %v", err)

		defer resp.Body.Close()

		require.Equalf(t, http.StatusCreated, resp.StatusCode, "expected status %d, got %d", http.StatusCreated, resp.StatusCode)

		var createdTask models.Task

		err = json.NewDecoder(resp.Body).Decode(&createdTask)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		require.Equalf(t, task.Title, createdTask.Title, "expected title %s, got %s", task.Title, createdTask.Title)
	})

	t.Run("unhappy path - empty title", func(t *testing.T) {
		t.Parallel()

		env := testutils.SetupIntegrationTest(t)

		task := models.CreateTaskRequest{
			Title:       "",
			Description: "Missing title field",
			Status:      "Todo",
		}

		body, err := json.Marshal(task)
		require.NoErrorf(t, err, "failed to marshal task request: %v", err)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		resp, err := env.Server.Handle(http.MethodPost, "/tasks", bytes.NewReader(body), headers)
		require.NoErrorf(t, err, "failed to send post request: %v", err)

		defer resp.Body.Close()

		require.Equalf(t, http.StatusBadRequest, resp.StatusCode, "expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	})
}
