package migrate

import (
	"app/internal/models"
	"app/internal/services/userservice"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Apply migrations if required (queries will be added later)
func migrate(ctx context.Context, pgPool *pgxpool.Pool) error {
	// TODO: Add migration logic
	return nil
}

// Initialise the database (actual queries will be added later)
func initialise(ctx context.Context, pgPool *pgxpool.Pool) error {

	// Create users table
	fmt.Print("Creating User table... ")
	_, err := pgPool.Exec(context.Background(), `
CREATE TABLE app_user (
	user_id SERIAL PRIMARY KEY,
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
    user_data JSONB DEFAULT '{}'::JSONB NOT NULL
);
`)
	if err != nil {
		return err
	}
	fmt.Println("done")

	userService := userservice.NewUserService(pgPool)

	// add the system user with a random password
	fmt.Print("Creating system user... ")

	var newAPIUser = models.NewAPIUser{
		Username: "system",
		Permissions: models.UserPermissions{
			UserAdmin: models.UserAdminPermissions{
				Access: true,
			},
		},
	}

	password, err := userService.CreateAPIUser(ctx, newAPIUser)
	if err != nil {
		return err
	}

	fmt.Println("done")
	fmt.Println("System user details:\n\tusername: system\n\tpassword: " + password)

	//
	// END OF INITIALISATION
	//

	fmt.Println("done")

	// TODO: Add initialisation logic
	return nil
}

// Create the database if it does not exist
func createDatabase(ctx context.Context, conn *pgx.Conn, dbName string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, dbName))
	return err
}
