package migrate

import (
	"app/models/usermodel"
	"database/sql"
	"fmt"
)

func checkInitialisationRequired(tx *sql.Tx) (bool, error) {

	// Check if the initialisation has already been done
	var requiresInitialisation bool
	query := `
SELECT NOT EXISTS (
	SELECT 1 FROM sqlite_master WHERE type='table' AND name='User'
)`
	err := tx.QueryRow(query).Scan(&requiresInitialisation)
	if err != nil {
		return false, err
	}

	return requiresInitialisation, nil
}

func initialiseDB(tx *sql.Tx) error {

	//
	// START OF INITIALISATION
	//

	fmt.Println("Initialising database...")

	// Create users table
	fmt.Print("Creating User table... ")
	_, err := tx.Exec(`
CREATE TABLE User (
	UserID INTEGER PRIMARY KEY AUTOINCREMENT, 
	IsAPIUser BOOLEAN DEFAULT FALSE NOT NULL,
	Username TEXT NOT NULL UNIQUE,
	Email TEXT UNIQUE, 
	FirstName TEXT,
	LastName TEXT,
	Created DATETIME DEFAULT CURRENT_TIMESTAMP,
	LastLogin DATETIME DEFAULT NULL,
	HashedPassword TEXT NOT NULL,
	FailedLoginAttempts INTEGER DEFAULT 0 NOT NULL,
	LoginBlockedUntil DATETIME DEFAULT NULL,
	Permissions JSON DEFAULT '{}' NOT NULL,
	UserData JSON DEFAULT '{}' NOT NULL
);
`)
	if err != nil {
		return err
	}
	fmt.Println("done")

	// add the system user with a random password
	fmt.Print("Creating system user... ")

	password, err := usermodel.GenerateRandomPassword(24)
	if err != nil {
		panic(err)
	}

	var userToAdd = usermodel.NewAPIUser{
		Username: "system",
		Password: password,
		Permissions: usermodel.UserPermissions{
			UserAdmin: usermodel.UserAdminPermissions{
				Access: true,
			},
		},
	}

	err = usermodel.AddAPIUser(tx, userToAdd)
	if err != nil {
		return err
	}

	fmt.Println("done")
	fmt.Println("System user details:\n\tusername: system\n\tpassword: " + password)

	//
	// END OF INITIALISATION
	//

	fmt.Println("done")

	return nil
}
