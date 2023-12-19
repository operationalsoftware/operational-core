package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func AddUser(db *sql.DB, user User) error {
	insertUserStmt := `
INSERT INTO users 
	(username, email, first_name, last_name)
VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(insertUserStmt, user.Username, user.Email, user.FirstName, user.LastName)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
