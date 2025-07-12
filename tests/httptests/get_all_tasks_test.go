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
	t.Run("happy path - get all tasks", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		tasksToCreate := []models.CreateTaskRequest{
			{
				Title:       "First Task",
				Description: "Description 1",
				Status:      "todo",
			},
			{
				Title:       "Second Task",
				Description: "Description 2",
				Status:      "in_progress",
			},
		}

		for _, task := range tasksToCreate {
			body, err := json.Marshal(task)
			require.NoError(t, err)

			resp, err := http.Post(env.ServerURL+"/tasks", "application/json", bytes.NewReader(body))
			require.NoError(t, err)

			resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)
		}

		resp, err := http.Get(env.ServerURL + "/tasks")
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var tasks []models.Task

		err = json.NewDecoder(resp.Body).Decode(&tasks)
		require.NoError(t, err)

		require.Len(t, tasks, 2)

		titles := []string{tasks[0].Title, tasks[1].Title}
		require.Contains(t, titles, "First Task")
		require.Contains(t, titles, "Second Task")
	})

	t.Run("unhappy path - no tasks in db", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		resp, err := http.Get(env.ServerURL + "/tasks")
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var tasks []models.Task

		err = json.NewDecoder(resp.Body).Decode(&tasks)
		require.NoError(t, err)

		require.Len(t, tasks, 0)
	})
}
