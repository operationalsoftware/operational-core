package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetUsers(db *sql.DB) []User {
	users := []User{}
	getUsersQuery := `
SELECT user_id, first_name, last_name, email, username FROM users
	`

	rows, err := db.Query(getUsersQuery)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return users
}
