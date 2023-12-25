package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func EditUser(db *sql.DB, user User, id string) error {
	editUserQuery := `
UPDATE users SET first_name = ?, last_name = ?, email = ?, username = ? WHERE user_id = ?
	`

	_, err := db.Exec(editUserQuery, user.FirstName, user.LastName, user.Email, user.Username, id)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
