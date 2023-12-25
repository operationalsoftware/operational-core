package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetUser(db *sql.DB, id string) User {
	user := User{}
	getUserQuery := `
SELECT user_id, first_name, last_name, email, username FROM users WHERE user_id = ?
	`

	rows, err := db.Query(getUserQuery, id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var dbUser User
		err := rows.Scan(&dbUser.UserId, &dbUser.FirstName, &dbUser.LastName, &dbUser.Email, &dbUser.Username)
		if err != nil {
			log.Fatal(err)
		}
		user = dbUser
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return user
}
