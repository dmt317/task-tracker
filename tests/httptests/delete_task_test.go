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

func TestDeleteTask(t *testing.T) {
	t.Run("happy path - delete task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)
		createReq := models.CreateTaskRequest{
			Title:       "Task to delete",
			Description: "To be removed",
			Status:      "todo",
		}
		body, _ := json.Marshal(createReq)
		resp, err := http.Post(env.ServerURL+"/tasks", "application/json", bytes.NewReader(body))

		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created models.Task

		err = json.NewDecoder(resp.Body).Decode(&created)

		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodDelete, env.ServerURL+"/tasks/"+created.ID, http.NoBody)
		client := &http.Client{}
		resp, err = client.Do(req)

		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, err = http.Get(env.ServerURL + "/tasks/" + created.ID)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("unhappy path - delete non-existent task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		req, _ := http.NewRequest(http.MethodDelete, env.ServerURL+"/tasks/nonexistent-id", http.NoBody)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
