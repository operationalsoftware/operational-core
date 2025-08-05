package migrate

import (
	"app/assets"
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/db"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

func Run() error {
	pgEnv := db.LoadPostgresEnv()

	ctx := context.Background()

	postgresConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres",
		pgEnv.User, pgEnv.Password, pgEnv.Host, pgEnv.Port)

	postgresConn, err := pgx.Connect(ctx, postgresConnStr)
	if err != nil {
		return fmt.Errorf("Failed to connect to postgres database: %v", err)
	}
	defer postgresConn.Close(ctx)

	err = ensureDatabaseExists(ctx, postgresConn, pgEnv.Database)
	if err != nil {
		return fmt.Errorf("error ensuring database exists: %v", err)
	}

	applicationConnStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		pgEnv.User, pgEnv.Password, pgEnv.Host, pgEnv.Port, pgEnv.Database)

	applicationConn, err := pgx.Connect(ctx, applicationConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to application database %s: %v", pgEnv.Database, err)
	}
	defer applicationConn.Close(ctx)

	err = ensureMigrationLogTableExists(ctx, applicationConn)
	if err != nil {
		return fmt.Errorf("error ensuring migration log table exists: %v", err)
	}

	err = ensureMigrationsApplied(ctx, applicationConn)
	if err != nil {
		return fmt.Errorf("error ensuring migrations applied: %v", err)
	}

	err = ensureSystemUserExists(ctx, applicationConn)
	if err != nil {
		return fmt.Errorf("error ensuring system user exists: %v", err)
	}

	return nil

}

func ensureDatabaseExists(ctx context.Context, conn *pgx.Conn, dbName string) error {

	var exists bool
	err := conn.QueryRow(ctx, `
SELECT EXISTS (
	SELECT FROM pg_database WHERE datname = $1
)
	`, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error querying for application database %s: %w", dbName, err)
	}

	if exists {
		return nil // Database already exists, nothing to do
	}

	fmt.Printf("Database %s does not exist. Creating... ", dbName)
	_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err != nil {
		return fmt.Errorf("error creating database %s: %w", dbName, err)
	}
	fmt.Println("done.")

	return nil
}

func ensureMigrationLogTableExists(ctx context.Context, conn *pgx.Conn) error {

	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS migration_log (
	script_id TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	script TEXT NOT NULL
);
	`)

	if err != nil {
		return fmt.Errorf("error creating migration_log table: %w", err)
	}

	return nil

}

func ensureMigrationsApplied(ctx context.Context, conn *pgx.Conn) error {

	scriptsDirName := "/internal/migrate/scripts"

	unappliedMigrationScriptIDs, err := getUnappliedMigrationScriptIDs(ctx, conn, scriptsDirName)
	if err != nil {
		return fmt.Errorf("error getting unapplied migration script ids: %v", err)
	}

	for _, scriptID := range unappliedMigrationScriptIDs {
		fmt.Printf("Applying migration: %s... ", scriptID)
		err := applyMigrationScript(ctx, conn, scriptsDirName, scriptID)
		if err != nil {
			return err
		}
		fmt.Println("done.")
	}

	return nil
}

func getUnappliedMigrationScriptIDs(ctx context.Context, conn *pgx.Conn, scriptsDirName string) ([]string, error) {
	allMigrationScriptIDs, err := getAllMigrationScriptIDs(scriptsDirName)
	if err != nil {
		return []string{}, fmt.Errorf("error getting all migration script ids: %v", err)
	}

	highestAppliedScriptID, err := getHighestAppliedMigrationScriptID(ctx, conn)
	if err != nil {
		return []string{}, fmt.Errorf("error getting highest applied migration script id: %v", err)
	}

	highestAppliedIdx := slices.Index(allMigrationScriptIDs, highestAppliedScriptID)

	unappliedMigrationScriptIDs := allMigrationScriptIDs[highestAppliedIdx+1:]

	return unappliedMigrationScriptIDs, nil
}

func getAllMigrationScriptIDs(scriptsDirName string) ([]string, error) {

	scriptIDs := []string{}
	for path := range assets.Assets.Files {
		if strings.HasPrefix(path, scriptsDirName) &&
			strings.HasSuffix(path, ".sql") {
			scriptID := path[len(scriptsDirName)+1:]
			scriptIDs = append(scriptIDs, scriptID)
		}
	}

	sort.Strings(scriptIDs)

	return scriptIDs, nil
}

func getHighestAppliedMigrationScriptID(ctx context.Context, conn *pgx.Conn) (string, error) {
	var scriptID string
	err := conn.QueryRow(ctx, `
SELECT
	script_id
FROM
	migration_log
ORDER BY
	script_id DESC
LIMIT 1
		`).Scan(&scriptID)

	if err == pgx.ErrNoRows {
		// no migrations applied yet
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to query highest applied migration script id: %v", err)
	}

	return scriptID, err
}

func applyMigrationScript(ctx context.Context, conn *pgx.Conn, scriptsDirName, scriptID string) error {

	assetPath := path.Join(scriptsDirName, scriptID)
	scriptFile, err := assets.Assets.Open(assetPath)
	if err != nil {
		fmt.Printf("error opening file %s: %s\n", assetPath, err)
		panic(err)
	}
	defer scriptFile.Close()

	scriptFileBytes, err := io.ReadAll(scriptFile)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string
	script := string(scriptFileBytes)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beggining transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, script)
	if err != nil {
		return fmt.Errorf("error running migration script %s: %v", scriptID, err)
	}

	_, err = tx.Exec(ctx, `
INSERT INTO
	migration_log (
		script_id,
		script
	) VALUES (
		$1,
		$2
	);
		`, scriptID, script)
	if err != nil {
		return fmt.Errorf("error adding entry to migration_log: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func ensureSystemUserExists(ctx context.Context, conn *pgx.Conn) error {
	fmt.Print("Creating system user... ")

	userRepository := repository.NewUserRepository()

	systemUser, err := userRepository.GetUserByUsername(ctx, conn, "system")
	if err != nil {
		return fmt.Errorf("error getting system user by username: %v", err)
	}
	if systemUser != nil {
		// already exists
		return nil
	}

	systemPassword := os.Getenv("SYSTEM_USER_PASSWORD")
	if systemPassword == "" {
		return fmt.Errorf("system password not set in .env file")
	}

	var newAPIUser = model.NewAPIUser{
		Username: "system",
		Password: systemPassword,
		Permissions: model.UserPermissions{
			UserAdmin: model.UserAdminPermissions{
				Access: true,
			},
		},
	}

	password, err := userRepository.CreateAPIUser(ctx, conn, newAPIUser)
	if err != nil {
		return err
	}

	fmt.Println("done.")
	fmt.Println("System user details:\n\tusername: system\n\tpassword: " + password)

	return nil
}
