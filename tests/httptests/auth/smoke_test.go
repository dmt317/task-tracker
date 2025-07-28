package auth

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"task-tracker/internal/config"
)

func TestPostgresConnection(t *testing.T) {
	config := config.LoadConfig()

	conn, err := pgx.Connect(context.Background(), config.AuthDBConn)

	require.NoErrorf(t, err, "error connecting to auth DB: %v", err)

	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "SELECT 1;")

	require.NoErrorf(t, err, "error executing query to auth DB: %v", err)
}
