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

func TestUpdateTask(t *testing.T) {
	t.Parallel()

	t.Run("happy path - update task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		task := models.CreateTaskRequest{
			Title:       "Old Title",
			Description: "Old Description",
			Status:      "todo",
		}
		body, err := json.Marshal(task)
		require.NoErrorf(t, err, "failed to marshal task request: %v", err)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		resp, err := env.Server.Handle(http.MethodPost, "/tasks", bytes.NewReader(body), headers)
		require.NoErrorf(t, err, "failed to send request: %v", err)

		defer resp.Body.Close()

		require.Equalf(t, http.StatusCreated, resp.StatusCode, "expected status %d, got %d", http.StatusCreated, resp.StatusCode)

		var created models.Task

		err = json.NewDecoder(resp.Body).Decode(&created)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		updateTask := models.UpdateTaskRequest{
			Title:       "New Title",
			Description: "New Description",
			Status:      "in_progress",
		}
		updateBody, err := json.Marshal(updateTask)
		require.NoErrorf(t, err, "failed to marshal task request: %v", err)

		updateResp, err := env.Server.Handle(http.MethodPatch, "/tasks/"+created.ID, bytes.NewReader(updateBody), headers)
		require.NoErrorf(t, err, "failed to send request: %v", err)

		defer updateResp.Body.Close()

		require.Equal(t, http.StatusNoContent, updateResp.StatusCode, "expected status %d, got %d", http.StatusNoContent, resp.StatusCode)

		resp, err = env.Server.Handle(http.MethodGet, "/tasks/"+created.ID, http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send request: %v", err)

		defer resp.Body.Close()

		require.Equalf(t, http.StatusOK, resp.StatusCode, "expected status %d, got %d", http.StatusOK, resp.StatusCode)

		var updatedTask models.Task

		err = json.NewDecoder(resp.Body).Decode(&updatedTask)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		require.Equal(t, updateTask.Title, updatedTask.Title)
		require.Equal(t, updateTask.Description, updatedTask.Description)
		require.Equal(t, updateTask.Status, updatedTask.Status)

		env.CleanUpTest(t)
	})

	t.Run("unhappy path - update non-existent task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		updateTask := models.UpdateTaskRequest{
			Title:       "Test",
			Description: "Desc",
			Status:      "todo",
		}
		body, err := json.Marshal(updateTask)
		require.NoErrorf(t, err, "failed to marshal task request: %v", err)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		resp, err := env.Server.Handle(http.MethodPatch, "/tasks/"+"non-existent", bytes.NewReader(body), headers)
		require.NoErrorf(t, err, "failed to send patch request: %v", err)

		defer resp.Body.Close()

		require.Equalf(t, http.StatusNotFound, resp.StatusCode, "expected status %d, got %d", http.StatusNotFound, resp.StatusCode)

		env.CleanUpTest(t)
	})
}
