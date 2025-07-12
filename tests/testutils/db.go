package testutils

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	connStr = "postgres://postgres:postgres@localhost:5432/tasktracker?sslmode=disable"
)

func createTestDatabase(t *testing.T, dbName string) {
	t.Helper()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)

	if err != nil {
		t.Fatalf("error connecting Postgres: %v", err)
	}

	defer conn.Close(ctx)

	query := fmt.Sprintf(`CREATE DATABASE "%v"`, dbName)
	_, err = conn.Exec(ctx, query)

	if err != nil {
		t.Fatalf("error creating database %s: %v", dbName, err)
	}
}

func dropTestDatabase(t *testing.T, dbName string) {
	t.Helper()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)

	if err != nil {
		t.Fatalf("error connecting Postgres: %v", err)
		return
	}

	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, fmt.Sprintf(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, dbName))

	if err != nil {
		t.Fatalf("error terminating connections to database %s: %v", dbName, err)
		return
	}

	query := fmt.Sprintf(`DROP DATABASE "%v"`, dbName)
	_, err = conn.Exec(ctx, query)

	if err != nil {
		t.Fatalf("error dropping database %s: %v", dbName, err)
		return
	}
}

func runMigrations(t *testing.T, dbName string) {
	t.Helper()

	conn := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
	cmd := exec.Command("migrate", "-path", "../../migrations", "-database", conn, "up")

	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("error running migrations for database %s: %v\n%s", dbName, err, output)
	}
}
