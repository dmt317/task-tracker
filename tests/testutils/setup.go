package testutils

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"task-tracker/internal/config"
	"task-tracker/internal/server"
)

type TestEnv struct {
	DBName    string
	ServerURL string
}

func SetupIntegrationTest(t *testing.T) *TestEnv {
	t.Helper()

	dbName := uuid.New().String()
	createTestDatabase(t, dbName)
	runMigrations(t, dbName)

	cfg := config.Config{
		DBConn:     fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName),
		ServerPort: "8080",
		InMemory:   "False",
	}

	server := server.NewHTTPServer(cfg)

	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	waitForServer(t, "http://localhost:"+cfg.ServerPort)

	t.Cleanup(func() {
		ctx, shutdown := context.WithTimeout(context.Background(), 2*time.Second)
		defer shutdown()

		err := server.Stop(ctx)
		if err != nil {
			t.Fatalf("failed to shutdown test server: %v", err)
		}

		dropTestDatabase(t, dbName)
	})

	return &TestEnv{
		DBName:    dbName,
		ServerURL: "http://localhost:" + cfg.ServerPort,
	}
}

func waitForServer(t *testing.T, url string) {
	t.Helper()

	const retries = 20

	const delay = 250 * time.Millisecond

	for i := 0; i < retries; i++ {
		resp, err := http.Get(url + "/swagger")
		if err == nil && resp.StatusCode < 500 {
			resp.Body.Close()
			return
		}

		time.Sleep(delay)
	}

	t.Fatalf("server did not start on %s after %d attempts", url, retries)
}
