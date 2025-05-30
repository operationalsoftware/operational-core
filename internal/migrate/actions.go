package migrate

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// Apply migrations if required (queries will be added later)
func migrate(ctx context.Context, tx pgx.Tx) error {

	// Team table
	_, err := tx.Exec(ctx, `
CREATE TABLE team (
	team_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	team_name TEXT NOT NULL,
	is_archived BOOLEAN NOT NULL DEFAULT FALSE
);
		`)
	if err != nil {
		return err
	}
	// Team table end

	// Andon Issue table
	_, err = tx.Exec(context.Background(), `
CREATE TABLE andon_issue (
    andon_issue_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    issue_name TEXT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id INTEGER REFERENCES andon_issue(andon_issue_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES app_user(user_id),
    updated_at TIMESTAMPTZ,
    updated_by INTEGER REFERENCES app_user(user_id)
);
		`)
	if err != nil {
		return err
	}
	// Andon Issue table end

	return nil
}

// Initialise the database (actual queries will be added later)
func initialise(ctx context.Context, tx pgx.Tx) error {

	tx, err := tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	fmt.Println("Checking for system password in .env... ")
	systemPassword := os.Getenv("SYSTEM_USER_PASSWORD")
	if systemPassword == "" {
		return fmt.Errorf("system password not set in .env file")
	}

	// Create users table
	fmt.Print("Creating User table... ")
	_, err = tx.Exec(context.Background(), `
CREATE TABLE app_user (
	user_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    is_api_user BOOLEAN DEFAULT FALSE NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE, 
    first_name TEXT,
    last_name TEXT,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMPTZ DEFAULT NULL,
    hashed_password TEXT NOT NULL,
    failed_login_attempts INTEGER DEFAULT 0 NOT NULL,
    login_blocked_until TIMESTAMPTZ DEFAULT NULL,
    permissions JSONB DEFAULT '{}'::JSONB NOT NULL,
    user_data JSONB DEFAULT '{}'::JSONB NOT NULL,
	session_duration_minutes INT
);
`)
	if err != nil {
		return err
	}
	fmt.Println("done")

	userRepository := repository.NewUserRepository()

	// add the system user with a random password
	fmt.Print("Creating system user... ")

	var newAPIUser = model.NewAPIUser{
		Username: "system",
		Password: systemPassword,
		Permissions: model.UserPermissions{
			UserAdmin: model.UserAdminPermissions{
				Access: true,
			},
		},
	}

	password, err := userRepository.CreateAPIUser(ctx, tx, newAPIUser)
	if err != nil {
		return err
	}

	fmt.Println("done")
	fmt.Println("System user details:\n\tusername: system\n\tpassword: " + password)

	// Create Recent search table
	_, err = tx.Exec(context.Background(), `
CREATE TABLE recent_search (
	recent_search_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	search_term TEXT NOT NULL,
	search_entities TEXT[] NOT NULL,
	user_id INT REFERENCES app_user(user_id),
	last_searched_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT unique_search_per_user UNIQUE (search_term, search_entities, user_id)
);
		`)
	if err != nil {
		return err
	}
	// Recent search table end

	// Team table
	_, err = tx.Exec(context.Background(), `
CREATE TABLE team (
	team_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	team_name TEXT NOT NULL,
	is_archived BOOLEAN NOT NULL DEFAULT FALSE
);
		`)
	if err != nil {
		return err
	}
	// Team table end

	// Andon Issue table
	_, err = tx.Exec(context.Background(), `
CREATE TABLE andon_issue (
    andon_issue_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    issue_name TEXT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id INTEGER REFERENCES andon_issue(andon_issue_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER NOT NULL REFERENCES app_user(user_id),
    updated_at TIMESTAMPTZ,
    updated_by INTEGER REFERENCES app_user(user_id)
);
		`)
	if err != nil {
		return err
	}
	// Andon Issue table end

	//
	// END OF INITIALISATION
	//

	fmt.Println("done")

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Create the database if it does not exist
func createDatabase(ctx context.Context, conn *pgx.Conn, dbName string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, dbName))
	return err
}
