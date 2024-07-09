package migrate

import (
	"app/db"
	"app/models/usermodel"
	"fmt"
	"log"
)

func initialise() bool {
	db := db.UseDB()

	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	//
	// START OF INITIALISATION
	//

	fmt.Println("Initialising database...")

	// Create users table
	fmt.Print("Creating User table... ")
	_, err = db.Exec(`
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
		return false
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
		return false
	}

	fmt.Println("done")
	fmt.Println("System user details:\n\tusername: system\n\tpassword: " + password)

	//
	// END OF INITIALISATION
	//
	fmt.Print("Database initialised, committing changes...")
	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("done")

	return true
}
