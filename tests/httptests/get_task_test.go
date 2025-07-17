package httptests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"task-tracker/internal/models"
	"task-tracker/tests/testutils"
)

func TestGetTask(t *testing.T) {
	t.Parallel()

	t.Run("happy path - get existing task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		task := models.CreateTaskRequest{
			Title:       "Integration Test Task",
			Description: "Some description",
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

		var createdTask models.Task

		err = json.NewDecoder(resp.Body).Decode(&createdTask)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		getResp, err := env.Server.Handle(http.MethodGet, "/tasks/"+createdTask.ID, http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send get request: %v", err)

		defer getResp.Body.Close()

		require.Equalf(t, http.StatusOK, getResp.StatusCode, "expected status %d, got %d", http.StatusOK, getResp.StatusCode)

		var fetched models.Task

		err = json.NewDecoder(getResp.Body).Decode(&fetched)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		require.Equal(t, createdTask.ID, fetched.ID)
		require.Equal(t, createdTask.Title, fetched.Title)
		require.Equal(t, createdTask.Description, fetched.Description)
		require.Equal(t, createdTask.Status, fetched.Status)

		env.CleanUpTest(t)
	})

	t.Run("unhappy path - task not found", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		nonExistentID := uuid.New().String()

		resp, err := env.Server.Handle(http.MethodGet, "/tasks/"+nonExistentID, http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send get request: %v", err)

		defer resp.Body.Close()

		require.Equalf(t, http.StatusNotFound, resp.StatusCode, "expected status %d, got %d", http.StatusNotFound, resp.StatusCode)

		_, err = io.ReadAll(resp.Body)
		require.NoErrorf(t, err, "failed to read response: %v", err)

		env.CleanUpTest(t)
	})
}
