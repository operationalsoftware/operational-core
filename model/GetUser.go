package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetUser(db *sql.DB, id int) User {
	getUserQuery := `
SELECT
	UserID,
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName
FROM
	User
WHERE
	UserID = ?
	`

	var user User
	err := db.QueryRow(getUserQuery, id).Scan(
		&user.UserId,
		&user.IsAPIUser,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		log.Panic(err)
	}

	return user
}
