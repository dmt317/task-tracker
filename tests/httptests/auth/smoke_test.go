package auth

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestPostgresConnection(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connStr)

	require.NoErrorf(t, err, "error connecting to auth DB: %v", err)

	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "SELECT 1;")

	require.NoErrorf(t, err, "error executing query to auth DB: %v", err)
}
