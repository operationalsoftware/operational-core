package migrate

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/db"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Apply migrations if required (queries will be added later)
func migrate(ctx context.Context, pgPool *pgxpool.Pool) error {
	// Create users table
	_, err := pgPool.Exec(context.Background(), `
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
    user_data JSONB DEFAULT '{}'::JSONB NOT NULL
);
`)
	if err != nil {
		return err
	}

	// Populate users table
	sqLiteDB := db.UseDB()

	rows, err := sqLiteDB.Query(`
SELECT
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName,
	Created,
	LastLogin,
	HashedPassword,
	FailedLoginAttempts,
	LoginBlockedUntil,
	Permissions,
	UserData
FROM
	User
ORDER BY UserID ASC;
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			isAPIUser           bool
			username            string
			email               sql.NullString
			firstName           sql.NullString
			lastName            sql.NullString
			created             time.Time
			lastLogin           sql.NullTime
			hashedPassword      string
			failedLoginAttempts int
			loginBlockedUntil   sql.NullTime
			permissions         string
			userData            string
		)

		rows.Scan(
			&isAPIUser,
			&username,
			&email,
			&firstName,
			&lastName,
			&created,
			&lastLogin,
			&hashedPassword,
			&failedLoginAttempts,
			&loginBlockedUntil,
			&permissions,
			&userData,
		)

		pgPool.Exec(
			ctx,
			`
INSERT INTO app_user (
	is_api_user,
	username,
	email,
	first_name,
	last_name,
	created,
	last_login,
	hashed_password,
	failed_login_attempts,
	login_blocked_until,
	permissions,
	user_data
)
VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
			`,
			isAPIUser,
			username,
			email,
			firstName,
			lastName,
			created,
			lastLogin,
			hashedPassword,
			failedLoginAttempts,
			loginBlockedUntil,
			permissions,
			userData,
		)
	}

	// Users table end

	return nil
}

// Initialise the database (actual queries will be added later)
func initialise(ctx context.Context, pgPool *pgxpool.Pool) error {

	tx, err := pgPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

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
    user_data JSONB DEFAULT '{}'::JSONB NOT NULL
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
