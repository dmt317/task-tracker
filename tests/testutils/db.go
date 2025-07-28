package testutils

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

const (
	connStr = "postgres://postgres:postgres@localhost:5432/tasktracker?sslmode=disable"
)

func createTestDatabase(t *testing.T, dbName string) {
	t.Helper()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)

	require.NoErrorf(t, err, "error connecting Postgres: %v", err)

	defer conn.Close(ctx)

	query := fmt.Sprintf(`CREATE DATABASE "%v"`, dbName)
	_, err = conn.Exec(ctx, query)

	require.NoErrorf(t, err, "error creating database %s: %v", dbName, err)
}

func dropTestDatabase(t *testing.T, dbName string) {
	t.Helper()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)

	require.NoErrorf(t, err, "error connecting Postgres: %v", err)

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, fmt.Sprintf(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, dbName))

	require.NoErrorf(t, err, "error terminating connections to db %s: %v", dbName, err)

	query := fmt.Sprintf(`DROP DATABASE "%v"`, dbName)
	_, err = conn.Exec(ctx, query)

	require.NoErrorf(t, err, "error dropping database %s: %v", dbName, err)
}

func runMigrations(t *testing.T, dbName string) {
	t.Helper()

	conn := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
	cmd := exec.Command("migrate", "-path", "../../../migrations", "-database", conn, "up")

	output, err := cmd.CombinedOutput()

	require.NoErrorf(t, err, "error running migrations for db %s: %v\n%s", dbName, err, output)
}
