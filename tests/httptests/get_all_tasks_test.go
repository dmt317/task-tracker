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

func TestGetAllTasks(t *testing.T) {
	t.Parallel()

	t.Run("happy path - get all tasks", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		tasks := []models.CreateTaskRequest{
			{
				Title:       "First Task",
				Description: "Description 1",
				Status:      "todo",
			},
			{
				Title:       "Second Task",
				Description: "Description 2",
				Status:      "todo",
			},
		}

		headers := map[string]string{
			"Content-Type": "application/json",
		}

		for _, task := range tasks {
			body, err := json.Marshal(task)
			require.NoErrorf(t, err, "failed to marshal task request: %v", err)

			resp, err := env.Server.Handle(http.MethodPost, "/tasks", bytes.NewReader(body), headers)
			require.NoErrorf(t, err, "failed to send post request: %v", err)

			require.Equalf(t, http.StatusCreated, resp.StatusCode, "expected status %d, got %d", http.StatusCreated, resp.StatusCode)

			resp.Body.Close()
		}

		resp, err := env.Server.Handle(http.MethodGet, "/tasks", http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send get request: %v", err)
		defer resp.Body.Close()

		require.Equalf(t, http.StatusOK, resp.StatusCode, "expected status %d, got %d", http.StatusOK, resp.StatusCode)

		var gottenTasks []models.Task

		err = json.NewDecoder(resp.Body).Decode(&gottenTasks)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		require.Len(t, gottenTasks, 2)

		titles := []string{gottenTasks[0].Title, gottenTasks[1].Title}
		require.Contains(t, titles, tasks[0].Title)
		require.Contains(t, titles, tasks[1].Title)

		env.CleanUpTest(t)
	})

	t.Run("unhappy path - no tasks in db", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		resp, err := env.Server.Handle(http.MethodGet, "/tasks", http.NoBody, nil)
		require.NoErrorf(t, err, "failed to send get request: %v", err)
		defer resp.Body.Close()

		require.Equalf(t, http.StatusOK, resp.StatusCode, "expected status %d, got %d", http.StatusOK, resp.StatusCode)

		var gottenTasks []models.Task

		err = json.NewDecoder(resp.Body).Decode(&gottenTasks)
		require.NoErrorf(t, err, "failed to decode response: %v", err)

		require.Len(t, gottenTasks, 0)

		env.CleanUpTest(t)
	})
}
