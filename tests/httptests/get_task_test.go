package httptests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"task-tracker/internal/models"
	"task-tracker/tests/testutils"
)

func TestGetTask(t *testing.T) {
	t.Run("happy path - get existing task", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		createReq := models.CreateTaskRequest{
			Title:       "Integration Test Task",
			Description: "Some description",
			Status:      "todo",
		}

		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		resp, err := http.Post(env.ServerURL+"/tasks", "application/json", bytes.NewReader(reqBody))
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created models.Task

		err = json.NewDecoder(resp.Body).Decode(&created)
		require.NoError(t, err)

		getResp, err := http.Get(fmt.Sprintf("%s/tasks/%s", env.ServerURL, created.ID))
		require.NoError(t, err)

		defer getResp.Body.Close()

		require.Equal(t, http.StatusOK, getResp.StatusCode)

		var fetched models.Task

		err = json.NewDecoder(getResp.Body).Decode(&fetched)
		require.NoError(t, err)

		require.Equal(t, created.ID, fetched.ID)
		require.Equal(t, created.Title, fetched.Title)
		require.Equal(t, created.Description, fetched.Description)
		require.Equal(t, created.Status, fetched.Status)
	})

	t.Run("unhappy path - task not found", func(t *testing.T) {
		env := testutils.SetupIntegrationTest(t)

		nonExistentID := uuid.New().String()

		resp, err := http.Get(fmt.Sprintf("%s/tasks/%s", env.ServerURL, nonExistentID))
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		t.Logf("server responded with: %s", string(body))
	})
}
