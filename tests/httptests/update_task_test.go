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
	t.Run("happy path - update task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		createReq := models.CreateTaskRequest{
			Title:       "Old Title",
			Description: "Old Description",
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

		updateReq := models.UpdateTaskRequest{
			Title:       "New Title",
			Description: "New Description",
			Status:      "in_progress",
		}
		updateBody, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest(http.MethodPatch, env.ServerURL+"/tasks/"+created.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, err = http.Get(env.ServerURL + "/tasks/" + created.ID)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var updated models.Task

		err = json.NewDecoder(resp.Body).Decode(&updated)
		require.NoError(t, err)

		require.Equal(t, "New Title", updated.Title)
		require.Equal(t, "New Description", updated.Description)
		require.Equal(t, "in_progress", updated.Status)
	})

	t.Run("unhappy path - update non-existent task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		updateReq := models.UpdateTaskRequest{
			Title:       "Test",
			Description: "Desc",
			Status:      "todo",
		}
		body, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest(http.MethodPatch, env.ServerURL+"/tasks/nonexistent-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
