package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"task-tracker/internal/config"
	"task-tracker/internal/server"
)

type TestEnv struct {
	DBName string
	Server *server.HTTPServer
}

func SetupIntegrationTest(t *testing.T) *TestEnv {
	t.Helper()

	dbName := uuid.New().String()
	createTestDatabase(t, dbName)
	runMigrations(t, dbName)

	cfg := config.Config{
		DBConn:     fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName),
		ServerPort: "5050",
		InMemory:   "False",
	}

	server := server.NewHTTPServer(cfg)

	err := server.ConfigureServer(context.Background())
	if err != nil {
		panic(err)
	}

	env := &TestEnv{
		DBName: dbName,
		Server: server,
	}

	env.CleanUpTest(t)

	return env
}

func (env *TestEnv) CleanUpTest(t *testing.T) {
	t.Cleanup(func() {
		dropTestDatabase(t, env.DBName)
	})
}
