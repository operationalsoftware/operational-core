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

	// 	// Alter column for app_user table
	// 	_, err := tx.Exec(context.Background(), `
	// CREATE TABLE file (
	// 	file_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	// 	file_name TEXT,
	// 	mime_type TEXT,
	// 	file_ext TEXT,
	// 	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	// );
	// 	`)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// end

	// Create stock_item table
	_, err := tx.Exec(context.Background(), `
CREATE TABLE stock_item (
	stock_code TEXT PRIMARY KEY,
	description TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
	`)
	if err != nil {
		return err
	}
	// end

	// Create stock_item table
	_, err = tx.Exec(context.Background(), `
CREATE TABLE stock_item_change (
	stock_item_change_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	stock_code TEXT NOT NULL REFERENCES stock_item(stock_code) ON UPDATE CASCADE, 
	stock_code_history TEXT, 
	description TEXT,
	change_by INT REFERENCES app_user(user_id),
	changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
	`)
	if err != nil {
		return err
	}

	// Create sku_config table
	_, err = tx.Exec(context.Background(), `
CREATE TABLE sku_config (
	sku_config_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	sku_field TEXT NOT NULL,
	label TEXT NOT NULL,
	code TEXT NOT NULL,
	CONSTRAINT unique_sku_field_label UNIQUE (sku_field, label),
  	CONSTRAINT unique_sku_field_code UNIQUE (sku_field, code)
);
	`)
	if err != nil {
		return err
	}

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
	// end

	// Create Stock ledger tables
	_, err = tx.Exec(context.Background(), `
			CREATE TABLE stock_transaction (
				stock_transaction_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
				transaction_type TEXT NOT NULL,
				transaction_by INT NOT NULL REFERENCES app_user(user_id),
				transaction_note TEXT NOT NULL,
				timestamp TIMESTAMPTZ NOT NULL
			);
		`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `
			CREATE TABLE stock_transaction_entry (
				stock_transaction_entry_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
				account TEXT NOT NULL,
				stock_code TEXT NOT NULL,
				location TEXT NOT NULL,
				bin TEXT NOT NULL,
				lot_number TEXT NOT NULL,
				quantity NUMERIC NOT NULL,
				running_total NUMERIC NOT NULL,
				stock_transaction_id INT NOT NULL REFERENCES stock_transaction(stock_transaction_id)
			);
		`)
	if err != nil {
		return err
	}
	// end Stock ledger

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
