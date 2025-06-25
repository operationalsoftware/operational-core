package migrate

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Check if the database exists
func checkDatabaseExists(ctx context.Context, conn *pgx.Conn, dbName string) (bool, error) {
	var exists bool
	query := `
SELECT EXISTS (
	SELECT 1 FROM pg_database WHERE datname = $1
)`
	err := conn.QueryRow(ctx, query, dbName).Scan(&exists)
	return exists, err
}

// Check if the database is initialised (by checking if the app_users table exists)
func checkInitialised(ctx context.Context, pgPool *pgxpool.Pool) (bool, error) {
	var exists bool
	query := `
SELECT EXISTS (
	SELECT 1 FROM information_schema.tables WHERE table_name = 'app_user'
)`
	err := pgPool.QueryRow(ctx, query).Scan(&exists)
	return exists, err
}

// Check if migrations are required (query will be added later)
func checkMigrationRequired(ctx context.Context, tx pgx.Tx) (bool, error) {
	// TODO: Add query logic
	var exists bool

	err := tx.QueryRow(ctx, `
SELECT EXISTS (
	SELECT 1
	FROM information_schema.columns
	WHERE table_schema = 'public'
	AND table_name = 'user_event'
	AND column_name = 'user_event_id'
)
`).Scan(&exists)
	if err != nil {
		return false, err
	}

	return !exists, nil

}
