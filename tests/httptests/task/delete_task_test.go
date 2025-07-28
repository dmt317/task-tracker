package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"task-tracker/internal/models"
	"task-tracker/tests/testutils"
)

func TestDeleteTask(t *testing.T) {
	t.Run("happy path - delete task", func(t *testing.T) {
		t.Parallel()

		env := testutils.SetupIntegrationTest(t)

		task := models.CreateTaskRequest{
			Title:       "Task to delete",
			Description: "To be removed (",
			Status:      "todo",
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

		var created models.Task

		err = json.NewDecoder(resp.Body).Decode(&created)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		resp, err = env.Server.Handle(http.MethodDelete, "/tasks/"+created.ID, http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send delete request: %v", err)
		defer resp.Body.Close()

		require.Equalf(t, http.StatusNoContent, resp.StatusCode, "expected status %d, got %d", http.StatusNoContent, resp.StatusCode)

		resp, err = env.Server.Handle(http.MethodGet, "/tasks/"+created.ID, http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send get request: %v", err)
		defer resp.Body.Close()

		require.Equalf(t, http.StatusNotFound, resp.StatusCode, "expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	})

	t.Run("unhappy path - delete non-existent task", func(t *testing.T) {
		t.Parallel()

		env := testutils.SetupIntegrationTest(t)

		resp, err := env.Server.Handle(http.MethodDelete, "/tasks/nonexistent-id", http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send delete request: %v", err)
		defer resp.Body.Close()

		require.Equalf(t, http.StatusNotFound, resp.StatusCode, "expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	})
}
